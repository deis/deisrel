package actions

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"
	"github.com/deis/deisrel/git"
	"github.com/google/go-github/github"
)

const (
	// ShortFlag is the cli flag that indicates whether to show short or long SHAs
	ShortFlag = "short"
)

// GetShas is the CLI action for getting github shas of all of the Deis Workflow repos
func GetShas(ghClient *github.Client) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		transformFunc := git.NoTransform
		if c.Bool(ShortFlag) {
			transformFunc = git.ShortSHATransformNoErr
		}

		reposAndShas, err := git.GetSHAs(ghClient, repoNames, transformFunc, c.String(RefFlag))
		if err != nil {
			log.Fatal(err)
		}
		for _, repoAndSha := range reposAndShas {
			fmt.Printf("%s - %s\n", repoAndSha.Name, repoAndSha.SHA)
		}
		return nil
	}
}
