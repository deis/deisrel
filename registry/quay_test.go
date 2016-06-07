package registry

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/deisrel/testutil"
)

func TestCheckExistenceOnQuayFound(t *testing.T) {
	ts := testutil.NewTestServer()
	defer ts.Close()

	testImageAndTag := ImageAndTag{
		Image: "myImage",
		Tag:   "myTag",
	}

	ts.Mux.HandleFunc(fmt.Sprintf("/api/v1/repository/%s/tag/%s/images", testImageAndTag.Image, testImageAndTag.Tag), func(w http.ResponseWriter, r *http.Request) {
		if got := r.Method; got != "GET" {
			t.Errorf("Request method: %v, want GET", got)
		}
	})

	fakeQuayRegistry := NewFakeQuayRegistry(ts)
	err := fakeQuayRegistry.CheckExistence(testImageAndTag)
	assert.NoErr(t, err)
}

func TestCheckExistenceOnQuayNotFound(t *testing.T) {
	ts := testutil.NewTestServer()
	defer ts.Close()

	testImageAndTag := ImageAndTag{
		Image: "myImage",
		Tag:   "myTag",
	}

	ts.Mux.HandleFunc(fmt.Sprintf("/api/v1/repository/%s/tag/%s/images", testImageAndTag.Image, testImageAndTag.Tag), func(w http.ResponseWriter, r *http.Request) {
		if got := r.Method; got != "GET" {
			t.Errorf("Request method: %v, want GET", got)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		resp := `{
		  "status": 404,
		  "error_message": "Not Found",
		  "title": "not_found",
		  "error_type": "not_found",
		  "detail": "Not Found",
		  "type": "https://quay.io/api/v1/error/not_found"
		}`

		fmt.Fprintf(w, resp)
	})

	fakeQuayRegistry := NewFakeQuayRegistry(ts)
	err := fakeQuayRegistry.CheckExistence(testImageAndTag)

	expectedErr := ErrTagNotFound{
		imageAndTag: testImageAndTag,
		registry:    "quay.io",
	}
	assert.Err(t, expectedErr, err)
}
