package changelog

import (
	"testing"

	"github.com/arschles/assert"
)

func TestCommitFocus(t *testing.T) {
	// find a valid focus
	assert.Equal(t, commitFocus("fix(builder): some stuff"), "builder", "focus")
	// find a missing focus
	assert.Equal(t, commitFocus("fix: some stuff"), "*", "focus")
}

func TestCommitTitle(t *testing.T) {
	// find a valid title
	assert.Equal(t, commitTitle("fix(builder): some stuff"), "some stuff", "title")
	// ensure the whole thing is dumped when there's no parseable title
	assert.Equal(t, commitTitle("fix(builder) some stuff"), "fix(builder) some stuff", "title")
	// ensure everything after the first ':' is returned
	assert.Equal(t, commitTitle("fix(builder): stuff1: stuff2"), "stuff1: stuff2", "title")
}
