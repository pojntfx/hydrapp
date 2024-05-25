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

type Renderer interface {
	Render(templateOverride string) (filePath string, fileContent []byte, err error)
}

type renderer struct {
	filePath string
	template string
	data     interface{}
}

func NewRenderer(
	filePath string,
	template string,
	data interface{},
) Renderer {
	return &renderer{
		filePath,
		template,
		data,
	}
}

func (r *renderer) Render(templateOverride string) (filePath string, fileContent []byte, err error) {
	titler := cases.Title(language.English)

	t, err := template.
		New(r.filePath).
		Funcs(template.FuncMap{
			"LastRelease": func(releases []Release) Release {
				return releases[0]
			},
			"Titlecase": func(title string) string {
				return titler.String(title)
			},
			"DeveloperID": func(appID string) string {
				parts := strings.Split(appID, ".")
				if len(parts) > 1 {
					parts = parts[0 : len(parts)-2]
				}

				return strings.Join(parts, ".")
			},
		}).
		Parse(func() string {
			if strings.TrimSpace(templateOverride) != "" {
				return templateOverride
			}

			return r.template
		}())
	if err != nil {
		return "", []byte{}, err
	}

	buf := bytes.NewBuffer([]byte{})
	if err := t.Execute(buf, r.data); err != nil {
		return "", []byte{}, err
	}

	return r.filePath, buf.Bytes(), err
}
