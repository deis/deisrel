package docker

import (
	"strconv"
	"testing"
	"time"

	"github.com/arschles/assert"
)

func TestRetagNoImages(t *testing.T) {
	memCl := NewMemoryClient()
	pairs := []ImageTagPair{}
	pairCh, errCh, doneCh := RetagImages(memCl, pairs)

	select {
	case <-doneCh:
	case pair := <-pairCh:
		t.Errorf("a pair was returned when none were passed (%s)", pair)
	case err := <-errCh:
		t.Errorf("an error was returned when no pairs were passed (%s)", err)
	case <-time.After(1 * time.Millisecond):
		t.Errorf("the done channel wasn't closed immediately when no pairs were passed")
	}
}

func TestRetagImages(t *testing.T) {
	const newRepo = "newRepo"
	const newTag = "newTag"
	memCl := NewMemoryClient()
	images := createTestImages()
	pairs := make([]ImageTagPair, len(images))
	for i, img := range images {
		target := *img
		target.SetRepo(newRepo)
		target.SetTag("newTag")
		pairs[i] = ImageTagPair{Source: img, Target: &target}
	}
	pairCh, errCh, doneCh := RetagImages(memCl, pairs)
	retPairs := []ImageTagPair{}
	errs := []ErrTag{}
	for {
		done := false
		select {
		case <-doneCh:
			done = true
		case pair := <-pairCh:
			retPairs = append(retPairs, pair)
		case err := <-errCh:
			errs = append(errs, err)
		}
		if done {
			break
		}
	}

	assert.Equal(t, len(errs), 0, "number of errors")
	assert.Equal(t, len(retPairs), len(pairs), "number of retags returned")
	assert.Equal(t, len(memCl.Retags), len(pairs), "number of retags received")
	for i, pair := range pairs {
		assert.Equal(t, pair, retPairs[i], "retagged image "+strconv.Itoa(i))
		assert.Equal(t, pair, memCl.Retags[i], "retagged image "+strconv.Itoa(i))
	}
}
