package git

import (
	"fmt"
	"sync"

	"github.com/google/go-github/github"
)

// GetSHAs returns a slice of the latest RepoAndSha for each repository, using ghClient.
// It transforms each SHA that it finds using the transform function. It uses the ref parameter
// to determine the SHA or branch from which to list commits
func GetSHAs(ghClient *github.Client, repos []string, transform func(string) string, ref string) ([]RepoAndSha, error) {
	outCh := make(chan RepoAndSha)
	errCh := make(chan error)
	doneCh := make(chan struct{})
	var wg sync.WaitGroup
	for _, repo := range repos {
		wg.Add(1)
		ch := make(chan RepoAndSha)
		ech := make(chan error)
		go func(repo string) {
			commitsListOpts := &github.CommitsListOptions{
				SHA: ref,
				ListOptions: github.ListOptions{
					Page:    1,
					PerPage: 1,
				},
			}
			repoCommits, _, err := ghClient.Repositories.ListCommits("deis", repo, commitsListOpts)
			if err != nil {
				ech <- fmt.Errorf("Error listing commits for repo %s (%s)", repo, err)
				return
			}
			if len(repoCommits) < 1 {
				ech <- fmt.Errorf("No commits found for repo %s", repo)
				return
			}
			repoCommit := repoCommits[0]
			sha := transform(*repoCommit.SHA)
			ch <- RepoAndSha{Name: repo, SHA: sha}
		}(repo)
		go func() {
			defer wg.Done()
			select {
			case e := <-ech:
				errCh <- e
			case o := <-ch:
				outCh <- o
			}
		}()
	}
	go func() {
		wg.Wait()
		close(doneCh)
	}()

	ret := []RepoAndSha{}
	for {
		select {
		case <-doneCh:
			return ret, nil
		case str := <-outCh:
			ret = append(ret, str)
		case err := <-errCh:
			return nil, err
		}
	}
}
