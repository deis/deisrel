package git

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/deisrel/testutil"
)

func TestCreateBranches(t *testing.T) {
	const (
		orgName  = "deis"
		repoName = "controller"
		branch   = "testbranch"
		commit   = "aa218f56b14c9653891f9e74264a383fa43fefbd"
	)
	ts := testutil.NewTestServer()
	defer ts.Close()
	ts.Mux.HandleFunc(fmt.Sprintf("/repos/%s/%s/git/refs", orgName, repoName), func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, "POST", "request method")
		ret := fmt.Sprintf(`
    {
      "ref": "refs/heads/%s",
      "url": "https://api.github.com/repos/%s/%s/git/refs/heads/%s",
      "object": {
        "type": "commit",
        "sha": "%s",
        "url": "https://api.github.com/repos/%s/%s/git/commits/%s"
      }
    }`, branch, orgName, repoName, branch, commit, orgName, repoName, commit)
		fmt.Fprintf(w, ret)
	})

	rasl := []RepoAndSha{
		RepoAndSha{Name: repoName, SHA: "master"},
	}
	retRasl, err := CreateBranches(ts.Client, branch, rasl)
	assert.NoErr(t, err)
	assert.Equal(t, len(retRasl), len(rasl), "number of returned RepoAndSha structs")
	assert.Equal(t, retRasl[0], rasl[0], "returned RepoAndSha")
}
