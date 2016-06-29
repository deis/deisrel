package helm

import (
	"path/filepath"

	"github.com/google/go-github/github"
)

// GetWorkflowChart returns a helmChart that represents the workflow chart and
// the files in it that need updating for a release
func GetWorkflowChart(ghClient *github.Client) (*Chart, error) {
	return newChart(ghClient, "workflow-dev", []string{
		"README.md",
		"Chart.yaml",
	})
}

func GetWorkflowE2EChart(ghClient *github.Client) (*Chart, error) {
	return newChart(ghClient, "workflow-dev-e2e", []string{
		"README.md",
		"Chart.yaml",
		filepath.Join("tpl", "generate_params.toml"),
	})
}

func GetRouterChart(ghClient *github.Client) (*Chart, error) {
	return newChart(ghClient, "router-dev", []string{
		"README.md",
		"Chart.yaml",
	})
}
