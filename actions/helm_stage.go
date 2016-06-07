package actions

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/arschles/sys"
	"github.com/codegangsta/cli"
	"github.com/deis/deisrel/git"
	"github.com/google/go-github/github"
)

var (
	ourFS = sys.RealFS()
	ourFP = sys.RealFP()
)

func helmStage(ghClient *github.Client, c *cli.Context, helmChart helmChart) {
	var opt github.RepositoryContentGetOptions
	opt.Ref = c.GlobalString(RefFlag)
	org := c.GlobalString(GHOrgFlag)
	repo := c.GlobalString(RepoFlag)
	stagingDir := c.GlobalString(StagingDirFlag)

	if stagingDir == "" {
		stagingDir = filepath.Join(defaultStagingPath, helmChart.Name)
	}
	// create stagingDir and 'tpl' subdir for staging files
	if err := createDir(ourFS, filepath.Join(stagingDir, "tpl")); err != nil {
		log.Fatalf("Error creating dir %s (%s)", filepath.Join(stagingDir, "tpl"), err)
	}

	// gather helmChart.Files from GitHub needing release string updates
	ghFiles, err := downloadFiles(ghClient, org, repo, &opt, helmChart)
	if err != nil {
		log.Fatalf("Error downloading contents of %v (%s)", helmChart.Files, err)
	}
	stageFiles(ourFS, ghFiles, stagingDir)

	if err := updateFilesWithRelease(ourFP, ourFS, deisRelease, stagingDir); err != nil {
		log.Fatalf("Error updating files with release '%s' (%s)", deisRelease.Short, err)
	}

	// stage 'tpl/generate_params.toml' with latest git shas for each component
	defaultParamsComponentAttrs := genParamsComponentAttrs{
		Org:        c.GlobalString(OrgFlag),
		PullPolicy: c.GlobalString(PullPolicyFlag),
		Tag:        c.GlobalString(TagFlag),
	}
	paramsComponentMap := getParamsComponentMap(ghClient, defaultParamsComponentAttrs, helmChart.Template, c.GlobalString(RefFlag))
	generateParams(ourFS, stagingDir, paramsComponentMap, helmChart)
}

func createDir(fs sys.FS, dirName string) error {
	_, err := fs.Stat(dirName)
	if os.IsNotExist(err) {
		if err := fs.MkdirAll(dirName, os.ModePerm); err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func downloadFiles(
	ghClient *github.Client,
	org,
	repo string,
	opt *github.RepositoryContentGetOptions,
	helmChart helmChart) ([]git.File, error) {

	ret := make([]git.File, 0, len(helmChart.Files))
	for _, fileName := range helmChart.Files {
		relativeFilePath := filepath.Join(helmChart.Name, fileName)
		readCloser, err := git.DownloadContents(ghClient, org, repo, relativeFilePath, opt)
		if err != nil {
			return nil, err
		}
		ret = append(ret, git.File{ReadCloser: readCloser, Name: fileName})
	}
	return ret, nil
}

func stageFiles(fs sys.FS, ghFiles []git.File, stagingDir string) {
	for _, ghFile := range ghFiles {
		readCloser := ghFile.ReadCloser

		localFilePath := filepath.Join(stagingDir, ghFile.Name)
		f, err := fs.Create(localFilePath)
		if err != nil {
			log.Fatalf("Error creating file %s (%s)", localFilePath, err)
		}

		if _, err := io.Copy(f, readCloser); err != nil {
			log.Fatalf("Error writing contents to file %s (%s)", localFilePath, err)
		}
		log.Printf("File %s staged in '%s'", ghFile.Name, stagingDir)
		defer readCloser.Close()
		defer f.Close()
	}
}

func updateFilesWithRelease(fp sys.FP, fs sys.FS, release releaseName, walkPath string) error {
	if release.Full == "" || release.Short == "" {
		log.Printf("WORKFLOW_RELEASE (%s) and/or WORKFLOW_RELEASE_SHORT (%s) not provided so not amending staged files.", release.Full, release.Short)
	} else {
		if err := fp.Walk(walkPath, getReleaseWalker().handlerFunc(fs, release)); err != nil {
			return err
		}
	}
	return nil
}
