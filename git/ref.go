package git

import (
	"fmt"

	"github.com/google/go-github/github"
)

func newBranchReference(orgName, repoName, branchName, sha string) *github.Reference {
	return &github.Reference{
		Ref: github.String(fmt.Sprintf("refs/heads/%s", branchName)),
		URL: github.String(fmt.Sprintf("https://api.github.com/repos/%s/%s/git/refs/heads/%s", orgName, repoName, branchName)),
		Object: &github.GitObject{
			Type: github.String("commit"),
			SHA:  github.String(sha),
			URL:  github.String(fmt.Sprintf("https://api.github.com/repos/%s/%s/git/commits/%s", orgName, repoName, sha)),
		},
	}
}
