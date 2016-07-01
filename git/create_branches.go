package git

import (
	"fmt"
	"sync"

	"github.com/google/go-github/github"
)

// CreateBranches creates branches called branchName in all of the repos listed in reposAndSHAs,
// each from the given SHA in that element. Returns a slice of RepoAndSha, each with the repo name
// and sha given. The returned slice will not necessarily be in the same order as reposAndSHAs.
// Finally, this func returns a nil slice and a non-nil error if any create branch operation failed
func CreateBranches(ghClient *github.Client, branchName string, reposAndSHAs []RepoAndSha) ([]RepoAndSha, error) {
	var wg sync.WaitGroup
	rasCh := make(chan RepoAndSha)
	errCh := make(chan error)
	doneCh := make(chan struct{})
	for _, ras := range reposAndSHAs {
		refName := fmt.Sprintf("refs/heads/%s", branchName)
		wg.Add(1)
		go func(ras RepoAndSha, refName string) {
			defer wg.Done()
			if _, _, err := ghClient.Git.CreateRef(
				"deis",
				ras.Name,
				newBranchReference("deis", ras.Name, branchName, ras.SHA),
			); err != nil {
				errCh <- err
				return
			}
			rasCh <- ras
		}(ras, refName)
	}

	go func() {
		wg.Wait()
		close(doneCh)
	}()

	ret := []RepoAndSha{}
	for {
		select {
		case err := <-errCh:
			return nil, err
		case ras := <-rasCh:
			ret = append(ret, ras)
		case <-doneCh:
			return ret, nil
		}
	}
}
