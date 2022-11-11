package renderers

import (
	"bytes"
	"strings"
	"text/template"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

func (r *Renderer) Render(templateOverride string) (filePath string, fileContent string, err error) {
	titler := cases.Title(language.English)

	t, err := template.
		New(r.filePath).
		Funcs(template.FuncMap{
			"LastRelease": func(releases []Release) Release {
				return releases[len(releases)-1]
			},
			"Titlecase": func(title string) string {
				return titler.String(title)
			},
		}).
		Parse(func() string {
			if strings.TrimSpace(templateOverride) != "" {
				return templateOverride
			}

			return r.template
		}())
	if err != nil {
		return "", "", err
	}

	buf := bytes.NewBuffer([]byte{})
	if err := t.Execute(buf, r.data); err != nil {
		return "", "", err
	}

	return r.filePath, buf.String(), err
}
