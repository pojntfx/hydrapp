package renderers

import (
	"bytes"
	"text/template"
)

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
	t, err := template.New(r.filePath).Parse(r.template)
	if err != nil {
		return "", "", err
	}

	buf := bytes.NewBuffer([]byte{})
	if err := t.Execute(buf, r.data); err != nil {
		return "", "", err
	}

	return r.filePath, buf.String(), err
}
