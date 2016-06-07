package actions

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/arschles/sys"
)

type releaseWalker struct {
	filepath.WalkFunc
	sys.FakeFS
}

func getReleaseWalker() *releaseWalker {
	return &releaseWalker{}
}

func (r *releaseWalker) handlerFunc(fs sys.FS, release releaseName) filepath.WalkFunc {
	return func(path string, fi os.FileInfo, err error) error {
		return r.walk(path, release, fi, err, fs)
	}
}

func (r *releaseWalker) walk(path string, release releaseName, fi os.FileInfo, err error, fs sys.FS) error {
	if err != nil {
		return err
	}

	if fi.IsDir() {
		return nil
	}

	read, err := fs.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading file %s (%s)", path, err)
	}

	newContents := strings.Replace(string(read), "-dev", "-"+release.Short, -1)
	newContents = strings.Replace(newContents, "v2.0.0", release.Full, -1)

	if _, err := fs.WriteFile(path, []byte(newContents), 0); err != nil {
		log.Fatalf("Error writing contents to file %s (%s)", path, err)
	}

	fmt.Printf("File '%s' updated with release '%s'\n", path, release.Short)
	return nil
}
