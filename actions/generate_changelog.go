package actions

import (
	"io"
	"log"
	"sync"

	"github.com/codegangsta/cli"
	"github.com/deis/deisrel/changelog"
	"github.com/google/go-github/github"
)

// GenerateChangelog is the CLI action for creating an aggregated changelog from all of the Deis Workflow repos.
func GenerateChangelog(client *github.Client, dest io.Writer) func(*cli.Context) error {
	return func(c *cli.Context) error {
		oldTag := c.Args().Get(0)
		newTag := c.Args().Get(1)
		if oldTag == "" || newTag == "" {
			log.Fatal("Usage: changelog global <old-release> <new-release>")
		}
		vals, errs := generateChangelogVals(client, oldTag, newTag)
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

func generateChangelogVals(client *github.Client, oldTag, newTag string) ([]changelog.Values, []error) {
	var wg sync.WaitGroup
	done := make(chan bool)
	valsCh := make(chan changelog.Values)
	errCh := make(chan error)
	defer close(errCh)
	for _, name := range repoNames {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			vals := &changelog.Values{OldRelease: oldTag, NewRelease: newTag}
			_, err := changelog.SingleRepoVals(client, vals, newTag, name, true)
			if err != nil {
				errCh <- err
				return
			}
			valsCh <- *vals
		}(name)
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
