package actions

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/codegangsta/cli"
	"github.com/google/go-github/github"
)

// TODO: breakup and move to 'git' package https://github.com/deis/deisrel/issues/60
type repoAndSha struct {
	repoName string
	sha      string
}

func noTransform(s string) string       { return s }
func shortShaTransform(s string) string { return s[:7] }
func imageTagTransform(s string) string { return fmt.Sprintf("git-%s", shortShaTransform(s)) }

func getShas(ghClient *github.Client, repos []string, transform func(string) string, ref string) ([]repoAndSha, error) {
	outCh := make(chan repoAndSha)
	errCh := make(chan error)
	doneCh := make(chan struct{})
	var wg sync.WaitGroup
	for _, repo := range repos {
		wg.Add(1)
		ch := make(chan repoAndSha)
		ech := make(chan error)
		go func(repo string) {
			commitsListOpts := &github.CommitsListOptions{
				SHA: ref,
				ListOptions: github.ListOptions{
					Page:    1,
					PerPage: 1,
				},
			}
			repoCommits, _, err := ghClient.Repositories.ListCommits("deis", repo, commitsListOpts)
			if err != nil {
				ech <- fmt.Errorf("Error listing commits for repo %s (%s)", repo, err)
				return
			}
			if len(repoCommits) < 1 {
				ech <- fmt.Errorf("No commits found for repo %s", repo)
				return
			}
			repoCommit := repoCommits[0]
			sha := transform(*repoCommit.SHA)
			ch <- repoAndSha{repoName: repo, sha: sha}
		}(repo)
		go func() {
			defer wg.Done()
			select {
			case e := <-ech:
				errCh <- e
			case o := <-ch:
				outCh <- o
			}
		}()
	}
	go func() {
		wg.Wait()
		close(doneCh)
	}()

	ret := []repoAndSha{}
	for {
		select {
		case <-doneCh:
			return ret, nil
		case str := <-outCh:
			ret = append(ret, str)
		case err := <-errCh:
			return nil, err
		}
	}
}

func getLastTag(ghClient *github.Client, repos []string) (map[string]*github.RepositoryTag, error) {
	for _, repo := range repos {
		_, _, err := ghClient.Repositories.ListTags("deis", repo, nil)
		if err != nil {
			return make(map[string]*github.RepositoryTag), err
		}
	}
	return make(map[string]*github.RepositoryTag), nil
}

func downloadContents(ghClient *github.Client, org, repo, filepath string, opt *github.RepositoryContentGetOptions) (io.ReadCloser, error) {
	rc, err := ghClient.Repositories.DownloadContents(org, repo, filepath, opt)
	if err != nil {
		return nil, err
	}
	return rc, nil
}

func GitTag(client *github.Client) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		tag := c.Args().Get(0)
		shaFilepath := c.String(ShaFilepathFlag)

		if tag == "" {
			log.Fatal("Usage: deisrel git tag <options> <tag>")
		}

		repos, err := getShas(client, allGitRepoNames, noTransform, c.String(RefFlag))
		if err != nil {
			log.Fatal(err)
		}
		if shaFilepath != "" {
			// update the latest shas with the shas in shaFilepath
			reposFromFile, err := getShasFromFilepath(shaFilepath)
			if err != nil {
				log.Fatal(err)
			}
			// merge the latest shas with the shas in shaFilePath, since the file may only
			// specify a subset of the latest repos
			for _, repo := range repos {
				for _, updatedRepo := range reposFromFile {
					if updatedRepo.repoName == repo.repoName {
						repo.sha = updatedRepo.sha
					}
				}
			}
		}
		fmt.Println("=== Repos")
		for _, repo := range repos {
			fmt.Printf("%s: %s\n", repo.repoName, repo.sha)
		}
		var ok bool = true
		if !c.Bool(YesFlag) {
			var err error
			ok, err = prompt()
			if err != nil {
				log.Fatal(err)
			}
		}

		if ok {
			if err := createGitTag(client, repos, tag); err != nil {
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

func getShasFromFilepath(path string) ([]repoAndSha, error) {
	ret := []repoAndSha{}
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open %s: %s", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.ContainsRune(line, '=') {
			repoParts := strings.SplitN(line, "=", 2)
			ret = append(ret, repoAndSha{
				repoName: repoParts[0],
				sha:      repoParts[1],
			})
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed reading %s: %s", path, err)
	}
	return ret, nil
}

// createGitTag tags each repository at the given SHA with the tag supplied.
func createGitTag(client *github.Client, repos []repoAndSha, tag string) error {
	for _, repo := range repos {
		ref := &github.Reference{
			Ref: github.String("refs/tags/" + tag),
			Object: &github.GitObject{
				SHA: github.String(repo.sha),
			},
		}
		_, _, err := client.Git.CreateRef("deis", repo.repoName, ref)
		// GitHub returns HTTP 422 Unprocessable Entity when a field is invalid,
		// such as when a reference already exists or the sha does not exist
		// https://developer.github.com/v3/#client-errors
		if err != nil && !strings.Contains(err.Error(), "Reference already exists") {
			return err
		}
		fmt.Printf("%s: created tag %s\n", repo.repoName, tag)
	}
	return nil
}
