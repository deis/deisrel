package changelog

import (
	"text/template"
)

const (
	tplStr = `### {{.OldRelease}} -> {{.NewRelease}}
{{ if (len .Features) gt 0 }}
#### Features
{{range .Features}}- {{.}}
{{end}}
{{ end -}}
{{ if (len .Fixes) gt 0 -}}
#### Fixes
{{range .Fixes}}- {{.}}
{{end}}
{{ end }}
{{ if (len .Documentation) gt 0 }}
#### Documentation
{{range .Documentation}}- {{.}}
{{end}}
{{end}}
{{ if (len .Maintenance) gt 0 }}
#### Maintenance
{{range .Maintenance}}- {{.}}
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

// MergeValues merges all of the slices in vals together into a single Values struct which has OldRelease set to oldRel and NewRelease set to newRel
func MergeValues(oldRel, newRel string, vals []Values) *Values {
	ret := &Values{OldRelease: oldRel, NewRelease: newRel}
	for _, val := range vals {
		ret.Features = append(ret.Features, val.Features...)
		ret.Fixes = append(ret.Fixes, val.Fixes...)
		ret.Documentation = append(ret.Documentation, val.Documentation...)
		ret.Maintenance = append(ret.Maintenance, val.Maintenance...)
	}
	return ret
}
