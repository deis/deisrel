package registry

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/deisrel/testutil"
)

func TestCheckExistenceOnHubFound(t *testing.T) {
	ts := testutil.NewTestServer()
	defer ts.Close()

	testImageAndTag := ImageAndTag{
		Image: "myImage",
		Tag:   "myTag",
	}

	ts.Mux.HandleFunc(fmt.Sprintf("/v2/%s/tags/list", testImageAndTag.Image), func(w http.ResponseWriter, r *http.Request) {
		if got := r.Method; got != "GET" {
			t.Errorf("Request method: %v, want GET", got)
		}
		resp := `{
			"name": "` + testImageAndTag.Image + `",
			"tags": [
				"` + testImageAndTag.Tag + `"
			]
		}`
		fmt.Fprintf(w, resp)
	})

	fakeHubClient := NewFakeHubRegistry(t, ts)
	err := fakeHubClient.CheckExistence(testImageAndTag)
	assert.NoErr(t, err)
}

func TestCheckExistenceOnHubNotFound(t *testing.T) {
	ts := testutil.NewTestServer()
	defer ts.Close()

	testImageAndTag := ImageAndTag{
		Image: "myImage",
		Tag:   "myTag",
	}

	ts.Mux.HandleFunc(fmt.Sprintf("/v2/%s/tags/list", testImageAndTag.Image), func(w http.ResponseWriter, r *http.Request) {
		if got := r.Method; got != "GET" {
			t.Errorf("Request method: %v, want GET", got)
		}
		resp := `{
			"name": "` + testImageAndTag.Image + `",
			"tags": [
				"notMyTag"
			]
		}`
		fmt.Fprintf(w, resp)
	})

	fakeHubClient := NewFakeHubRegistry(t, ts)
	err := fakeHubClient.CheckExistence(testImageAndTag)

	expectedErr := ErrTagNotFound{
		imageAndTag: testImageAndTag,
		registry:    fakeHubClient.Client.URL,
	}
	assert.Err(t, expectedErr, err)
}
