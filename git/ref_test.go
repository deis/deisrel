package git

import (
	"fmt"
	"testing"

	"github.com/arschles/assert"
)

func TestNewBranchReference(t *testing.T) {
	const (
		orgName    = "testorg"
		repoName   = "testrepo"
		branchName = "testbranch"
		sha        = "testsha"
	)

	ref := newBranchReference(orgName, repoName, branchName, sha)
	assert.Equal(t, *ref.Ref, fmt.Sprintf("refs/heads/%s", branchName), "sha")
	assert.Equal(
		t,
		*ref.URL,
		fmt.Sprintf("https://api.github.com/repos/%s/%s/git/refs/heads/%s", orgName, repoName, branchName),
		"url",
	)
	assert.Equal(t, *ref.Object.Type, "commit", "object type")
	assert.Equal(t, *ref.Object.SHA, sha, "sha")
	assert.Equal(
		t,
		*ref.Object.URL,
		fmt.Sprintf("https://api.github.com/repos/%s/%s/git/commits/%s", orgName, repoName, sha),
		"url",
	)
}
