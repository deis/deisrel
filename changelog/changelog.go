package changelog

import (
	"text/template"
)

const (
	tplStr = `{{.OldRelease}} -> {{.NewRelease}}

{{ if (len .Features) gt 0 }}
# Features

{{range .Features}} - {{.}}
{{end}}
{{end}}

{{ if (len .Fixes) gt 0 }}
# Fixes

{{range .Fixes}} - {{.}}
{{end}}
{{end}}

{{ if (len .Documentation) gt 0 }}
# Documentation

{{range .Documentation}} - {{.}}
{{end}}
{{end}}

{{ if (len .Maintenance) gt 0 }}
# Maintenance

{{range .Maintenance}} - {{.}}
{{end}}
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
