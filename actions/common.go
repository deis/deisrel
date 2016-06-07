package actions

import (
	"log"
	"os"
	"path/filepath"
	"text/template"
)

const (
	// TagFlag represents the '-tag' flag
	TagFlag = "tag"
	// PullPolicyFlag represents the '-pull-policy' flag
	PullPolicyFlag = "pull-policy"
	// OrgFlag represents the '-org' flag
	OrgFlag = "org"
	// ShaFilepathFlag represents the --sha-filepath flag
	ShaFilepathFlag = "sha-filepath"
	// YesFlag represents the --yes flag
	YesFlag = "yes"
	// RepoFlag represents the '-repo' flag
	RepoFlag = "repo"
	// RefFlag represents the '-ref' flag (for specifying a SHA, branch or tag)
	RefFlag = "ref"
	// GHOrgFlag represents the '-ghOrg' flag
	GHOrgFlag = "ghOrg"
	// StagingDirFlag represents the '-stagingDir' flag
	StagingDirFlag = "stagingDir"
)

const (
	generateParamsFileName = "generate_params.toml"
)

type helmChart struct {
	Name     string
	Template *template.Template
	Files    []string
}

type releaseName struct {
	Full  string
	Short string
}

var (
	// TODO: https://github.com/deis/deisrel/issues/12
	repoToComponentNames = map[string][]string{
		"builder":          {"Builder"},
		"controller":       {"Controller"},
		"dockerbuilder":    {"DockerBuilder"},
		"fluentd":          {"FluentD"},
		"monitor":          {"InfluxDB", "Grafana", "Telegraf"},
		"logger":           {"Logger"},
		"minio":            {"Minio"},
		"postgres":         {"Database"},
		"registry":         {"Registry"},
		"router":           {"Router"},
		"slugbuilder":      {"SlugBuilder"},
		"slugrunner":       {"SlugRunner"},
		"stdout-metrics":   {"StdoutMetrics"},
		"workflow-e2e":     {"WorkflowE2E"},
		"workflow-manager": {"WorkflowManager"},
	}

	componentToImageName = map[string]string{
		"Builder":         "builder",
		"Controller":      "controller",
		"DockerBuilder":   "dockerbuilder",
		"FluentD":         "fluentd",
		"InfluxDB":        "influxdb",
		"Grafana":         "grafana",
		"Telegraf":        "telegraf",
		"Logger":          "logger",
		"Minio":           "minio",
		"Database":        "postgres",
		"Registry":        "registry",
		"Router":          "router",
		"SlugBuilder":     "slugbuilder",
		"SlugRunner":      "slugrunner",
		"StdoutMetrics":   "stdout-metrics",
		"WorkflowE2E":     "workflow-e2e",
		"WorkflowManager": "workflow-manager",
	}

	repoNames = getRepoNames(repoToComponentNames)

	// additionalGitRepoNames represents the repo names lacking representation
	// in any helm chart, yet still requiring updates during each Workflow
	// release, including changelog generation and creation of git tags.
	additionalGitRepoNames = []string{"workflow", "charts"}

	// allGitRepoNames represent all GitHub repo names needing git-based updates for a release
	allGitRepoNames = append(repoNames, additionalGitRepoNames...)

	componentNames = getComponentNames(repoToComponentNames)

	deisRelease = releaseName{
		Full:  os.Getenv("WORKFLOW_RELEASE"),
		Short: os.Getenv("WORKFLOW_RELEASE_SHORT"),
	}
	defaultStagingPath = getFullPath("staging")

	// RouterChart represents the router chart and its files needing updating
	// for a release
	RouterChart = helmChart{
		Name:     "router-dev",
		Template: generateParamsRouterTpl,
		Files: []string{
			"README.md",
			"Chart.yaml",
		},
	}

	// WorkflowChart represents the workflow chart and its files needing updating
	// for a release
	WorkflowChart = helmChart{
		Name:     "workflow-dev",
		Template: generateParamsTpl,
		Files: []string{
			"README.md",
			"Chart.yaml",
		},
	}

	// WorkflowE2EChart represents the workflow e2e chart and its files needing updating
	// for a release
	WorkflowE2EChart = helmChart{
		Name:     "workflow-dev-e2e",
		Template: generateParamsE2ETpl,
		Files: []string{
			"README.md",
			"Chart.yaml",
			filepath.Join("tpl", "workflow-e2e-pod.yaml"),
		},
	}
)

func getRepoNames(repoToComponentNames map[string][]string) []string {
	repoNames := make([]string, 0, len(repoToComponentNames))
	for repoName := range repoToComponentNames {
		repoNames = append(repoNames, repoName)
	}
	return repoNames
}

func getComponentNames(repoToComponentNames map[string][]string) []string {
	var ret []string
	for _, componentNames := range repoToComponentNames {
		for _, componentName := range componentNames {
			ret = append(ret, componentName)
		}
	}
	return ret
}

func getFullPath(dirName string) string {
	currentWorkingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working dir (%s)", err)
	}
	return filepath.Join(currentWorkingDir, dirName)
}
