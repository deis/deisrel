package docker

import (
	"fmt"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/deisrel/git"
)

func createTestImages() []*Image {
	return []*Image{
		&Image{registry: "quay.io", repo: "arschles", name: "myimg", tag: "abc"},
		&Image{registry: "registry1", repo: "repo1", name: "name1", tag: "tag1"},
		&Image{registry: "registry2", repo: "repo2", name: "name2", tag: "tag2"},
	}
}

func TestParseImageFromName(t *testing.T) {
	type testCase struct {
		Name     string
		Expected *Image
		Err      error
	}
	testCases := []testCase{
		testCase{
			Name:     "quay.io/arschles/myimg:abc",
			Expected: &Image{registry: "quay.io", repo: "arschles", name: "myimg", tag: "abc"},
			Err:      nil,
		},
	}
	for i, testCase := range testCases {
		parsed, err := ParseImageFromName(testCase.Name)
		if err != testCase.Err {
			t.Errorf("case %d expected error %s but got %s", i, testCase.Err, err)
			continue
		}
		if *parsed != *testCase.Expected {
			t.Errorf("case %d expected %s but got %s", i, *testCase.Expected, *parsed)
			continue
		}
	}
}

func TestParseImageFromRepoAndSha(t *testing.T) {
	registries := []string{"", DockerHubRegistry, "quay.io"}
	const org = "regorg"
	ras := git.RepoAndSha{Name: "testRepo", SHA: "testSHA"}
	imgs, err := ParseImageFromRepoAndSha(registries, org, ras)
	assert.NoErr(t, err)
	assert.Equal(t, len(imgs), len(registries), "number of returned images")
	for i, img := range imgs {
		expectedReg := registries[i]
		if expectedReg == DockerHubRegistry {
			expectedReg = ""
		}
		assert.Equal(t, img.registry, expectedReg, fmt.Sprintf("registry for image %d", i))
	}
}
