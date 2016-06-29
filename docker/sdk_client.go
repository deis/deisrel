package docker

import (
	"errors"
	"io/ioutil"

	dlib "github.com/fsouza/go-dockerclient"
)

var (
	errNotYetImplemented = errors.New("not yet implemented")
)

type sdkClient struct {
	cl *dlib.Client
}

// NewSDKClient creates a new Client that uses the SDK under the hood
func NewSDKClient(cl *dlib.Client) Client {
	return &sdkClient{cl: cl}
}

// Pull is the client interface implementation.
func (s *sdkClient) Pull(img *Image) error {
	return s.cl.PullImage(dlib.PullImageOptions{
		Repository:   img.FullWithoutTag(),
		Tag:          img.tag,
		OutputStream: ioutil.Discard,
	}, dlib.AuthConfiguration{})
}

// Push is the client interface implementation. It's currently not implemented and will return an
// error indicating that fact. See https://github.com/deis/deisrel/issues/108 for details.
func (s *sdkClient) Push(img *Image) error {
	return errNotYetImplemented
}

// Retag is the client interface implementation. It's currently not implemented and will return
// an error indicating that fact. See https://github.com/deis/deisrel/issues/108 for details.
func (s *sdkClient) Retag(src *Image, tar *Image) error {
	return errNotYetImplemented
}
