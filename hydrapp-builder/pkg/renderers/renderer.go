package renderers

import (
	"bytes"
	"text/template"
)

type Renderer[T interface{}] struct {
	filePath string
	template string
	data     T
}

func NewRenderer[T interface{}](
	filePath string,
	template string,
	data T,
) *Renderer[T] {
	return &Renderer[T]{
		filePath,
		template,
		data,
	}
}

func (r *Renderer[T]) Render() (filePath string, fileContent string, err error) {
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
