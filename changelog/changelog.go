package changelog

import (
	"text/template"
)

const (
	tplStr = `{{.OldRelease}} -> {{.NewRelease}}

# Features

{{range .Features}} - {{.}}
{{else}}No new features for this release.
{{end}}

# Fixes

{{range .Fixes}} - {{.}}
{{else}}No bug fixes for this release.
{{end}}

# Documentation

{{range .Documentation}} - {{.}}
{{else}}No new documentation for this release.
{{end}}

# Maintenance

{{range .Maintenance}} - {{.}}
{{else}}No maintenance required for this release.
{{end}}`
)

var (
	// Tpl is the standard changelog template. Execute it with a Values struct
	Tpl = template.Must(template.New("changelog").Parse(tplStr))
)

// Values represents the values that are required to render a changelog
type Values struct {
	OldRelease    string
	NewRelease    string
	Features      []string
	Fixes         []string
	Documentation []string
	Maintenance   []string
}
