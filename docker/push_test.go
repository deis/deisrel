package docker

import (
	"strconv"
	"testing"

	"github.com/arschles/assert"
)

func TestPushNoImages(t *testing.T) {
	memCl := NewMemoryClient()
	images := []*Image{}
	errs := PushImages(memCl, images)
	assert.Equal(t, len(errs), 0, "number of errors")
	assert.Equal(t, len(memCl.Pushes), 0, "number of pushes received")
}

func TestPushImages(t *testing.T) {
	memCl := NewMemoryClient()
	images := createTestImages()
	errs := PushImages(memCl, images)
	assert.Equal(t, len(errs), 0, "number of errors")
	assert.Equal(t, len(memCl.Pushes), len(images), "number of pushes received")
	for i, img := range images {
		assert.Equal(t, *img, *memCl.Pushes[i], "pushed image "+strconv.Itoa(i))
	}
}
