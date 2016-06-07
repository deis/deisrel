package git

import (
	"fmt"
	"strings"

	"github.com/google/go-github/github"
)

// CreateTags tags each repository at the corresponding SHA with the tag supplied
func CreateTags(client *github.Client, repos []RepoAndSha, tag string) error {
	for _, repo := range repos {
		ref := &github.Reference{
			Ref: github.String("refs/tags/" + tag),
			Object: &github.GitObject{
				SHA: github.String(repo.SHA),
			},
		}
		_, _, err := client.Git.CreateRef("deis", repo.Name, ref)
		// GitHub returns HTTP 422 Unprocessable Entity when a field is invalid,
		// such as when a reference already exists or the sha does not exist
		// https://developer.github.com/v3/#client-errors
		if err != nil && !strings.Contains(err.Error(), "Reference already exists") {
			return err
		}
		fmt.Printf("%s: created tag %s\n", repo.Name, tag)
	}
	return nil
}
