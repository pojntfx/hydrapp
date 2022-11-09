package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"text/template"

	_ "embed"

	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/builders"
	cconfig "github.com/pojntfx/hydrapp/hydrapp-builder/pkg/config"
)

//go:embed INSTALLATION.md.tpl
var installationTemplate string

type installationData struct {
	AndroidRepoURL  string
	AppName         string
	MacOSBinaryURL  string
	MacOSBinaryName string
	AppID           string
	Flatpaks        []flatpak
}

type flatpak struct {
	Architecture string
	RepoURL      string
}

func main() {
	config := flag.String("config", "hydrapp.yaml", "Config file to use")
	branchID := flag.String("branch-id", "", `Branch ID to build the app as, i.e. main (for an app ID like "myappid.main" and baseURL like "mybaseurl/main"`)
	branchName := flag.String("branch-name", "", `Branch name to build the app as, i.e. Main (for an app name like "myappname (Main)"`)

	flag.Parse()

	content, err := ioutil.ReadFile(*config)
	if err != nil {
		panic(err)
	}

	cfg, err := cconfig.Parse(content)
	if err != nil {
		panic(err)
	}

	t, err := template.
		New("INSTALLATION.md").
		Parse(installationTemplate)
	if err != nil {
		panic(err)
	}

	appName := builders.GetAppNameForBranch(cfg.App.Name, *branchName)
	macOSBinaryName := builders.GetAppIDForBranch(cfg.App.ID, *branchID) + ".darwin.dmg"

	flatpaks := []flatpak{}
	for _, f := range cfg.Flatpak {
		flatpaks = append(flatpaks, flatpak{
			Architecture: f.Architecture,
			RepoURL:      cfg.App.BaseURL + builders.GetPathForBranch(f.Path, *branchID) + "/hydrapp.flatpakrepo",
		})
	}

	buf := bytes.NewBuffer([]byte{})
	if err := t.Execute(buf, installationData{
		AppID:           builders.GetAppIDForBranch(cfg.App.ID, *branchID),
		AppName:         appName,
		AndroidRepoURL:  cfg.App.BaseURL + builders.GetPathForBranch(cfg.APK.Path, *branchID),
		MacOSBinaryURL:  cfg.App.BaseURL + builders.GetPathForBranch(cfg.DMG.Path, *branchID) + "/" + macOSBinaryName,
		MacOSBinaryName: macOSBinaryName,
		Flatpaks:        flatpaks,
	}); err != nil {
		panic(err)
	}

	fmt.Println(buf.String())
}
