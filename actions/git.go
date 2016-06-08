package actions

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/deis/deisrel/git"
	"github.com/google/go-github/github"
)

// GitTag creates the CLI action for the 'deisrel git tag' command
func GitTag(client *github.Client) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		tag := c.Args().Get(0)
		shaFilepath := c.String(ShaFilepathFlag)

		if tag == "" {
			log.Fatal("Usage: deisrel git tag <options> <tag>")
		}

		repos, err := git.GetSHAs(client, allGitRepoNames, git.NoTransform, c.String(RefFlag))
		if err != nil {
			log.Fatal(err)
		}
		if shaFilepath != "" {
			// update the latest shas with the shas in shaFilepath
			reposFromFile, err := git.GetShasFromFilepath(shaFilepath)
			if err != nil {
				log.Fatal(err)
			}
			// merge the latest shas with the shas in shaFilePath, since the file may only
			// specify a subset of the latest repos
			for _, repo := range repos {
				for _, updatedRepo := range reposFromFile {
					if updatedRepo.Name == repo.Name {
						repo.SHA = updatedRepo.SHA
					}
				}
			}
		}
		fmt.Println("=== Repos")
		rasList := git.NewEmptyRepoAndShaList()
		for _, repo := range repos {
			rasList.Add(repo)
		}
		rasList.Sort()
		fmt.Println(rasList.String())

		ok := true
		if !c.Bool(YesFlag) {
			var err error
			ok, err = prompt()
			if err != nil {
				log.Fatal(err)
			}
		}

		if ok {
			if err := git.CreateTags(client, repos, tag); err != nil {
				log.Fatal(fmt.Errorf("could create tag %s: %v", tag, err))
			}
		}
		return nil
	}
}

func prompt() (bool, error) {
	acceptableAnswers := []string{
		"y",
		"yes",
		"yea",
		"yep",
		"sure",
		"ok",
		"okey-dokey",
		"affirmative",
		"aye aye, captain",
		"roger",
		"fo' shizzle",
		"totally",
		"oui",
		"sí",
		"`ae",
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Is this okay? (y/N) ")
	text, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	for _, ans := range acceptableAnswers {
		if strings.TrimSpace(strings.ToLower(text)) == ans {
			return true, nil
		}
	}
	return false, nil
}
