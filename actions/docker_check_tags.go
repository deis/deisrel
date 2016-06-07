package actions

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/coreos/go-quay/quay"
	"github.com/google/go-github/github"

	reg "github.com/deis/deisrel/registry"
)

const (
	quayTokenEnvVarName = "QUAY_AUTH_TOKEN"
)

// DockerCheckTags is the CLI action for checking that Docker image tags
// exist in registries we are interested in
func DockerCheckTags(ghClient *github.Client) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		ref := c.GlobalString(RefFlag)

		// setup hub registry client
		hubClient, err := reg.GetHub(reg.HubRegistryURL, "", "")
		if err != nil {
			return err
		}
		hub := reg.NewHubRegistry(hubClient)
		// setup quay registry client
		quayTkn := os.Getenv(quayTokenEnvVarName)
		if quayTkn == "" {
			log.Fatalf("'%s' env var required", quayTokenEnvVarName)
		}
		quay := reg.NewQuayRegistry(quay.NewHTTPClient(nil), reg.NewQuayAuth(quayTkn))

		// latest sha values should exist in 'git-<sha>' form as tags in registries
		repoAndShas, err := getShas(ghClient, repoNames, shortShaTransform, ref)
		if err != nil {
			return err
		}

		foundImgTags, errs := dockerCheckTags(ghClient, quay, hub, repoAndShas, ref)
		if len(errs) > 0 {
			for _, err := range errs {
				fmt.Printf("Error encountered attempting checking tag (%s)\n", err)
			}
		}

		fmt.Println("Successfully found the following image tags in Quay.io and DockerHub registries:")
		for _, foundImgTag := range foundImgTags {
			fmt.Println(foundImgTag)
		}
		return nil
	}
}
