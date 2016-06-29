package docker

import (
	"fmt"
	"strings"
	"sync"
)

// ErrPullImage is the error returned when an image couldn't be pulled
type ErrPullImage struct {
	Img *Image
	Err error
}

// Error is the error interface implementation
func (e ErrPullImage) Error() string {
	return fmt.Sprintf("pulling image %s (%s)", *e.Img, strings.TrimSpace(e.Err.Error()))
}

// PullImages pulls each image in images concurrently. The first returned channel receives on
// each image successfully pulled, and the second receives on each image that failed to pull
// for any reason. The total recieves across both channels will equal len(images), and the
// third channel will be closed only after all of those receives occur.
func PullImages(cl Client, images []*Image) (<-chan Image, <-chan ErrPullImage, <-chan struct{}) {
	succCh := make(chan Image)
	errCh := make(chan ErrPullImage)
	doneCh := make(chan struct{})
	var wg sync.WaitGroup
	for _, image := range images {
		wg.Add(1)
		go func(img *Image) {
			defer wg.Done()
			if err := cl.Pull(img); err != nil {
				errCh <- ErrPullImage{Img: img, Err: err}
				return
			}
			succCh <- *img
		}(image)
	}
	go func() {
		wg.Wait()
		close(doneCh)
	}()

	return succCh, errCh, doneCh
}
