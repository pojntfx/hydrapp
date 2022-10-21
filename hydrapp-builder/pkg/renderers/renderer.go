package renderers

import (
	"bytes"
	"text/template"
	"time"
)

type Release struct {
	Version     string    `json:"version" yaml:"version"`
	Date        time.Time `json:"date" yaml:"date"`
	Description string    `json:"description" yaml:"description"`
	Author      string    `json:"author" yaml:"author"`
	Email       string    `json:"email" yaml:"email"`
}

type Renderer struct {
	filePath string
	template string
	data     interface{}
}

func NewRenderer(
	filePath string,
	template string,
	data interface{},
) *Renderer {
	return &Renderer{
		filePath,
		template,
		data,
	}
}

func (r *Renderer) Render() (filePath string, fileContent string, err error) {
	t, err := template.
		New(r.filePath).
		Funcs(template.FuncMap{
			"LastRelease": func(releases []Release) Release {
				return releases[len(releases)-1]
			},
		}).
		Parse(r.template)
	if err != nil {
		return "", "", err
	}

	buf := bytes.NewBuffer([]byte{})
	if err := t.Execute(buf, r.data); err != nil {
		return "", "", err
	}

	return r.filePath, buf.String(), err
}
