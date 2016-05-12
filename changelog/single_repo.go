package changelog

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/deis/deisrel/git"
	"github.com/google/go-github/github"
)

// SingleRepoVals generates a changelog entry from vals.OldRelease to sha. It returns the commits that were unparseable (and had to be skipped) or any error encountered during the process. On a nil error, vals is filled in with all of the sorted changelog entries
func SingleRepoVals(client *github.Client, vals *Values, sha, name string) ([]string, error) {
	var skippedCommits []string
	commitCompare, resp, err := client.Repositories.CompareCommits("deis", name, vals.OldRelease, sha)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			return nil, errTagNotFoundForRepo{repoName: name, tagName: vals.OldRelease}
		}
		return nil, errCouldNotCompareCommits{old: vals.OldRelease, new: sha, err: err}
	}
	for _, commit := range commitCompare.Commits {
		if commit.Commit.Message == nil {
			continue
		}
		if commit.SHA == nil {
			continue
		}
		commitMessage := strings.Split(*commit.Commit.Message, "\n")[0]
		shortSHA, err := git.ShortSHATransform(*commit.SHA)
		if err != nil {
			return nil, err
		}
		focus := commitFocus(*commit.Commit.Message)
		title := commitTitle(*commit.Commit.Message)
		changelogMessage := fmt.Sprintf("%s %s: %s", shortSHA, focus, title)
		if strings.HasPrefix(commitMessage, "feat(") {
			vals.Features = append(vals.Features, changelogMessage)
		} else if strings.HasPrefix(commitMessage, "fix(") {
			vals.Fixes = append(vals.Fixes, changelogMessage)
		} else if strings.HasPrefix(commitMessage, "docs(") || strings.HasPrefix(commitMessage, "doc(") {
			vals.Documentation = append(vals.Documentation, changelogMessage)
		} else if strings.HasPrefix(commitMessage, "chore(") {
			vals.Maintenance = append(vals.Maintenance, changelogMessage)
		} else {
			skippedCommits = append(skippedCommits, *commit.SHA)
		}
	}
	return skippedCommits, nil
}
