package actions

import (
	"log"
	"os"
	"path/filepath"

	"github.com/deis/deisrel/git"
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
	// IncludeClosed represents the '--includeClosed' flag
	IncludeClosed = "includeClosed"
)

const (
	generateParamsFileName = "generate_params.toml"
)

type releaseName struct {
	Full  string
	Short string
}

var (
	// additionalGitRepoNames represents the repo names lacking representation
	// in any helm chart, yet still requiring updates during each Workflow
	// release, including changelog generation and creation of git tags.
	additionalGitRepoNames = []string{"workflow", "charts"}

	// allGitRepoNames represent all GitHub repo names needing git-based updates for a release
	allGitRepoNames = append(git.RepoNames(), additionalGitRepoNames...)

	repoNames      = git.RepoNames()
	componentNames = git.ComponentNames()
	// TODO: https://github.com/deis/deisrel/issues/12
	repoToComponentNames = git.RepoToComponentNames()

	deisRelease = releaseName{
		Full:  os.Getenv("WORKFLOW_RELEASE"),
		Short: os.Getenv("WORKFLOW_RELEASE_SHORT"),
	}
	defaultStagingPath = getFullPath("staging")
)

func getFullPath(dirName string) string {
	currentWorkingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working dir (%s)", err)
	}
	return filepath.Join(currentWorkingDir, dirName)
}
