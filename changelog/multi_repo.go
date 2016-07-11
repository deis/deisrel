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

// MultiRepoVals concurrently fetches the changelogs for each repo in reposAndValues
// (using SingleRepoVals for each) and returns the RepoAndValues for each of the aforementioned
// repos. On any failure, returns a nil slice and appropriate error. The ordering of the
// returned slice is undefined.
func MultiRepoVals(client *github.Client, reposAndValues []RepoAndValues) ([]RepoAndValues, error) {
	var wg sync.WaitGroup
	ravCh := make(chan RepoAndValues)
	errCh := make(chan error)
	doneCh := make(chan struct{})

	for _, rav := range reposAndValues {
		wg.Add(1)
		go func(rav RepoAndValues) {
			defer wg.Done()
			if _, err := SingleRepoVals(client, rav.Values, rav.RepoAndSHA.SHA, rav.RepoAndSHA.Name, false); err != nil {
				select {
				case errCh <- err:
				case <-doneCh:
				}
				return
			}
			select {
			case ravCh <- RepoAndValues{RepoAndSHA: rav.RepoAndSHA, Values: rav.Values}:
			case <-doneCh:
			}
		}(rav)
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
			return nil, err
		case <-doneCh:
			return ret, nil
		}
	}
}
