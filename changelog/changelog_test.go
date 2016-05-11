package changelog

import (
	"bytes"
	"strings"
	"testing"
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
				Documentation: documentation,
				Maintenance:   maintenance,
			},
			missing: "# Features",
		},
		testCase{
			vals: Values{
				OldRelease:    oldRelease,
				NewRelease:    newRelease,
				Features:      features,
				Fixes:         nil,
				Documentation: documentation,
				Maintenance:   maintenance,
			},
			missing: "# Fixes",
		},
		testCase{
			vals: Values{
				OldRelease:    oldRelease,
				NewRelease:    newRelease,
				Features:      features,
				Fixes:         fixes,
				Documentation: nil,
				Maintenance:   maintenance,
			},
			missing: "# Documentation",
		},
		testCase{
			vals: Values{
				OldRelease:    oldRelease,
				NewRelease:    newRelease,
				Features:      features,
				Fixes:         fixes,
				Documentation: documentation,
				Maintenance:   nil,
			},
			missing: "# Maintenance",
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
