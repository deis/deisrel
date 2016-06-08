package git

import (
	"fmt"
)

// RepoAndSha is the representation of a Git repo and a SHA in that repo
type RepoAndSha struct {
	Name string
	SHA  string
}

// ShortSHA returns the shortened SHA of r. If the SHA is already short, then returns just r.SHA
func (r RepoAndSha) ShortSHA() string {
	if len(r.SHA) < 8 {
		return r.SHA
	}
	return r.SHA[0:7]
}

// String is the fmt.Stringer interface implementation
func (r RepoAndSha) String() string {
	return fmt.Sprintf("%s: %s", r.Name, r.SHA)
}
