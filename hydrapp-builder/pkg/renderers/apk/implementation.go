package apk

import (
	_ "embed"

	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
)

//go:embed implementation.c.tpl
var implementationTemplate string

type implementationData struct{}

func NewImplementationRenderer() *renderers.Renderer {
	return renderers.NewRenderer("android.c", implementationTemplate, implementationData{})
}
