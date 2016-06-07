package actions

import (
	"strings"
	"testing"

	"github.com/arschles/assert"
)

var (
	ras1 = repoAndSha{repoName: "testRepo1", sha: "testSHA"}
	ras2 = repoAndSha{repoName: "testRepo2", sha: "testSHA"}
)

func TestNewEmptyRepoAndShaList(t *testing.T) {
	empty := newEmptyRepoAndShaList()
	assert.NotNil(t, empty.repoAndShaMap, "repo and sha map")
}

func TestRepoAndShaListAdd(t *testing.T) {
	rasl := newEmptyRepoAndShaList()
	rasl.Add(ras1)
	assert.Equal(t, len(rasl.repoNames), 1, "length of repo names list")
	assert.Equal(t, rasl.repoNames[0], ras1.repoName, "repo name in list")
	assert.Equal(t, rasl.repoAndShaMap[ras1.repoName].sha, ras1.sha, "repo sha in map")
}

func TestRepoAndShaListSort(t *testing.T) {
	rasl := newEmptyRepoAndShaList()
	rasl.Add(ras1)
	rasl.Add(ras2)
	rasl.Sort()
	assert.Equal(t, len(rasl.repoNames), 2, "length of repo names list")
	assert.Equal(t, rasl.repoNames[0], ras1.repoName, "name of first repo")
	assert.Equal(t, rasl.repoNames[1], ras2.repoName, "name of second repo")
}

func TestRepoAndShaListString(t *testing.T) {
	rasl := newEmptyRepoAndShaList()
	rasl.Add(ras1)
	rasl.Add(ras2)
	str := rasl.String()
	spl := strings.Split(str, "\n")
	assert.Equal(t, len(spl), 2, "length of newline-split string")
	assert.Equal(t, spl[0], ras1.String(), "first repoAndSha string")
	assert.Equal(t, spl[1], ras2.String(), "second repoAndSha string")
}
