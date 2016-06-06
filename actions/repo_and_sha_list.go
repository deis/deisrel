package actions

import (
	"sort"
	"strings"
)

type repoAndShaList struct {
	repoNames     sort.StringSlice
	repoAndShaMap map[string]repoAndSha
}

func newEmptyRepoAndShaList() *repoAndShaList {
	return &repoAndShaList{repoNames: nil, repoAndShaMap: make(map[string]repoAndSha)}
}

func (r *repoAndShaList) Sort() {
	r.repoNames.Sort()
}

func (r *repoAndShaList) Add(ras repoAndSha) {
	r.repoNames = append(r.repoNames, ras.repoName)
	r.repoAndShaMap[ras.repoName] = ras
}

func (r repoAndShaList) String() string {
	strs := make([]string, len(r.repoNames))
	for i, repoName := range r.repoNames {
		strs[i] = r.repoAndShaMap[repoName].String()
	}
	return strings.Join(strs, "\n")
}
