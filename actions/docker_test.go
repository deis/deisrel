package actions

import (
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/deisrel/testutil"

	reg "github.com/deis/deisrel/registry"
)

// TODO: place in 'docker' package, https://github.com/deis/deisrel/issues/60
func TestDockerCheckTags(t *testing.T) {
	ts := testutil.NewTestServer()
	defer ts.Close()

	checkExistenceFound := func(imgTag reg.ImageAndTag) error { return nil }
	fakeQuay := reg.NewFakeQuayRegistry(ts)
	fakeQuay.ExistenceChecker = checkExistenceFound
	fakeHub := reg.NewFakeHubRegistry(t, ts)
	fakeHub.ExistenceChecker = checkExistenceFound

	var repoAndShas []repoAndSha
	for _, repo := range repoNames { // use repoNames list from common.go
		repoAndShas = append(repoAndShas,
			repoAndSha{repoName: repo, sha: "1234abcd"})
	}

	foundImgTags, errs := dockerCheckTags(ts.Client, fakeQuay, fakeHub, repoAndShas, "ref")
	assert.Equal(t, errs, []error{}, "errs")
	// check against len of componentNames (common.go) as multiple components to a repo
	assert.Equal(t, len(foundImgTags), len(componentNames), "foundImgTags length")
}
