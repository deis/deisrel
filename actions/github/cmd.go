package github

import (
	"github.com/codegangsta/cli"
	"github.com/deis/deisrel/actions"
	"github.com/google/go-github/github"
)

// Command returns the CLI command for 'deisrel github ...' commands
func Command(ghClient *github.Client) cli.Command {
	return cli.Command{
		Name: "github",
		Subcommands: []cli.Command{
			cli.Command{
				Name:        "create-releases",
				Action:      createReleasesCmd(ghClient),
				Usage:       "Create Github releases for each known repository",
				Description: "This command creates a new release, tag & appropriate changelog for all repositories",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  releaseTagFlag,
						Usage: "The tag to create (or apply, if the tag already exists) for this release.",
					},
					cli.StringFlag{
						Name:  releaseTitleFlag,
						Usage: "The title of this release.",
					},
					cli.StringFlag{
						Name:  actions.RefFlag,
						Value: "master",
						Usage: "The Git ref from which to create the release",
					},
					cli.BoolFlag{
						Name:  actions.YesFlag,
						Usage: "If true, proceed with actually creating releases",
					},
				},
			},
		},
	}
}
