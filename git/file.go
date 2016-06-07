package git

import (
	"io"
)

// File is the representation of a downloaded file on github
type File struct {
	ReadCloser io.ReadCloser
	Name       string
}
