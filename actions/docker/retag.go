package docker

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"
	"github.com/deis/deisrel/actions"
	"github.com/deis/deisrel/docker"
	"github.com/deis/deisrel/git"
	"github.com/google/go-github/github"
)

const (
	newOrgFlag    = "new-org"
	defaultNewOrg = "deis"
)

func getAllReposAndShas(
	ghClient *github.Client,
	shaFilePath,
	ref string,
) ([]git.RepoAndSha, error) {
	var allReposAndShas []git.RepoAndSha
	if shaFilePath != "" {
		reposFromFile, err := git.GetShasFromFilepath(shaFilePath)
		if err != nil {
			return nil, fmt.Errorf("getting git SHAs from %s (%s)", shaFilePath, err)
		}
		allReposAndShas = reposFromFile
	}

	reposAndShas, err := git.GetSHAs(ghClient, git.RepoNames(), git.NoTransform, ref)
	if err != nil {
		return nil, fmt.Errorf("getting all SHAs from HEAD on each repository (%s)", err)
	}
	allReposAndShas = reposAndShas
	return allReposAndShas, nil
}

func ensureImages(dockerCl docker.Client, images []*docker.Image) {
	imgsCh, errCh, doneCh := docker.PullImages(dockerCl, images)
	for {
		select {
		case img := <-imgsCh:
			fmt.Printf("pulled %s\n", img.String())
		case err := <-errCh:
			fmt.Printf("error pulling %s (%s)\n", err.Img.String(), err.Err)
		case <-doneCh:
			return
		}
	}
}

func retagAll(dockerCl docker.Client, imageTagPairs []docker.ImageTagPair) {
	pairsCh, errCh, doneCh := docker.RetagImages(dockerCl, imageTagPairs)
	for {
		select {
		case pair := <-pairsCh:
			fmt.Printf("re-tagged %s to %s\n", pair.Source.String(), pair.Target.String())
		case err := <-errCh:
			fmt.Printf(
				"error re-tagging %s to %s (%s)\n",
				err.SourceImage.String(),
				err.TargetImage.String(),
				err.Err,
			)
		case <-doneCh:
			return
		}
	}
}

func pushTargets(dockerCl docker.Client, imageTagPairs []docker.ImageTagPair) {
	images := make([]*docker.Image, len(imageTagPairs))
	for i, imageTagPair := range imageTagPairs {
		images[i] = imageTagPair.Target
	}
	if err := docker.PushImages(dockerCl, images); err != nil {
		log.Printf("Error pushing (%s)", err)
		return
	}
}

func retagCmd(ghClient *github.Client, dockerCl docker.Client) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		newTag := c.Args().Get(0)
		if newTag == "" {
			log.Fatal("This command should have 1 argument to specify the new tag to use")
		}
		newOrg := c.String(newOrgFlag)
		if newOrg == "" {
			newOrg = defaultNewOrg
		}
		// only prompt to push new images if the yes flag was false
		shaFilepath := c.String(actions.ShaFilepathFlag)
		ref := c.String(actions.RefFlag)
		promptPush := !c.Bool(actions.YesFlag)

		allReposAndShas, err := getAllReposAndShas(ghClient, shaFilepath, ref)
		if err != nil {
			log.Fatalf("Error getting all git SHAs (%s)", err)
		}

		repoAndShaList := git.NewRepoAndShaListFromSlice(allReposAndShas)
		repoAndShaList.Sort()
		images, err := docker.ParseImagesFromRepoAndShaList(docker.DeisCIDockerOrg, repoAndShaList)
		if err != nil {
			log.Fatalf("Error parsing docker images (%s)", err)
		}

		fmt.Printf("Pulling %d images\n", len(images))
		ensureImages(dockerCl, images)
		fmt.Println("Re-tagging images...")
		imageTagPairs := docker.CreateImageTagPairsFromTransform(images, func(img docker.Image) *docker.Image {
			img.SetRepo(newOrg)
			img.SetTag(newTag)
			return &img
		})
		retagAll(dockerCl, imageTagPairs)
		fmt.Println("done")

		if promptPush {
			fmt.Println("Pushing new tags")
			pushTargets(dockerCl, imageTagPairs)
		} else {
			fmt.Println("Not pushing newly tagged images")
		}
		return nil
	}
}
