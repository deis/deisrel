package actions

import (
	"log"

	"github.com/google/go-github/github"
)

func GenerateChangelogs(ghClient *github.Client) func(c *cli.Context) {
	return func(c *cli.Context) {
		log.Printf("TODO")
	}
}
