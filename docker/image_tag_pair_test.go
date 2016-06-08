package docker

import (
	"testing"

	"github.com/arschles/assert"
)

func TestCreateImageTagPairsFromTransform(t *testing.T) {
	const newRepo = "newrepo"
	xform := func(img Image) *Image {
		img.SetRepo(newRepo)
		return &img
	}
	imgs := createTestImages()
	pairs := CreateImageTagPairsFromTransform(imgs, xform)
	for i, pair := range pairs {
		assert.Equal(t, imgs[i].String(), pair.Source.String(), "source image")
		newImg := imgs[i]
		newImg.SetRepo(newRepo)
		assert.Equal(t, newImg.String(), pair.Target.String(), "target image")
	}
}
