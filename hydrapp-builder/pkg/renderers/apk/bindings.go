package apk

import (
	_ "embed"
	"strings"

	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
)

//go:embed bindings.go.tpl
var bindingsTemplate string

type bindingsData struct {
	AppID              string
	AppBackendPackage  string
	AppFrontendPackage string
}

func NewBindingsRenderer(
	appID string,
	appBackendPkg string,
	appFrontendPkg string,
) *renderers.Renderer {
	return renderers.NewRenderer("main_android.go", bindingsTemplate, bindingsData{strings.Replace(appID, ".", "_", -1), appBackendPkg, appFrontendPkg})
}
