package branches

import (
	"github.com/codegangsta/cli"
	"github.com/deis/deisrel/actions"
	"github.com/google/go-github/github"
)

const (
	branchNameFlag = "name"
)

// Command returns the CLI command for all 'deisrel branches ...' commands
func Command(ghClient *github.Client) cli.Command {
	return cli.Command{
		Name: "branches",
		Subcommands: []cli.Command{
			cli.Command{
				Name:        "create",
				Usage:       "Create branches on all repositories that are part of the Deis Workflow platform.",
				Description: "This command creates branches on all repositories that are part of the Deis Workflow platform",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  branchNameFlag,
						Value: "",
						Usage: "The name of the branch to create on all repositories",
					},
					cli.BoolFlag{
						Name:  actions.YesFlag,
						Usage: "Whether to proceed with branch creation. Pass false to do a dry run",
					},
					cli.StringFlag{
						Name:  actions.RefFlag,
						Value: "master",
						Usage: "The ref (branch or SHA) to branch from.",
					},
				},
				Action: createCmd(ghClient),
			},
		},
	}
}
