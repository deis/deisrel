package docker

import (
	"sync"
)

// PushImages pushes docker images, each of which is based on the items in images
func PushImages(cl Client, images []*Image) []error {
	var wg sync.WaitGroup
	errCh := make(chan error)
	doneCh := make(chan struct{})
	for _, image := range images {
		wg.Add(1)
		go func(img *Image) {
			defer wg.Done()
			if err := cl.Push(img); err != nil {
				errCh <- err
				return
			}
		}(image)
	}
	go func() {
		wg.Wait()
		close(doneCh)
	}()
	errs := []error{}
	for {
		select {
		case <-doneCh:
			return errs
		case err := <-errCh:
			errs = append(errs, err)
		}
	}
}
