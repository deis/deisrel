package actions

import (
	"fmt"
	"sync"

	"github.com/google/go-github/github"

	reg "github.com/deis/deisrel/registry"
)

const (
	ciRepoOrg    = "deisci"
	prodRepoOrg  = "deis"
	monitorImage = "monitor"
)

// TODO: place in 'docker' package, https://github.com/deis/deisrel/issues/60
func dockerCheckTags(ghClient *github.Client, quay reg.Registry, hub reg.Registry, repoAndShas []repoAndSha, ref string) ([]string, []error) {
	var wg sync.WaitGroup
	done := make(chan bool)
	resultsCh := make(chan string)
	errCh := make(chan error)

	for _, rs := range repoAndShas {
		wg.Add(1)
		go func(rs repoAndSha) {
			defer wg.Done()
			var imgTags []reg.ImageAndTag

			// a given GitHub repo may contain multiple components,
			// each publishing images to their own quay repos
			for _, componentName := range repoToComponentNames[rs.repoName] {
				imgTags = append(imgTags,
					reg.ImageAndTag{
						Image: fmt.Sprintf("%s/%s", ciRepoOrg, componentToImageName[componentName]),
						Tag:   imageTagTransform(rs.sha),
					})
			}
			for _, imgTag := range imgTags {
				if err := quay.CheckExistence(imgTag); err != nil {
					if err != nil {
						errCh <- err
						return
					}
				}
				if err := hub.CheckExistence(imgTag); err != nil {
					if err != nil {
						errCh <- err
						return
					}
				}
				resultsCh <- fmt.Sprintf("%s:%s", imgTag.Image, imgTag.Tag)
			}
		}(rs)
	}
	go func() {
		// wait for all fetches from github to be complete before returning
		wg.Wait()
		close(done)
	}()

	results := []string{}
	errs := []error{}
	for {
		select {
		case <-done:
			return results, errs
		case result := <-resultsCh:
			results = append(results, result)
		case err := <-errCh:
			errs = append(errs, err)
		}
	}
}
