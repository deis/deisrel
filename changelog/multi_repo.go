package changelog

import (
	"sync"

	"github.com/deis/deisrel/git"
	"github.com/google/go-github/github"
)

// RepoAndValues stores the RepoAndSha and its changelog values
type RepoAndValues struct {
	RepoAndSHA git.RepoAndSha
	Values     *Values
}

// MultiRepoVals concurrently fetches the changelogs for each repo in reposAndSHAs and returns
// the RepoAndValues for each of the aforementioned repos. On any failure, returns a nil slice
// and appropriate error. The ordering of the returned slice is undefined.
func MultiRepoVals(client *github.Client, reposAndSHAs []git.RepoAndSha) ([]RepoAndValues, error) {
	var wg sync.WaitGroup
	ravCh := make(chan RepoAndValues)
	errCh := make(chan error)
	doneCh := make(chan struct{})

	for _, ras := range reposAndSHAs {
		wg.Add(1)
		go func(ras git.RepoAndSha) {
			defer wg.Done()
			vals := new(Values)
			if _, err := SingleRepoVals(client, vals, ras.SHA, ras.Name, false); err != nil {
				select {
				case errCh <- err:
				case <-doneCh:
				}
				return
			}
			select {
			case ravCh <- RepoAndValues{RepoAndSHA: ras, Values: vals}:
			case <-doneCh:
			}
		}(ras)
	}
	go func() {
		wg.Wait()
		close(doneCh)
	}()

	var ret []RepoAndValues
	for {
		select {
		case rav := <-ravCh:
			ret = append(ret, rav)
		case err := <-errCh:
			close(doneCh)
			return nil, err
		case <-doneCh:
			return ret, nil
		}
	}
}
