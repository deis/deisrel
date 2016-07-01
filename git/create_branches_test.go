package git

import (
	"encoding/json"
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
		retRef := newBranchReference(orgName, repoName, branch, commit)
		assert.NoErr(t, json.NewEncoder(w).Encode(retRef))
	})

	rasl := []RepoAndSha{
		RepoAndSha{Name: repoName, SHA: "master"},
	}
	retRasl, err := CreateBranches(ts.Client, branch, rasl)
	assert.NoErr(t, err)
	assert.Equal(t, len(retRasl), len(rasl), "number of returned RepoAndSha structs")
	assert.Equal(t, retRasl[0], rasl[0], "returned RepoAndSha")
}
