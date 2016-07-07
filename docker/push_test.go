package docker

import (
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
	pushMap := make(map[string]struct{})
	for _, push := range memCl.Pushes {
		pushMap[push.String()] = struct{}{}
	}
	for i, img := range images {
		_, ok := pushMap[img.String()]
		assert.True(t, ok, "image %d (%s) wasn't pushed", i, img.String())
	}
}
