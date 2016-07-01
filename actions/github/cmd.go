package github

import (
	"github.com/codegangsta/cli"
	"github.com/google/go-github/github"
)

// Command returns the CLI command for 'deisrel github ...' commands
func Command(ghClient *github.Client) cli.Command {
	return cli.Command{
		Name: "github",
		Subcommands: []cli.Command{
			cli.Command{
				Name:   "create-releases",
				Action: createReleasesCmd(ghClient),
			},
		},
	}
}
