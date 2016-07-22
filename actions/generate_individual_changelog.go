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
		sha := c.Args().Get(2)
		vals := &changelog.Values{
			OldRelease: c.Args().Get(1),
			NewRelease: c.Args().Get(3),
		}
		if vals.OldRelease == "" || vals.NewRelease == "" || sha == "" || repoName == "" {
			log.Fatal("Usage: changelog individual <repo> <old-release> <sha> <new-release>")
		}
		skippedCommits, err := changelog.SingleRepoVals(client, vals, sha, repoName, false)

		if len(skippedCommits) > 0 {
			for _, ci := range skippedCommits {
				fmt.Fprintln(os.Stderr, "skipping commit", ci)
			}
		}

		if err != nil {
			log.Fatalf("could not generate changelog: %s", err)
		}
		if err := changelog.Tpl.Execute(dest, vals); err != nil {
			log.Fatalf("could not template changelog: %s", err)
		}
		return nil
	}
}
