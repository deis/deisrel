package changelog

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/arschles/assert"
	"github.com/deis/deisrel/git"
	"github.com/deis/deisrel/testutil"
)

func TestMultiRepoVals(t *testing.T) {
	ts := testutil.NewTestServer()
	defer ts.Close()
	components := []string{"controller", "builder", "monitor", "logger"}
	for _, component := range components {
		ts.Mux.HandleFunc(fmt.Sprintf("/repos/deis/%s/compare/b...h", component), func(w http.ResponseWriter, r *http.Request) {
			if got := r.Method; got != "GET" {
				t.Errorf("Request method: %v, want GET", got)
			}
			fmt.Fprintf(w, `{
  		  "base_commit": {
  		    "sha": "s",
  		    "commit": {
  		      "author": { "name": "n" },
  		      "committer": { "name": "n" },
  		      "message": "m",
  		      "tree": { "sha": "t" }
  		    },
  		    "author": { "login": "n" },
  		    "committer": { "login": "l" },
  		    "parents": [ { "sha": "s" } ]
  		  },
  		  "status": "s",
  		  "ahead_by": 1,
  		  "behind_by": 2,
  		  "total_commits": 1,
  		  "commits": [
  		    {
  		      "sha": "abc1234567890",
  		      "commit": { "author": { "name": "n" }, "message": "feat(deisrel): new feature!" },
  		      "author": { "login": "l" },
  		      "committer": { "login": "l" },
  		      "parents": [ { "sha": "s" } ]
  		    },
  		    {
  		      "sha": "abc2345678901",
  		      "commit": { "author": { "name": "n" }, "message": "fix(deisrel): bugfix!" },
  		      "author": { "login": "l" },
  		      "committer": { "login": "l" },
  		      "parents": [ { "sha": "s" } ]
  		    },
  		    {
  		      "sha": "abc3456789012",
  		      "commit": { "author": { "name": "n" }, "message": "docs(deisrel): new docs!" },
  		      "author": { "login": "l" },
  		      "committer": { "login": "l" },
  		      "parents": [ { "sha": "s" } ]
  		    },
  		    {
  		      "sha": "abc4567890123",
  		      "commit": { "author": { "name": "n" }, "message": "doc(deisrel): new docs!" },
  		      "author": { "login": "l" },
  		      "committer": { "login": "l" },
  		      "parents": [ { "sha": "s" } ]
  		    },
  		    {
  		      "sha": "abc5678901234",
  		      "commit": { "author": { "name": "n" }, "message": "chore(deisrel): boring chore" },
  		      "author": { "login": "l" },
  		      "committer": { "login": "l" },
  		      "parents": [ { "sha": "s" } ]
  		    }
  		  ],
  		  "files": [ { "filename": "f" } ]
  		}`)
		})
	}

	reposAndSHAs := make([]git.RepoAndSha, len(components))
	for i, component := range components {
		reposAndSHAs[i] = git.RepoAndSha{
			Name: component,
			SHA:  "h",
		}
	}
	repoAndVals, err := MultiRepoVals(ts.Client, reposAndSHAs)
	assert.NoErr(t, err)
	assert.Equal(t, len(repoAndVals), len(reposAndSHAs), "number of returned repos")
}
