package actions

import (
	"fmt"
)

type repoAndSha struct {
	repoName string
	sha      string
}

func (r repoAndSha) String() string {
	return fmt.Sprintf("%s: %s", r.repoName, r.sha)
}
