package docker

// ImageTagPair represents a source image that will be re-tagged to a target image
type ImageTagPair struct {
	Source *Image
	Target *Image
}

// CreateImageTagPairsFromTransform returns a slice of ImageTagPairs where each
// source is images[i] and each target is transform(images[i])
func CreateImageTagPairsFromTransform(images []*Image, transform func(Image) *Image) []ImageTagPair {
	ret := make([]ImageTagPair, len(images))
	for i, img := range images {
		ret[i] = ImageTagPair{Source: img, Target: transform(*img)}
	}
	return ret
}
