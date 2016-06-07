package registry

import (
	"log"
	"net/http"
	"testing"

	"github.com/coreos/go-quay/quay"
	"github.com/deis/deisrel/testutil"

	httptransport "github.com/go-openapi/runtime/client"
	hub "github.com/heroku/docker-registry-client/registry"
)

// FakeQuayRegistry extends Registry particular to the quay.io api
type FakeQuayRegistry struct {
	Client           *quay.Client
	Auth             QuayAuth
	ExistenceChecker ExistenceChecker
}

// CheckExistence is the implentation of said function for a FakeQuayRegistry
func (fqr *FakeQuayRegistry) CheckExistence(imgTag ImageAndTag) error {
	if fqr.ExistenceChecker != nil {
		return fqr.ExistenceChecker(imgTag)
	}
	return NewQuayRegistry(fqr.Client, fqr.Auth).CheckExistence(imgTag)
}

// NewFakeQuayRegistry returns a FakeQuayRegistry using the TestServer ts provided
func NewFakeQuayRegistry(ts *testutil.TestServer) *FakeQuayRegistry {
	return &FakeQuayRegistry{
		Client:           quay.New(httptransport.New(ts.Host, "/", []string{"http"}), nil),
		Auth:             nil,
		ExistenceChecker: nil,
	}
}

// FakeHubRegistry extends Registry particular to the DockerHub api
type FakeHubRegistry struct {
	Client           *hub.Registry
	ExistenceChecker ExistenceChecker
}

// CheckExistence is the implentation of said function for a FakeQuayRegistry
func (fhr *FakeHubRegistry) CheckExistence(imgTag ImageAndTag) error {
	if fhr.ExistenceChecker != nil {
		return fhr.ExistenceChecker(imgTag)
	}
	return NewHubRegistry(fhr.Client).CheckExistence(imgTag)
}

// NewFakeHubRegistry returns a FakeHubRegistry using the TestServer ts provided
func NewFakeHubRegistry(t *testing.T, ts *testutil.TestServer) *FakeHubRegistry {
	// handle automatic registry ping
	ts.Mux.HandleFunc("/v2/", func(w http.ResponseWriter, r *http.Request) {})

	hub, err := hub.NewInsecure(ts.Server.URL, "username", "password")
	if err != nil {
		log.Fatalf("Error creating new hub (%s)", err)
	}

	return &FakeHubRegistry{Client: hub, ExistenceChecker: nil}
}
