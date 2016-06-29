package actions

import (
	"github.com/codegangsta/cli"
	"github.com/deis/deisrel/helm"
	"github.com/google/go-github/github"
)

// HelmStageE2E is the cli handler for generating a release helm chart for deis-e2e
func HelmStageE2E(ghClient *github.Client) func(*cli.Context) error {
	return func(c *cli.Context) error {
		e2eChart, err := helm.GetWorkflowE2EChart(ghClient)
		if err != nil {
			return err
		}
		helmStage(ghClient, c, *e2eChart)
		return nil
	}
}
