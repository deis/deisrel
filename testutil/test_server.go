package testutil

import (
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/google/go-github/github"
)

// TestServer represents a test HTTP server along with a github.Client and Mux
type TestServer struct {
	Server *httptest.Server
	Client *github.Client
	Mux    *http.ServeMux
	Host   string
}

// NewTestServer sets up a test HTTP server along with a github.Client that is
// configured to talk to that test server.  Tests should register handlers on
// Mux which provide mock responses for the API method being tested.
func NewTestServer() *TestServer {
	// test server
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	// github client configured to use test server
	client := github.NewClient(nil)
	url, _ := url.Parse(server.URL)
	client.BaseURL = url
	client.UploadURL = url

	return &TestServer{
		Server: server,
		Client: client,
		Mux:    mux,
		Host:   url.Host,
	}
}

// Close closes the test HTTP server.
func (t *TestServer) Close() {
	t.Server.Close()
}
