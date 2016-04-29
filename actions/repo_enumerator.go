package actions

import (
	"sync"

	"github.com/google/go-github/github"
)

func enumerateRepos(
	ghClient *github.Client,
	repos []string,
	mapper func(*github.Client, int, string) error,
) []error {

	errCh := make(chan error)
	var wg sync.WaitGroup
	for i, repo := range repos {
		wg.Add(1)
		ech := make(chan error)
		go func(repoNum int, repoName string) {
			ech <- mapper(ghClient, repoNum, repoName)
		}(i, repo)
		go func() {
			defer wg.Done()
			errCh <- <-ech
		}()
	}
	go func() {
		wg.Wait()
		close(errCh)
	}()
	errs := []error{}
	for err := range errCh {
		errs = append(errs, err)
	}
	return errs
}
