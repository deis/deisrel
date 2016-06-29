package git

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// GetShasFromFilepath opens the file at path and reads all the git repos and SHAs in the file.
// Returns a slice of all the RepoAndShas that correspond to each entry in the file, or an empty
// slice and non-nil error if there was any error
func GetShasFromFilepath(path string) ([]RepoAndSha, error) {
	ret := []RepoAndSha{}
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
			ret = append(ret, RepoAndSha{
				Name: repoParts[0],
				SHA:  repoParts[1],
			})
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed reading %s: %s", path, err)
	}
	return ret, nil
}
