package git

import (
	"sync"

	"github.com/google/go-github/github"
)

// ReleaseInfo represents the data needed to create a single Github release
type ReleaseInfo struct {
	RepoName string
	Title    string
	Tag      string
	Body     string
}

// CreateRelease creates a single Github release
func CreateRelease(ghClient *github.Client, info ReleaseInfo) error {
	_, _, err := ghClient.Repositories.CreateRelease("deis", info.RepoName, &github.RepositoryRelease{
		TagName: github.String(info.Tag),
		Name:    github.String(info.Title),
		Body:    github.String(info.Body),
	})
	return err
}

// CreateReleases concurrently creates Github releases according to each element in infos.
// Returns a non-nil error if there was a problem creating _any_ Github release. If a non-nil
// error was returned, some releases may have already been created
func CreateReleases(ghClient *github.Client, infos []ReleaseInfo) error {
	var wg sync.WaitGroup
	errCh := make(chan error)
	doneCh := make(chan struct{})
	for _, info := range infos {
		wg.Add(1)
		go func(info ReleaseInfo) {
			defer wg.Done()
			if err := CreateRelease(ghClient, info); err != nil {
				select {
				case errCh <- err:
				case <-doneCh:
				}
				return
			}
		}(info)
	}
	go func() {
		wg.Wait()
		close(doneCh)
	}()
	for {
		select {
		case err := <-errCh:
			close(doneCh)
			return err
		case <-doneCh:
			return nil
		}
	}
}
