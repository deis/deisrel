package changelog

import (
	"bytes"
	"strings"
	"testing"

	"github.com/arschles/assert"
)

const (
	oldRelease = "old"
	newRelease = "new"
)

var (
	features      = []string{"feature1"}
	fixes         = []string{"fix1"}
	documentation = []string{"doc1"}
	maintenance   = []string{"maint1"}
	refactor      = []string{"ref1"}
)

func TestTemplate(t *testing.T) {
	type testCase struct {
		vals    Values
		missing string
	}
	testCases := []testCase{
		testCase{
			vals: Values{
				OldRelease:    oldRelease,
				NewRelease:    newRelease,
				Features:      features,
				Fixes:         fixes,
				Refactors:     refactor,
				Documentation: documentation,
				Maintenance:   maintenance,
			},
			missing: "",
		},
		testCase{
			vals: Values{
				OldRelease:    oldRelease,
				NewRelease:    newRelease,
				Features:      nil,
				Fixes:         fixes,
				Refactors:     refactor,
				Documentation: documentation,
				Maintenance:   maintenance,
			},
			missing: "#### Features",
		},
		testCase{
			vals: Values{
				OldRelease:    oldRelease,
				NewRelease:    newRelease,
				Features:      features,
				Fixes:         nil,
				Refactors:     refactor,
				Documentation: documentation,
				Maintenance:   maintenance,
			},
			missing: "#### Fixes",
		},
		testCase{
			vals: Values{
				OldRelease:    oldRelease,
				NewRelease:    newRelease,
				Features:      features,
				Fixes:         fixes,
				Refactors:     refactor,
				Documentation: nil,
				Maintenance:   maintenance,
			},
			missing: "#### Documentation",
		},
		testCase{
			vals: Values{
				OldRelease:    oldRelease,
				NewRelease:    newRelease,
				Features:      features,
				Fixes:         fixes,
				Refactors:     refactor,
				Documentation: documentation,
				Maintenance:   nil,
			},
			missing: "#### Maintenance",
		},
		testCase{
			vals: Values{
				OldRelease:    oldRelease,
				NewRelease:    newRelease,
				Features:      features,
				Fixes:         fixes,
				Refactors:     nil,
				Documentation: documentation,
				Maintenance:   maintenance,
			},
			missing: "#### Refactors",
		},
	}

	for i, testCase := range testCases {
		var buf bytes.Buffer
		if err := Tpl.Execute(&buf, testCase.vals); err != nil {
			t.Errorf("Error executing template %d (%s)", i, err)
			continue
		}
		if len(testCase.missing) > 0 {
			outStr := string(buf.Bytes())
			if strings.Contains(outStr, testCase.missing) {
				t.Errorf("Expected [%s] to be missing from the rendered template %d, but found it", testCase.missing, i)
			}
		}
	}
}

func TestMergeValues(t *testing.T) {
	val1 := Values{Features: []string{"feat1"}}
	val2 := Values{Fixes: []string{"fix1"}, Features: []string{"feat2"}}
	res := MergeValues("old", "new", []Values{val1, val2})
	assert.Equal(t, res.OldRelease, "old", "old release")
	assert.Equal(t, res.NewRelease, "new", "new release")
	assert.Equal(t, len(res.Features), 2, "length of features slice")
	assert.Equal(t, len(res.Fixes), 1, "length of fixes slice")
	assert.Equal(t, res.Features, []string{"feat1", "feat2"}, "features slice")
	assert.Equal(t, res.Fixes, []string{"fix1"}, "fixes slice")
}
