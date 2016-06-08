package docker

import (
	"strconv"
	"testing"
	"time"

	"github.com/arschles/assert"
)

func TestPullNoImages(t *testing.T) {
	memCl := NewMemoryClient()
	images := []*Image{}
	imgCh, errCh, doneCh := PullImages(memCl, images)
	select {
	case <-doneCh:
	case img := <-imgCh:
		t.Errorf("an image was returned when none were passed (%s)", img)
	case err := <-errCh:
		t.Errorf("an error was returned when no images were passed (%s)", err)
	case <-time.After(1 * time.Millisecond):
		t.Errorf("the done channel wasn't closed immediately when no images were passed")
	}
}

func TestPullImages(t *testing.T) {
	memCl := NewMemoryClient()
	images := createTestImages()
	imgCh, errCh, doneCh := PullImages(memCl, images)
	var retImages []Image
	var errs []ErrPullImage
	for {
		done := false
		select {
		case <-doneCh:
			done = true
		case img := <-imgCh:
			retImages = append(retImages, img)
		case err := <-errCh:
			errs = append(errs, err)
		}
		if done {
			break
		}
	}

	assert.Equal(t, len(retImages), len(images), "number of returned images")
	assert.Equal(t, len(retImages), len(memCl.Pulls), "number of pulled images")
	assert.Equal(t, len(errs), 0, "number of errors")
	for i, img := range retImages {
		assert.Equal(t, img, *images[i], "image "+strconv.Itoa(i))
		assert.Equal(t, img, *memCl.Pulls[i], "image "+strconv.Itoa(i))
	}
}
