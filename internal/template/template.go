package template

import (
	"bytes"
	"fmt"
	"text/template"
)

type FileData struct {
	Path    string
	Content string
}

const DefaultTemplate = `
=== {{.Path}} ===

{{.Content}}
`

type Template struct {
	tmpl *template.Template
}

func New() (*Template, error) {
	tmpl, err := template.New("file").Parse(DefaultTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}
	return &Template{tmpl: tmpl}, nil
}

func (t *Template) Format(data FileData) (string, error) {
	var buf bytes.Buffer
	if err := t.tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}
	return buf.String(), nil
}
