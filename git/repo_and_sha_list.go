package git

import (
	"sort"
	"strings"
)

// RepoAndShaList is a sortable list of RepoAndSha structs, sortable on the repo name
type RepoAndShaList struct {
	repoNames     sort.StringSlice
	repoAndShaMap map[string]RepoAndSha
}

// NewEmptyRepoAndShaList creates a new RepoAndShaList with no items in it
func NewEmptyRepoAndShaList() *RepoAndShaList {
	return &RepoAndShaList{repoNames: nil, repoAndShaMap: make(map[string]RepoAndSha)}
}

// NewRepoAndShaListFromSlice returns a new RepoAndShaList with all items in rasl added (in order)
func NewRepoAndShaListFromSlice(rasl []RepoAndSha) *RepoAndShaList {
	ret := NewEmptyRepoAndShaList()
	for _, ras := range rasl {
		ret.Add(ras)
	}
	return ret
}

// Sort sorts the internal repo list by name
func (r *RepoAndShaList) Sort() {
	r.repoNames.Sort()
}

// Slice returns a slice of the RepoAndShas that are stored internally, in the same order
func (r *RepoAndShaList) Slice() []RepoAndSha {
	ret := make([]RepoAndSha, len(r.repoAndShaMap))
	for i, name := range r.repoNames {
		ret[i] = r.repoAndShaMap[name]
	}
	return ret
}

// Add adds a new RepoAndSha to the internal list
func (r *RepoAndShaList) Add(ras RepoAndSha) {
	r.repoNames = append(r.repoNames, ras.Name)
	r.repoAndShaMap[ras.Name] = ras
}

// String is the fmt.Stringer interface implementation
func (r RepoAndShaList) String() string {
	strs := make([]string, len(r.repoNames))
	for i, repoName := range r.repoNames {
		strs[i] = r.repoAndShaMap[repoName].String()
	}
	return strings.Join(strs, "\n")
}
