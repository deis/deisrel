package git

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

	repoNames      = getRepoNames(repoToComponentNames)
	componentNames = getComponentNames(repoToComponentNames)
)

// RepoNames returns a slice of known repository names
func RepoNames() []string {
	return repoNames
}

// ComponentNames returns a slice of known component names
func ComponentNames() []string {
	return componentNames
}

// RepoToComponentNames returns a mapping from each repository name to its known component names
func RepoToComponentNames() map[string][]string {
	return repoToComponentNames
}

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
