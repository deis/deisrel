package helm

import (
	"io"
	"text/template"

	"github.com/google/go-github/github"
)

type GenerateParamsData struct {
}

// Chart represents a helm chart that should be rendered by deisrel. It has a name, a template
// to be used to render the generate_params.toml file, and a slice of files that should also be
// downloaded
type Chart struct {
	Name     string
	template *template.Template
	Files    []string
}

// RenderGenerateParamsTpl renders the Chart's generate_params template to w using paramsMap and returns the error
func (c Chart) RenderGenerateParamsTpl(w io.Writer, paramsMap GenParamsComponentMap) error {
	return c.template.Execute(w, paramsMap)
}

// newChart creates a new chart with the given name and files
func newChart(ghClient *github.Client, name string, files []string) (*Chart, error) {
	genParamsTpl, err := getGenerateParamsTpl(ghClient, name)
	if err != nil {
		return nil, err
	}
	return &Chart{
		Name:     name,
		template: genParamsTpl,
		Files:    files,
	}, nil
}
