package actions

import (
	"github.com/codegangsta/cli"
	"github.com/deis/deisrel/helm"
	"github.com/google/go-github/github"
)

// HelmStageWorkflow is the cli handler for generating a release helm chart for workflow
func HelmStageWorkflow(ghClient *github.Client) func(*cli.Context) error {
	return func(c *cli.Context) error {
		workflowChart, err := helm.GetWorkflowChart(ghClient)
		if err != nil {
			return err
		}
		helmStage(ghClient, c, *workflowChart)
		return nil
	}
}
