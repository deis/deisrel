package github

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"
	"github.com/deis/deisrel/actions"
	"github.com/deis/deisrel/changelog"
	"github.com/deis/deisrel/git"
	"github.com/google/go-github/github"
)

const (
	releaseTagFlag   = "release-tag"
	releaseTitleFlag = "release-title"
)

func createReleasesCmd(ghClient *github.Client) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		releaseTag := c.String(releaseTagFlag)
		releaseTitle := c.String(releaseTitleFlag)
		dryRun := c.Bool(actions.YesFlag)
		ref := c.String(actions.RefFlag)
		repoNames := git.RepoNames()
		reposAndValues := make([]changelog.RepoAndValues, len(repoNames))
		for i, repoName := range repoNames {
			reposAndValues[i] = changelog.RepoAndValues{
				RepoAndSHA: git.RepoAndSha{Name: repoName, SHA: ref},
			}
		}
		fmt.Println("Fetching changelogs for all repos")
		repoAndVals, err := changelog.MultiRepoVals(ghClient, reposAndValues)
		if err != nil {
			log.Fatalf("Error fetching changelog values (%s)", err)
		}

		relInfos := make([]git.ReleaseInfo, len(repoAndVals))
		for i, repoAndVal := range repoAndVals {
			body, err := repoAndVal.Values.RenderTplToString()
			if err != nil {
				log.Fatalf("Error rendering changelog for repository %s (%s)", repoAndVal.RepoAndSHA.Name, err)
			}
			relInfos[i] = git.ReleaseInfo{
				RepoName: repoAndVal.RepoAndSHA.Name,
				Title:    releaseTitle,
				Tag:      releaseTag,
				Body:     body,
			}
		}

		if dryRun {
			fmt.Println("Not creating releases. Dry run flag present")
		} else {
			if err := git.CreateReleases(ghClient, relInfos); err != nil {
				log.Fatalf("Error creating GH releases (%s)", err)
			}
		}

		return nil
	}
}
