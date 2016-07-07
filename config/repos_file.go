package config

import (
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const (
	// DefaultRepoNamesFile is the default name of the repositories file
	DefaultRepoNamesFile = "deisrel-repo-component-names.yaml"
)

// RepoToComponentNamesFile is the struct representation of the yaml file that holds the mapping
// from each repository to the components that it holds
type RepoToComponentNamesFile struct {
	RepoToComponentNames map[string][]string `yaml:"repoToComponentNames"`
}

// DecodeRepoToComponentNames decodes yamlFile into a RepoToComponentNamesFile. If there was
// a decode error, returns nil and the appropriate error
func DecodeRepoToComponentNames(yamlFile io.Reader) (*RepoToComponentNamesFile, error) {
	all, err := ioutil.ReadAll(yamlFile)
	if err != nil {
		return nil, err
	}
	ret := new(RepoToComponentNamesFile)

	if err := yaml.Unmarshal(all, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// GetRepoNames returns all of the repo names in an undefined order
func (r *RepoToComponentNamesFile) GetRepoNames() []string {
	ret := make([]string, len(r.RepoToComponentNames))
	i := 0
	for repoName := range r.RepoToComponentNames {
		ret[i] = repoName
		i++
	}
	return ret
}

// GetComponentNames returns all of the component names in an undefined order
func (r *RepoToComponentNamesFile) GetComponentNames() []string {
	var ret []string
	for _, compNameSlice := range r.RepoToComponentNames {
		for _, compName := range compNameSlice {
			ret = append(ret, compName)
		}
	}
	return ret
}
