package git

import (
	"fmt"
	"testing"

	"github.com/arschles/assert"
)

func TestRepoAndShaString(t *testing.T) {
	ras := RepoAndSha{Name: "testRepo", SHA: "testSHA"}
	assert.Equal(t, ras.String(), fmt.Sprintf("%s: %s", ras.Name, ras.SHA), "string representation of repo and sha")
}
