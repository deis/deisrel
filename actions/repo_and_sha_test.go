package actions

import (
	"fmt"
	"testing"

	"github.com/arschles/assert"
)

func TestRepoAndShaString(t *testing.T) {
	ras := repoAndSha{repoName: "testRepo", sha: "testSHA"}
	assert.Equal(t, ras.String(), fmt.Sprintf("%s: %s", ras.repoName, ras.sha), "string representation of repo and sha")
}
