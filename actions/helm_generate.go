package actions

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/arschles/sys"
	"github.com/deis/deisrel/git"
	"github.com/google/go-github/github"
)

type nopWriteCloser struct {
	io.Writer
}

func (n nopWriteCloser) Close() error {
	return nil
}

// NopWriteCloser returns a WriteCloser with a no-op Close method wrapping
// the provided Writer w.
func NopWriteCloser(w io.Writer) io.WriteCloser {
	return nopWriteCloser{w}
}

func generateParams(fs sys.FS, whereTo string, paramsComponentMap genParamsComponentMap, helmChart helmChart) error {
	executeTo, err := executeToStaging(fs, filepath.Join(whereTo, "tpl"))
	if err != nil {
		log.Fatalf("Error creating staging file (%s)", err)
	}
	defer executeTo.Close()

	return helmChart.Template.Execute(executeTo, paramsComponentMap)
}

func executeToStaging(fs sys.FS, stagingSubDir string) (io.WriteCloser, error) {
	fs.MkdirAll(stagingSubDir, os.ModePerm)
	return fs.Create(filepath.Join(stagingSubDir, generateParamsFileName))
}

func getParamsComponentMap(ghClient *github.Client, defaultParamsComponentAttrs genParamsComponentAttrs, template *template.Template, ref string) genParamsComponentMap {
	paramsComponentMap := createParamsComponentMap()

	if template == generateParamsE2ETpl {
		repoNames = []string{"workflow-e2e"}
		componentNames = []string{"WorkflowE2E"}
	}
	for _, componentName := range componentNames {
		paramsComponentMap[componentName] = defaultParamsComponentAttrs
	}

	if defaultParamsComponentAttrs.Tag == "" {
		// gather latest sha for each repo via GitHub api
		reposAndShas, err := git.GetSHAs(ghClient, repoNames, git.ShortSHATransformNoErr, ref)
		if err != nil {
			log.Fatalf("No tag given and couldn't fetch sha from GitHub (%s)", err)
		} else if len(reposAndShas) < 1 {
			log.Fatalf("No tag given and no shas returned from GitHub for %s", defaultParamsComponentAttrs.Org)
		}

		// a given repo may track multiple components; update each component Tag accordingly
		for _, repoAndSha := range reposAndShas {
			repoComponentNames := repoToComponentNames[repoAndSha.Name]
			paramsComponentAttrs := defaultParamsComponentAttrs
			for _, componentName := range repoComponentNames {
				paramsComponentAttrs.Tag = "git-" + repoAndSha.SHA
				paramsComponentMap[componentName] = paramsComponentAttrs
			}
		}
	}

	return paramsComponentMap
}
