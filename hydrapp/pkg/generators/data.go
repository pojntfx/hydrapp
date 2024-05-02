package generators

type GoModData struct {
	GoMod string
}

type GoMainData struct {
	GoMod string
}

type AndroidData struct {
	GoMod     string
	JNIExport string
}

type AppTSXData struct {
	AppName string
}

type IndexHTMLData struct {
	AppName string
}

type PackageJSONData struct {
	AppID          string
	AppDescription string
	ReleaseAuthor  string
	ReleaseEmail   string
	LicenseSPDX    string
}

type hydrappYAMLData struct {
	AppID string
}

type ProjectTypeOption struct {
	Name        string
	Description string
}

type CodeOfConductMDData struct {
	ReleaseEmail string
}

type ReadmeMDData struct {
	AppName        string
	AppSummary     string
	AppGitWeb      string
	AppDescription string
	AppBaseURL     string
	AppGit         string
	CurrentYear    string
	ReleaseAuthor  string
	LicenseSPDX    string
	Dir            string
}
