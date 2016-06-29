package docker

import (
	"fmt"
	"strings"
	"sync"
)

const (
	// DeisCIDockerOrg represents the "deisci" docker organization on quay and the docker hub
	DeisCIDockerOrg = "deisci"
	// DeisDockerOrg represents the "deis" docker organization on quay and the docker hub
	DeisDockerOrg = "deis"
)

// ErrTag is the error returned when an image couldn't be retagged
type ErrTag struct {
	SourceImage *Image
	TargetImage *Image
	Err         error
}

// Error is the error interface implementation
func (e ErrTag) Error() string {
	return fmt.Sprintf(
		"tagging image %s to new tag %s (%s)",
		e.SourceImage.String(),
		e.TargetImage.String(),
		strings.TrimSpace(e.Err.Error()),
	)
}

// RetagImages concurrently retags all of the images in pairs.
// The first returned chan recieves on each image that successfully is retagged and
// the second on each image that can't be retagged.
// The total number of receives on the first and second channels will equal len(pairs),
// and after all receives happen, the 3rd chan will be closed
func RetagImages(cl Client, pairs []ImageTagPair) (<-chan ImageTagPair, <-chan ErrTag, <-chan struct{}) {
	succCh := make(chan ImageTagPair)
	errCh := make(chan ErrTag)
	doneCh := make(chan struct{})

	var wg sync.WaitGroup
	for _, pair := range pairs {
		wg.Add(1)
		go func(p ImageTagPair) {
			defer wg.Done()
			if err := cl.Retag(p.Source, p.Target); err != nil {
				errCh <- ErrTag{SourceImage: p.Source, TargetImage: p.Target, Err: err}
				return
			}
			succCh <- p
		}(pair)
	}

	go func() {
		wg.Wait()
		close(doneCh)
	}()

	return succCh, errCh, doneCh
}
