package docker

import (
	"fmt"
	"strings"

	"github.com/deis/deisrel/git"
)

const (
	DockerHubRegistry = "index.docker.io"
)

// ErrInvalidImageName is the error returned when a func couldn't parse a string into an
// Image struct
type ErrInvalidImageName struct {
	Str string
}

// Error is the error interface implementation
func (e ErrInvalidImageName) Error() string {
	return fmt.Sprintf("%s is an invalid image name", e.Str)
}

// Image represents a single image name, including all information about its registry,
// repository and tag
type Image struct {
	registry string
	repo     string
	name     string
	tag      string
}

// ParseImageFromName parses a raw image name string into an Image
func ParseImageFromName(name string) (*Image, error) {
	spl := strings.Split(name, "/")
	splLast := strings.Split(spl[len(spl)-1], ":")
	tag := "latest"
	if len(splLast) > 1 {
		tag = splLast[1]
		spl[len(spl)-1] = splLast[0]
	}
	if len(spl) == 1 {
		// dockerhub trusted image
		return &Image{
			registry: "",
			repo:     "",
			name:     spl[0],
			tag:      tag,
		}, nil
	} else if len(spl) == 2 {
		// dockerhub image
		return &Image{
			registry: "",
			repo:     spl[0],
			name:     spl[1],
			tag:      tag,
		}, nil
	} else if len(spl) == 3 {
		// non-dockerhub image
		return &Image{
			registry: spl[0],
			repo:     spl[1],
			name:     spl[2],
			tag:      tag,
		}, nil
	}
	return nil, ErrInvalidImageName{Str: name}
}

// ParseImageFromRepoAndSha attempts to convert ras into a set of docker images, each of which
// points to a registry in registries and uses dockerRegistryOrg. Any registry that matches
// DockerHubRegistry will be converted to empty
func ParseImageFromRepoAndSha(
	dockerRegistries []string,
	dockerRegistryOrg string,
	ras git.RepoAndSha,
) ([]*Image, error) {
	ret := make([]*Image, len(dockerRegistries))
	for i, reg := range dockerRegistries {
		str := fmt.Sprintf("%s/%s:git-%s", dockerRegistryOrg, ras.Name, ras.ShortSHA())
		if reg == DockerHubRegistry {
			reg = ""
		}
		if reg != "" {
			// prepend the registry if it's non-empty
			str = fmt.Sprintf("%s/%s", reg, str)
		}
		img, err := ParseImageFromName(str)
		if err != nil {
			return nil, err
		}
		ret[i] = img
	}
	return ret, nil
}

// ParseImagesFromRepoAndShaList returns a slice of parsed Images in the same order as they
// appear in rasl.Slice(). Returns an empty slice and a non-nil error if any one of the
// git.RepoAndShas couldn't be parsed
func ParseImagesFromRepoAndShaList(
	dockerRegistries []string,
	dockerRegistryOrg string,
	rasl *git.RepoAndShaList,
) ([]*Image, error) {
	raslSlice := rasl.Slice()
	ret := []*Image{}
	for _, ras := range raslSlice {
		imgs, err := ParseImageFromRepoAndSha(dockerRegistries, dockerRegistryOrg, ras)
		if err != nil {
			return nil, err
		}
		for _, img := range imgs {
			ret = append(ret, img)
		}
	}
	return ret, nil
}

// FullWithoutTag returns the full image name without its tag
func (i Image) FullWithoutTag() string {
	return strings.Split(i.String(), ":")[0]
}

// String is the fmt.Stringer interface implementation. It returns the full image name and its tag
func (i Image) String() string {
	if i.registry != "" {
		return fmt.Sprintf("%s/%s/%s:%s", i.registry, i.repo, i.name, i.tag)
	} else if i.repo != "" {
		return fmt.Sprintf("%s/%s:%s", i.repo, i.name, i.tag)
	}
	return fmt.Sprintf("%s:%s", i.name, i.tag)
}

// SetRepo sets the image repository of i to repo
func (i *Image) SetRepo(repo string) {
	i.repo = repo
}

// SetTag sets the image tag of i to tag
func (i *Image) SetTag(tag string) {
	i.tag = tag
}
