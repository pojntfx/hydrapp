package apk

import (
	_ "embed"

	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
)

//go:embed header.h.tpl
var headerTemplate string

type headerData struct{}

func NewHeaderRenderer() *renderers.Renderer {
	return renderers.NewRenderer("android.h", headerTemplate, headerData{})
}
