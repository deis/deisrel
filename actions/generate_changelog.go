package actions

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"text/template"

	"github.com/codegangsta/cli"
	"github.com/google/go-github/github"
)

const changelogTpl string = `{{.OldRelease}} -> {{.NewRelease}}

# Features

{{range .Features}} - {{.}}
{{else}}No new features for this release.
{{end}}

# Fixes

{{range .Fixes}} - {{.}}
{{else}}No bug fixes for this release.
{{end}}

# Documentation

{{range .Documentation}} - {{.}}
{{else}}No new documentation for this release.
{{end}}

# Maintenance

{{range .Maintenance}} - {{.}}
{{else}}No maintenance required for this release.
{{end}}`

var changelogTemplate *template.Template = template.Must(template.New("changelog").Parse(changelogTpl))

type Changelog struct {
	OldRelease    string
	NewRelease    string
	Features      []string
	Fixes         []string
	Documentation []string
	Maintenance   []string
}

// GenerateChangelog is the CLI action for creating an aggregated changelog from all of the Deis Workflow repos.
func GenerateChangelog(client *github.Client, dest io.Writer) func(*cli.Context) error {
	return func(c *cli.Context) error {
		changelog := &Changelog{
			OldRelease: c.Args().Get(0),
			NewRelease: c.Args().Get(1),
		}
		if changelog.OldRelease == "" || changelog.NewRelease == "" {
			log.Fatal("Usage: changelog global <old-release> <new-release>")
		}
		if err := generateChangelog(client, changelog); err != nil {
			log.Fatalf("could not generate changelog: %s", err)
		}
		err := changelogTemplate.Execute(dest, changelog)
		if err != nil {
			log.Fatalf("could not template changelog: %s", err)
		}
		return nil
	}
}

func generateChangelog(client *github.Client, changelog *Changelog) error {
	var wg sync.WaitGroup
	done := make(chan bool)
	errCh := make(chan error)
	defer close(errCh)
	for _, name := range repoNames {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			commitCompare, resp, err := client.Repositories.CompareCommits("deis", name, changelog.OldRelease, changelog.NewRelease)
			if err != nil {
				if resp.StatusCode == http.StatusNotFound {
					log.Printf("tag does not exist for this repo; skipping %s", name)
					return
				}
				errCh <- fmt.Errorf("could not compare commits %s and %s: %s", changelog.OldRelease, changelog.NewRelease, err)
			}
			for _, commit := range commitCompare.Commits {
				commitMessage := strings.Split(*commit.Commit.Message, "\n")[0]
				changelogMessage := fmt.Sprintf("%s (%s) - %s: %s", shortShaTransform(*commit.SHA), name, commitFocus(*commit.Commit.Message), commitTitle(*commit.Commit.Message))
				if strings.HasPrefix(commitMessage, "feat(") {
					changelog.Features = append(changelog.Features, changelogMessage)
				} else if strings.HasPrefix(commitMessage, "fix(") {
					changelog.Fixes = append(changelog.Fixes, changelogMessage)
				} else if strings.HasPrefix(commitMessage, "docs(") || strings.HasPrefix(commitMessage, "doc(") {
					changelog.Documentation = append(changelog.Documentation, changelogMessage)
				} else if strings.HasPrefix(commitMessage, "chore(") {
					changelog.Maintenance = append(changelog.Maintenance, changelogMessage)
				} else {
					log.Printf("skipping commit %s from %s", *commit.SHA, name)
				}
			}
		}(name)
	}
	go func() {
		// wait for all fetches from github to be complete before returning
		wg.Wait()
		close(done)
	}()

	for {
		select {
		case <-done:
			return nil
		case err := <-errCh:
			return fmt.Errorf("could not generate changelog: %s", err)
		}
	}
}
