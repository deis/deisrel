package changelog

import (
	"fmt"
)

type errTagNotFoundForRepo struct {
	repoName string
	tagName  string
}

func (e errTagNotFoundForRepo) Error() string {
	return fmt.Sprintf("tag %s not found for repo %s", e.tagName, e.repoName)
}

type errCouldNotCompareCommits struct {
	old string
	new string
	err error
}

func (e errCouldNotCompareCommits) Error() string {
	return fmt.Sprintf("could not compare commits %s - %s (%s)", e.old, e.new, e.err)
}
