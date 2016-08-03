package actions

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/deis/deisrel/changelog"
	"github.com/google/go-github/github"
	"github.com/urfave/cli"
)

// GenerateIndividualChangelog is the CLI action for creating a changelog for a single repo
func GenerateIndividualChangelog(client *github.Client, dest io.Writer) func(*cli.Context) error {
	return func(c *cli.Context) error {
		repoName := c.Args().Get(0)
		sha := c.String("sha")
		vals := &changelog.Values{
			OldRelease: c.String("base-tag"),
			NewRelease: c.Args().Get(1),
		}
		if vals.NewRelease == "" || repoName == "" {
			log.Fatal("Usage: changelog individual <repo> <new-release>")
		}

		// If sha isn't set, use the latest commit on master
		if sha == "" {
			master, _, err := client.Repositories.GetBranch("deis", repoName, "master")
			if err != nil {
				return err
			}
			sha = *master.Commit.SHA
		}

		// If base-tag isn't set, use the most recent in the repository
		if vals.OldRelease == "" {
			tags, _, err := client.Repositories.ListTags("deis", repoName, nil)
			if err != nil {
				return err
			}

			if len(tags) < 1 {
				vals.OldRelease = "none"
			} else {
				vals.OldRelease = *tags[0].Name
			}
		}

		skippedCommits, err := changelog.SingleRepoVals(client, vals, sha, repoName, false)

		if len(skippedCommits) > 0 {
			for _, ci := range skippedCommits {
				fmt.Fprintln(os.Stderr, "skipping commit", ci)
			}
		}

		// Ecape secquences for color
		g := "\033[0;32m"
		b := "\033[0;34m"
		r := "\033[0m"
		fmt.Printf("\n%sCreating changelog for %s with tag %s through commit %s\n\n", g, b+repoName+g, b+vals.OldRelease+g, b+sha+r)

		if err != nil {
			log.Fatalf("could not generate changelog: %s", err)
		}
		if err := changelog.Tpl.Execute(dest, vals); err != nil {
			log.Fatalf("could not template changelog: %s", err)
		}
		return nil
	}
}
