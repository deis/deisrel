package actions

import (
	"bytes"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"testing"

	"github.com/arschles/assert"
	"github.com/arschles/sys"
	"github.com/deis/deisrel/git"
	"github.com/deis/deisrel/testutil"
)

func TestGenParamsComponentMapWorkflowE2EEmptyTag(t *testing.T) {
	ts := testutil.NewTestServer()
	defer ts.Close()

	org := "deis"
	repo := "workflow-e2e"
	sha := "123abc456def"
	ref := "foo-ref"

	ts.Mux.HandleFunc(fmt.Sprintf("/repos/%s/%s/commits", org, repo), func(w http.ResponseWriter, r *http.Request) {
		if got := r.Method; got != "GET" {
			t.Errorf("Request method: %v, want GET", got)
		}
		if got := r.URL.RequestURI(); !strings.Contains(got, fmt.Sprintf("&sha=%s", ref)) {
			t.Errorf("Request URL (%v) did not include %s in 'sha' specifier", got, ref)
		}
		resp := `[
			{
			 "url": "",
			 "sha": "` + sha + `",
			 "html_url": "",
			 "comments_url": "",
			 "commit": {
				 "url": "",
				 "author": {
					 "name": "",
					 "email": "",
					 "date": "2011-04-14T16:00:49Z"
				 },
				 "committer": {
					 "name": "",
					 "email": "",
					 "date": "2011-04-14T16:00:49Z"
				 },
				 "message": "",
				 "tree": {
					 "url": "",
					 "sha": "` + sha + `"
				 },
				 "comment_count": 0,
				 "verification": {
					 "verified": true,
					 "reason": "valid",
					 "signature": "",
					 "payload": "tree ` + sha + `\n..."
				 }
			 },
			 "author": {
				 "login": "",
				 "id": 1,
				 "avatar_url": "",
				 "gravatar_id": "",
				 "url": "",
				 "html_url": "",
				 "followers_url": "",
				 "following_url": "",
				 "gists_url": "",
				 "starred_url": "",
				 "subscriptions_url": "",
				 "organizations_url": "",
				 "repos_url": "",
				 "events_url": "",
				 "received_events_url": "",
				 "type": "User",
				 "site_admin": false
			 },
			 "committer": {
				 "login": "",
				 "id": 1,
				 "avatar_url": "",
				 "gravatar_id": "",
				 "url": "",
				 "html_url": "",
				 "followers_url": "",
				 "following_url": "",
				 "gists_url": "",
				 "starred_url": "",
				 "subscriptions_url": "",
				 "organizations_url": "",
				 "repos_url": "",
				 "events_url": "",
				 "received_events_url": "",
				 "type": "User",
				 "site_admin": false
			 },
			 "parents": [
				 {
					 "url": "",
					 "sha": "` + sha + `"
				 }
			 ]
			}
		]`
		fmt.Fprintf(w, resp)
	})

	testGenParamsComponentMap(t, ts, "", fmt.Sprintf("git-%s", git.ShortSHATransformNoErr(sha)), ref, WorkflowE2EChart)
}

func TestGenParamsComponentMapWorkflow(t *testing.T) {
	ts := testutil.NewTestServer()
	defer ts.Close()
	testGenParamsComponentMap(t, ts, "tag", "tag", "ref", WorkflowChart)
}

func TestGenParamsComponentMapE2E(t *testing.T) {
	ts := testutil.NewTestServer()
	defer ts.Close()
	testGenParamsComponentMap(t, ts, "tag", "tag", "ref", WorkflowE2EChart)
}

func testGenParamsComponentMap(t *testing.T, ts *testutil.TestServer, inputTag, expectedTag, ref string, helmChart helmChart) {
	defaultParamsComponentAttrs := genParamsComponentAttrs{
		Org:        "org",
		PullPolicy: "pullPolicy",
		Tag:        inputTag,
	}
	got := getParamsComponentMap(ts.Client, defaultParamsComponentAttrs, helmChart.Template, ref)

	want := createParamsComponentMap()
	wantedPCA := defaultParamsComponentAttrs
	wantedPCA.Tag = expectedTag
	if helmChart.Template == generateParamsE2ETpl {
		want["WorkflowE2E"] = wantedPCA
	} else {
		for _, componentName := range componentNames {
			want[componentName] = wantedPCA
		}
	}

	assert.Equal(t, got, want, fmt.Sprintf("generated paramsComponentMap for helm chart %s", helmChart.Name))
}

func TestGenerateParamsStageWorkflow(t *testing.T) {
	fakeFS := sys.NewFakeFS()
	stagingDir := filepath.Join(defaultStagingPath, "foo")
	defaultParamsComponentAttrs := genParamsComponentAttrs{
		Org:        "org",
		Tag:        "",
		PullPolicy: "stayGangsta",
	}
	paramsComponentMap := createParamsComponentMap()
	for _, componentName := range componentNames {
		paramsComponentMap[componentName] = defaultParamsComponentAttrs
	}

	err := generateParams(fakeFS, stagingDir, paramsComponentMap, WorkflowChart)
	assert.NoErr(t, err)

	expectedStagedFilepath := filepath.Join(stagingDir, "tpl/generate_params.toml")
	// verify file exists in fakeFS
	_, err = fakeFS.ReadFile(expectedStagedFilepath)
	assert.NoErr(t, err)

	actualFileContents, err := fakeFS.ReadFile(expectedStagedFilepath)
	assert.NoErr(t, err)
	expectedFileContents := new(bytes.Buffer)
	err = generateParamsTpl.Execute(expectedFileContents, paramsComponentMap)
	assert.NoErr(t, err)

	assert.Equal(t, actualFileContents, expectedFileContents.Bytes(), "staged file contents")

	// make sure each component name from canonical list exists in actualFileContents
	for _, componentName := range componentNames {
		lowerCasedComponentName := strings.ToLower(componentName)
		if lowerCasedComponentName != "workflowe2e" {
			if lowerCasedComponentName == "workflowmanager" {
				lowerCasedComponentName = "workflowManager"
			}
			assert.True(t,
				strings.Contains(string(actualFileContents), lowerCasedComponentName),
				fmt.Sprintf("component: %s not found!", lowerCasedComponentName))
		}
	}
}

func TestGenerateParamsStageE2E(t *testing.T) {
	fakeFS := sys.NewFakeFS()
	stagingDir := filepath.Join(defaultStagingPath, "foo")
	defaultParamsComponentAttrs := genParamsComponentAttrs{
		Org:        "org",
		Tag:        "",
		PullPolicy: "stayGangsta",
	}
	paramsComponentMap := createParamsComponentMap()
	componentName := "WorkflowE2E"
	paramsComponentMap[componentName] = defaultParamsComponentAttrs

	err := generateParams(fakeFS, stagingDir, paramsComponentMap, WorkflowE2EChart)
	assert.NoErr(t, err)

	expectedStagedFilepath := filepath.Join(stagingDir, "tpl/generate_params.toml")
	// verify file exists in fakeFS
	_, err = fakeFS.ReadFile(expectedStagedFilepath)
	assert.NoErr(t, err)

	actualFileContents, err := fakeFS.ReadFile(expectedStagedFilepath)
	assert.NoErr(t, err)
	expectedFileContents := new(bytes.Buffer)
	err = generateParamsE2ETpl.Execute(expectedFileContents, paramsComponentMap)
	assert.NoErr(t, err)

	assert.Equal(t, actualFileContents, expectedFileContents.Bytes(), "staged file contents")

	assert.True(t,
		strings.Contains(string(actualFileContents), "e2e"),
		fmt.Sprintln("component: e2e not found!"))
}

func TestExecuteToStaging(t *testing.T) {
	fakeFS := sys.NewFakeFS()
	stagingDir := filepath.Join(defaultStagingPath, "foo")

	_, err := executeToStaging(fakeFS, stagingDir)
	assert.NoErr(t, err)

	// just verify dir was created on fakeFS
	_, err = fakeFS.ReadFile(stagingDir)
	assert.NoErr(t, err)
}
