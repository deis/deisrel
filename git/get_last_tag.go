package git

import (
	"github.com/google/go-github/github"
)

// GetLastTag returns the latest tag for each of the repos in repos. It is not yet implemented
func GetLastTag(ghClient *github.Client, repos []string) (map[string]*github.RepositoryTag, error) {
	for _, repo := range repos {
		_, _, err := ghClient.Repositories.ListTags("deis", repo, nil)
		if err != nil {
			return make(map[string]*github.RepositoryTag), err
		}
	}
	return make(map[string]*github.RepositoryTag), nil
}
