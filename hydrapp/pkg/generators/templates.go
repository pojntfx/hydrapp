package generators

import (
	_ "embed"
)

var (
	//go:embed icon.png.tpl
	IconTpl []byte

	//go:embed go.mod.tpl
	GoModTpl string

	//go:embed main_react_panrpc.go.tpl
	GoMainReactPanrpcTpl string

	//go:embed main_vanillajs_forms.go.tpl
	GoMainVanillaJSFormsTpl string

	//go:embed main_vanillajs_rest.go.tpl
	GoMainVanillaJSRESTTpl string

	//go:embed android_react_panrpc.go.tpl
	AndroidReactPanrpcTpl string

	//go:embed android_vanillajs_forms.go.tpl
	AndroidVanillaJSFormsTpl string

	//go:embed android_vanillajs_rest.go.tpl
	AndroidVanillaJSRESTTpl string

	//go:embed .gitignore_react_panrpc.tpl
	GitignoreReactPanrpcTpl string

	//go:embed .gitignore_vanillajs_rest.tpl
	GitignoreVanillaJSRESTTpl string

	//go:embed backend_react_panrpc.go.tpl
	BackendReactPanrpcTpl string

	//go:embed backend_vanillajs_rest.go.tpl
	BackendVanillaJSRESTTpl string

	//go:embed frontend_react_panrpc.go.tpl
	FrontendReactPanrpcTpl string

	//go:embed frontend_vanillajs_rest.go.tpl
	FrontendVanillaJSRESTTpl string

	//go:embed frontend_vanillajs_forms.go.tpl
	FrontendVanillaJSFormsTpl string

	//go:embed App.tsx.tpl
	AppTSXTpl string

	//go:embed main.tsx.tpl
	MainTSXTpl string

	//go:embed index_react_panrpc.html.tpl
	IndexHTMLReactPanrpcTpl string

	//go:embed index_vanillajs_rest.html.tpl
	IndexHTMLVanillaJSRESTTpl string

	//go:embed index_vanillajs_forms.html.tpl
	IndexHTMLVanillaJSFormsTpl string

	//go:embed package_react_panrpc.json.tpl
	PackageJSONReactPanrpcTpl string

	//go:embed tsconfig.json.tpl
	TsconfigJSONTpl string

	//go:embed hydrapp.yaml.tpl
	HydrappYAMLTpl string

	//go:embed CODE_OF_CONDUCT.md.tpl
	CodeOfConductMDTpl string

	//go:embed README.md.tpl
	ReadmeMDTpl string
)
