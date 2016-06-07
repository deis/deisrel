package registry

import (
	"testing"

	"github.com/arschles/assert"
)

func TestSliceContains(t *testing.T) {
	testSlice := []string{"a", "b"}
	// slice contains
	assert.True(t, sliceContains(testSlice, "b"), "slice contains")
	// slice does not contain
	assert.False(t, sliceContains(testSlice, "c"), "slice doesn't contain")
}
