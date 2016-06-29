package actions

import (
	"github.com/codegangsta/cli"
	"github.com/deis/deisrel/helm"
	"github.com/google/go-github/github"
)

// HelmStageRouter is the cli handler for generating a release helm chart for router
func HelmStageRouter(ghClient *github.Client) func(*cli.Context) error {
	return func(c *cli.Context) error {
		routerChart, err := helm.GetRouterChart(ghClient)
		if err != nil {
			return err
		}
		helmStage(ghClient, c, *routerChart)
		return nil
	}
}
