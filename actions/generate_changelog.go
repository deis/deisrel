package actions

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"sync"

	"github.com/deis/deisrel/changelog"
	"github.com/google/go-github/github"
	"github.com/urfave/cli"
)

// GenerateChangelog is the CLI action for creating an aggregated changelog from all of the Deis Workflow repos.
func GenerateChangelog(client *github.Client, dest io.Writer) func(*cli.Context) error {
	return func(c *cli.Context) error {
		repoMapFile := c.Args().Get(0)
		oldTag := c.Args().Get(1)
		newTag := c.Args().Get(2)
		if repoMapFile == "" || oldTag == "" || newTag == "" {
			log.Fatal("Usage: changelog global <repo map> <old-release> <new-release>")
		}

		out, err := ioutil.ReadFile(repoMapFile)
		if err != nil {
			log.Fatal(err.Error())
		}

		repoMap := make(map[string][]string)
		err = json.Unmarshal(out, &repoMap)
		if err != nil {
			log.Fatal(err.Error())
		}

		vals, errs := generateChangelogVals(client, repoMap, oldTag, newTag)
		if len(errs) > 0 {
			for _, err := range errs {
				log.Printf("Error: %s", err)
			}
		}
		if err := changelog.Tpl.Execute(dest, changelog.MergeValues(oldTag, newTag, vals)); err != nil {
			log.Fatalf("could not template changelog: %s", err)
		}
		return nil
	}
}

func generateChangelogVals(client *github.Client, repoMap map[string][]string, oldTag, newTag string) ([]changelog.Values, []error) {
	var wg sync.WaitGroup
	done := make(chan bool)
	valsCh := make(chan changelog.Values)
	errCh := make(chan error)
	defer close(errCh)
	for repo := range repoMap {
		wg.Add(1)
		go func(repo string) {
			defer wg.Done()
			vals := &changelog.Values{OldRelease: oldTag, NewRelease: newTag}
			_, err := changelog.SingleRepoVals(client, vals, newTag, repo, true)
			if err != nil {
				errCh <- err
				return
			}
			valsCh <- *vals
		}(repo)
	}
	go func() {
		// wait for all fetches from github to be complete before returning
		wg.Wait()
		close(done)
	}()

	vals := []changelog.Values{}
	errs := []error{}
	for {
		select {
		case <-done:
			return vals, errs
		case val := <-valsCh:
			vals = append(vals, val)
		case err := <-errCh:
			errs = append(errs, err)
		}
	}
}
