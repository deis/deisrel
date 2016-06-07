package git

import (
	"io"

	"github.com/google/go-github/github"
)

// DownloadContents downloads filepath from org/repo on Github. If error is non-nil, the returned io.ReadCloser will contain the full contents of the file
func DownloadContents(
	ghClient *github.Client,
	org,
	repo,
	filepath string,
	opt *github.RepositoryContentGetOptions) (io.ReadCloser, error) {

	rc, err := ghClient.Repositories.DownloadContents(org, repo, filepath, opt)
	if err != nil {
		return nil, err
	}
	return rc, nil
}
