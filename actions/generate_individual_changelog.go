package actions

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/codegangsta/cli"
	"github.com/google/go-github/github"
)

const individualChangelogTplStr string = `{{.OldRelease}} -> {{.NewRelease}}

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

var individualChangelogTpl *template.Template = template.Must(template.New("changelog").Parse(individualChangelogTplStr))

type IndividualChangelog struct {
	OldRelease    string
	NewRelease    string
	Features      []string
	Fixes         []string
	Documentation []string
	Maintenance   []string
}

// GenerateIndividualChangelog is the CLI action for creating a changelog for a single repo
func GenerateIndividualChangelog(client *github.Client, dest io.Writer) func(*cli.Context) error {
	return func(c *cli.Context) error {
		repoName := c.Args().Get(0)
		changelog := &Changelog{
			OldRelease: c.Args().Get(1),
			NewRelease: c.Args().Get(2),
		}
		if changelog.OldRelease == "" || changelog.NewRelease == "" {
			log.Fatal("Usage: changelog individual <repo> <old-release> <new-release>")
		}
		if err := generateIndividualChangelog(client, changelog, repoName); err != nil {
			log.Fatalf("could not generate changelog: %s", err)
		}
		err := individualChangelogTpl.Execute(dest, changelog)
		if err != nil {
			log.Fatalf("could not template changelog: %s", err)
		}
		return nil
	}
}

func generateIndividualChangelog(client *github.Client, changelog *Changelog, name string) error {
	commitCompare, resp, err := client.Repositories.CompareCommits("deis", name, changelog.OldRelease, changelog.NewRelease)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			log.Printf("tag does not exist for this repo; skipping %s", name)
		}
		return fmt.Errorf("could not compare commits %s and %s: %s", changelog.OldRelease, changelog.NewRelease, err)
	}
	for _, commit := range commitCompare.Commits {
		commitMessage := strings.Split(*commit.Commit.Message, "\n")[0]
		changelogMessage := fmt.Sprintf("%s %s: %s", shortShaTransform(*commit.SHA), commitFocus(*commit.Commit.Message), commitTitle(*commit.Commit.Message))
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
	return nil
}
