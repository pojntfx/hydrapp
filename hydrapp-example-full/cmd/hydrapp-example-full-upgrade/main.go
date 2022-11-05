package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"runtime"
	"time"

	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/builders"
	cconfig "github.com/pojntfx/hydrapp/hydrapp-builder/pkg/config"
)

// See https://github.com/pojntfx/bagop/blob/main/main.go#L33
func getBinIdentifier(goOS, goArch string) string {
	if goOS == "windows" {
		return ".exe"
	}

	if goOS == "js" && goArch == "wasm" {
		return ".wasm"
	}

	return ""
}

// See https://github.com/pojntfx/bagop/blob/main/main.go#L45
func getArchIdentifier(goArch string) string {
	switch goArch {
	case "386":
		return "i686"
	case "amd64":
		return "x86_64"
	case "arm":
		return "armv7l" // Best effort, could also be `armv6l` etc. depending on `GOARCH`
	case "arm64":
		return "aarch64"
	default:
		return goArch
	}
}

type File struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Time string `json:"time"`
}

func main() {
	config := flag.String("config", "hydrapp.yaml", "Config file to use")
	branchID := flag.String("branch-id", "", `Branch ID to fetch the app for, i.e. main (for an app ID like "myappid.main" and baseURL like "mybaseurl/main"`)
	buildtime := flag.String("build-time", "2022-11-04T19:16:18+01:00", "Build timestamp of the current binary")

	flag.Parse()

	bt, err := time.Parse(time.RFC3339, *buildtime)
	if err != nil {
		panic(err)
	}

	content, err := ioutil.ReadFile(*config)
	if err != nil {
		panic(err)
	}

	cfg, err := cconfig.Parse(content)
	if err != nil {
		panic(err)
	}

	u, err := url.Parse(cfg.App.BaseURL)
	if err != nil {
		panic(err)
	}
	u.Path = path.Join(u.Path, cfg.Binaries.Path, *branchID, "index.json")

	res, err := http.DefaultClient.Get(u.String())
	if err != nil {
		panic(err)
	}

	var index []File
	if err := json.NewDecoder(res.Body).Decode(&index); err != nil {
		panic(err)
	}

	targetBinaryName := builders.GetAppIDForBranch(cfg.App.ID, *branchID) + "." + runtime.GOOS + "-" + getArchIdentifier(runtime.GOARCH) + getBinIdentifier(runtime.GOOS, runtime.GOARCH)
	downloadURL := ""

	for _, file := range index {
		if file.Name == targetBinaryName {
			nt, err := time.Parse(time.RFC3339, file.Time)
			if err != nil {
				panic(err)
			}

			if bt.Before(nt) {
				dlu, err := url.Parse(cfg.App.BaseURL)
				if err != nil {
					panic(err)
				}
				dlu.Path = path.Join(dlu.Path, cfg.Binaries.Path, *branchID, file.Name)

				downloadURL = dlu.String()
			}

			break
		}
	}

	fmt.Println(downloadURL)
}
