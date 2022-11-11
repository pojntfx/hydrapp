package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
	"text/template"

	_ "embed"

	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/builders"
	cconfig "github.com/pojntfx/hydrapp/hydrapp-builder/pkg/config"
	"github.com/pojntfx/hydrapp/hydrapp-utils/pkg/update"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	//go:embed INSTALLATION.md.tpl
	installationTemplate string

	errInvalidRPMDistro = errors.New("invalid RPM distro")
)

type installationData struct {
	AndroidRepoURL  string
	AppName         string
	MacOSBinaryURL  string
	MacOSBinaryName string
	AppID           string
	Flatpaks        []artifact
	MSIs            []artifact
	RPMs            []distroArtifact
	DEBs            []distroArtifact
	BinariesURL     string
}

type artifact struct {
	Architecture string
	URL          string
}

type distroArtifact struct {
	artifact
	DistroName    string
	DistroVersion string
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

	titler := cases.Title(language.English)

	t, err := template.
		New("INSTALLATION.md").
		Funcs(template.FuncMap{
			"Titlecase": func(title string) string {
				return titler.String(title)
			},
		}).
		Parse(installationTemplate)
	if err != nil {
		panic(err)
	}

	appName := builders.GetAppNameForBranch(cfg.App.Name, *branchName)
	macOSBinaryName := builders.GetAppIDForBranch(cfg.App.ID, *branchID) + ".darwin.dmg"

	flatpaks := []artifact{}
	for _, f := range cfg.Flatpak {
		flatpaks = append(flatpaks, artifact{
			Architecture: f.Architecture,
			URL:          cfg.App.BaseURL + builders.GetPathForBranch(f.Path, *branchID) + "/hydrapp.flatpakrepo",
		})
	}

	msis := []artifact{}
	for _, m := range cfg.MSI {
		msis = append(msis, artifact{
			Architecture: m.Architecture,
			URL:          cfg.App.BaseURL + builders.GetPathForBranch(m.Path, *branchID) + "/" + builders.GetAppIDForBranch(cfg.App.ID, *branchID) + ".windows-" + update.GetArchIdentifier(m.Architecture) + ".msi",
		})
	}

	rpms := []distroArtifact{}
	for _, r := range cfg.RPM {
		parts := strings.Split(r.Distro, "-")
		if len(parts) < 2 {
			panic(errInvalidRPMDistro)
		}

		rpms = append(rpms, distroArtifact{
			artifact: artifact{
				Architecture: r.Architecture,
				URL:          cfg.App.BaseURL + builders.GetPathForBranch(r.Path+"/"+parts[0]+"/"+parts[1], *branchID) + "/repodata/hydrapp.repo",
			},
			DistroName:    parts[0],
			DistroVersion: parts[1],
		})
	}

	debs := []distroArtifact{}
	for _, d := range cfg.DEB {
		debs = append(debs, distroArtifact{
			artifact: artifact{
				Architecture: d.Architecture,
				URL:          cfg.App.BaseURL + builders.GetPathForBranch(d.Path, *branchID),
			},
			DistroName:    d.OS,
			DistroVersion: d.Distro,
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
		MSIs:            msis,
		RPMs:            rpms,
		DEBs:            debs,
		BinariesURL:     cfg.App.BaseURL + builders.GetPathForBranch(cfg.Binaries.Path, *branchID),
	}); err != nil {
		panic(err)
	}

	fmt.Println(buf.String())
}
