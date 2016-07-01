package github

import (
	"github.com/codegangsta/cli"
	"github.com/google/go-github/github"
)

func createReleasesCmd(ghClient *github.Client) func(c *cli.Command) error {
	return nil
}
