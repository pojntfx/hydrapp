package apk

import (
	_ "embed"

	"github.com/pojntfx/hydrapp/hydrapp/pkg/renderers"
)

//go:embed header.h.tpl
var headerTemplate string

type headerData struct{}

func NewHeaderRenderer() renderers.Renderer {
	return renderers.NewRenderer("hydrapp_android.h", headerTemplate, headerData{})
}
