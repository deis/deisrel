package docker

import (
	"testing"
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
