package main

import (
	"log"
	"os"

	"github.com/arschles/deisrel/actions"
	"github.com/codegangsta/cli"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func main() {
	ghTkn := os.Getenv("GITHUB_TOKEN")
	if ghTkn == "" {
		log.Fatalf("'GITHUB_TOKEN' env var required")
	}
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: ghTkn})
	cl := oauth2.NewClient(oauth2.NoContext, ts)
	ghClient := github.NewClient(cl)

	app := cli.NewApp()
	app.Name = "deisrel"
	app.Usage = "Utilities for releasing a new Deis version"
	app.Commands = []cli.Command{
		cli.Command{
			Name: "git",
			Subcommands: []cli.Command{
				cli.Command{
					Name:   "shas",
					Action: actions.GetShas(ghClient),
					Flags: []cli.Flag{
						cli.BoolFlag{
							Name:  actions.ShortFlag,
							Usage: "Whether to show short 7 character shas",
						},
					},
				},
				cli.Command{
					Name:    "generate-changelogs",
					Aliases: []string{"gen-ch"},
					Action:  actions.GenerateChangelogs(ghClient),
				},
			},
		},
		cli.Command{
			Name: "helm-params",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  actions.TagFlag,
					Value: "",
					Usage: "The Docker tag to apply for all images. If empty, defaults to 'git-$SHORT_SHA' for each respective component",
				},
				cli.StringFlag{
					Name:  actions.PullPolicyFlag,
					Value: "IfNotPresent",
					Usage: "The 'imagePullPolicy' value to set on each image",
				},
				cli.StringFlag{
					Name:  actions.OrgFlag,
					Value: "deis",
					Usage: "The docker repository organization to set on each image",
				},
			},
			Subcommands: []cli.Command{
				cli.Command{
					Name:   "e2e",
					Action: actions.HelmGenerateE2E(ghClient),
				},
				cli.Command{
					Name:   "workflow",
					Action: actions.HelmGenerateWorkflow(ghClient),
				},
			},
		},
	}

	app.Run(os.Args)
}
