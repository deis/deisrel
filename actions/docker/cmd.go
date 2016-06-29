package docker

import (
	"github.com/codegangsta/cli"
	"github.com/deis/deisrel/actions"
	"github.com/deis/deisrel/docker"
	"github.com/google/go-github/github"
)

// Command returns the entire set of subcommands for the 'deisrel docker ...' command
func Command(ghClient *github.Client, dockerCl docker.Client) cli.Command {
	return cli.Command{
		Name: "docker",
		Subcommands: []cli.Command{
			cli.Command{
				Name:        "retag",
				Description: "This command pulls specific Docker images for each corresponding repository's Git SHA, then retags each one to a uniform release tag",
				Usage:       "Retag each specific Docker image from a repo-specific Git SHA to a uniform release tag",
				Flags: []cli.Flag{
					cli.BoolFlag{
						Name:  actions.YesFlag,
						Usage: "If true, skip the prompt to confirm that newly-tagged images will be pushed",
					},
					cli.StringFlag{
						Name:  actions.ShaFilepathFlag,
						Value: "",
						Usage: "the file path which to read in the shas to release",
					},
					cli.StringFlag{
						Name:  newOrgFlag,
						Usage: "The Docker registry organization for the new tagged images (default: deis)",
					},
					cli.StringFlag{
						Name:  actions.RefFlag,
						Value: "master",
						Usage: "Optional ref to add to GitHub repo request (can be SHA, branch or tag)",
					},
					cli.StringSliceFlag{
						Name:  registriesFlag,
						Value: &defaultDockerRegistriesStringSlice,
						Usage: "The docker registries to tag and push to. Use 'index.docker.io' to indicate the docker hub",
					},
				},
				Action: retagCmd(ghClient, dockerCl),
			},
		},
	}
}
