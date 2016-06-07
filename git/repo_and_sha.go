package git

import (
	"fmt"
)

// RepoAndSha is the representation of a Git repo and a SHA in that repo
type RepoAndSha struct {
	Name string
	SHA  string
}

// String is the fmt.Stringer interface implementation
func (r RepoAndSha) String() string {
	return fmt.Sprintf("%s: %s", r.Name, r.SHA)
}
