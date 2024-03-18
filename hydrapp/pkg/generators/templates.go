package generators

import (
	_ "embed"
)

var (
	//go:embed icon.png.tpl
	IconTpl []byte

	//go:embed go.mod.tpl
	GoModTpl string

	//go:embed main_panrpc.go.tpl
	GoMainpanrpcTpl string

	//go:embed main_forms.go.tpl
	GoMainFormsTpl string

	//go:embed main_rest.go.tpl
	GoMainRESTTpl string

	//go:embed android_panrpc.go.tpl
	AndroidpanrpcTpl string

	//go:embed android_forms.go.tpl
	AndroidFormsTpl string

	//go:embed android_rest.go.tpl
	AndroidRESTTpl string

	//go:embed .gitignore_panrpc_parcel.tpl
	GitignorepanrpcParcelTpl string

	//go:embed .gitignore_panrpc_cra.tpl
	GitignorepanrpcCRATpl string

	//go:embed .gitignore_rest.tpl
	GitignoreRESTTpl string

	//go:embed backend_panrpc.go.tpl
	BackendpanrpcTpl string

	//go:embed backend_rest.go.tpl
	BackendRESTTpl string

	//go:embed frontend_panrpc.go.tpl
	FrontendpanrpcTpl string

	//go:embed frontend_rest.go.tpl
	FrontendRESTTpl string

	//go:embed frontend_forms.go.tpl
	FrontendFormsTpl string

	//go:embed App.tsx.tpl
	AppTSXTpl string

	//go:embed main.tsx.tpl
	MainTSXTpl string

	//go:embed index_panrpc_parcel.html.tpl
	IndexHTMLpanrpcParcelTpl string

	//go:embed index_panrpc_cra.html.tpl
	IndexHTMLpanrpcCRATpl string

	//go:embed index_rest.html.tpl
	IndexHTMLRESTTpl string

	//go:embed index_forms.html.tpl
	IndexHTMLFormsTpl string

	//go:embed package_parcel.json.tpl
	PackageJSONParcelTpl string

	//go:embed package_cra.json.tpl
	PackageJSONCRATpl string

	//go:embed tsconfig.json.tpl
	TsconfigJSONTpl string

	//go:embed hydrapp.yaml.tpl
	HydrappYAMLTpl string

	//go:embed CODE_OF_CONDUCT.md.tpl
	CodeOfConductMDTpl string

	//go:embed README.md.tpl
	ReadmeMDTpl string
)
