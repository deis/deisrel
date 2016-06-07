package registry

import (
	"github.com/coreos/go-quay/quay"
	tag "github.com/coreos/go-quay/quay/tag"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// QuayRegistry implements the Registry interface.  It includes a Client, Auth
// and helper methods for interacting with the quay.io api
type QuayRegistry struct {
	Client *quay.Client
	Auth   QuayAuth
}

// NewQuayRegistry returns a pointer to a new QuayRegistry
func NewQuayRegistry(c *quay.Client, auth QuayAuth) *QuayRegistry {
	return &QuayRegistry{
		Client: c,
		Auth:   auth,
	}
}

// QuayAuth is an interface providing AuthenticateRequest
type QuayAuth interface {
	AuthenticateRequest(runtime.ClientRequest, strfmt.Registry) error
}

// NewQuayAuth returns a QuayAuth instance
func NewQuayAuth(authToken string) QuayAuth {
	return &realQuayAuth{AuthToken: authToken}
}

type realQuayAuth struct {
	AuthToken string
}

func (rqa *realQuayAuth) AuthenticateRequest(cr runtime.ClientRequest, reg strfmt.Registry) error {
	return cr.SetHeaderParam("Authorization", rqa.AuthToken)
}

// CheckExistence is the ExistenceChecker for QuayRegistry
func (q *QuayRegistry) CheckExistence(imageAndTag ImageAndTag) error {
	quayListTagImagesParams := &tag.ListTagImagesParams{
		Repository: imageAndTag.Image,
		Tag:        imageAndTag.Tag,
	}

	_, err := q.Client.Tag.ListTagImages(quayListTagImagesParams, q.Auth)
	switch err {
	case err.(*tag.ListTagImagesNotFound):
		return ErrTagNotFound{
			imageAndTag: imageAndTag,
			registry:    "quay.io",
		}
	default:
		return err
	}
}
