package androidmanifest

import (
	"bytes"
	_ "embed"
	"text/template"
)

const (
	path = "AndroidManifest.xml"
)

var (
	//go:embed AndroidManifest.xml
	tpl string
)

type data struct {
	AppID   string
	AppName string
}

type Renderer struct {
	data data
}

func NewRenderer(
	appID string,
	appName string,
) *Renderer {
	return &Renderer{
		data{
			appID,
			appName,
		},
	}
}

func (r *Renderer) Render() (filePath string, fileContent string, err error) {
	t, err := template.New(path).Parse(tpl)
	if err != nil {
		return "", "", err
	}

	buf := bytes.NewBuffer([]byte{})
	if err := t.Execute(buf, r.data); err != nil {
		return "", "", err
	}

	return path, buf.String(), err
}
