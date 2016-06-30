package branches

import (
	"fmt"
	"log"

	"github.com/codegangsta/cli"
	"github.com/deis/deisrel/actions"
	"github.com/deis/deisrel/git"
	"github.com/google/go-github/github"
)

func createCmd(ghClient *github.Client) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		branchName := c.String(branchNameFlag)
		if branchName == "" {
			log.Fatalf("Branch name not specified")
		}
		proceed := c.Bool(actions.YesFlag)
		ref := c.String(actions.RefFlag)

		repoNames := git.RepoNames()
		fmt.Printf("Getting SHAs for all repositories on '%s'\n", ref)
		reposAndSHAs, err := git.GetSHAs(ghClient, repoNames, git.NoTransform, ref)
		if err != nil {
			log.Fatalf("Error getting SHAs for repositories (%s)", err)
		}
		for _, ras := range reposAndSHAs {
			fmt.Printf("%s - %s\n", ras.Name, ras.SHA)
		}

		fmt.Println()
		fmt.Printf("Creating branch %s on all repositories\n", branchName)
		if proceed {
			rasl, err := git.CreateBranches(ghClient, branchName, reposAndSHAs)
			if err != nil {
				log.Fatalf("Error creating branches for repositories (%s)", err)
			}

			for _, ras := range rasl {
				fmt.Printf("Created branch %s on %s\n", branchName, ras.Name)
			}
		} else {
			fmt.Printf("Not creating branches. '%s' flag was false\n", actions.YesFlag)
		}
		return nil
	}
}
