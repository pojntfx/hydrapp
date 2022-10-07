package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/client"
	"github.com/pojntfx/hydrapp/example/pkg/executors"
)

var (
	errCouldNotParseTargets = errors.New("could not parse targets")
	errInvalidTarget        = errors.New("invalid target")
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	image := flag.String("image", "ghcr.io/pojntfx/hydrapp-build-rpm", "OCI image to use")
	pull := flag.Bool("pull", true, "Whether to pull the image or not")
	dst := flag.String("dst", filepath.Join(pwd, "out", "rpm"), "Output directory")
	appID := flag.String("app-id", "com.pojtinger.felicitas.hydrapp.example", "RPM app ID to use")
	gpgKeyContent := flag.String("gpg-key-content", "", "base64-encoded GPG key contents")
	gpgKeyPassword := flag.String("gpg-key-password", "", " base64-encoded password for the GPG key")
	gpgKeyID := flag.String("gpg-key-id", "", "ID of the GPG key to use")
	baseURL := flag.String("base-url", "https://pojntfx.github.io/hydrapp/rpm", "Base URL where the repo is to be hosted")
	targetsFlag := flag.String("targets", `[["rocky+epel-9", "x86_64"], ["fedora-36", "x86_64"], ["opensuse-tumbleweed", "x86_64"]]`, `List of distros and architectures to build for (in JSON format [["distro1", "architectureN"...]...])`)
	packageVersion := flag.String("package-version", "0.0.1", "RPM package version")
	packageSuffix := flag.String("package-suffix", "1.fc36", "RPM package suffix")

	flag.Parse()

	var rawTargets [][]string
	if err := json.Unmarshal([]byte(*targetsFlag), &rawTargets); err != nil {
		panic(errCouldNotParseTargets)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	if err := executors.DockerRunImage(
		ctx,
		cli,
		*image,
		*pull,
		true,
		*dst,
		map[string]string{
			"APP_ID":           *appID,
			"GPG_KEY_CONTENT":  *gpgKeyContent,
			"GPG_KEY_PASSWORD": *gpgKeyPassword,
			"GPG_KEY_ID":       *gpgKeyID,
			"BASE_URL":         *baseURL,
			"TARGETS": func() string {
				targets := ""
				for i, rawTarget := range rawTargets {
					if len(rawTarget) < 2 {
						panic(errInvalidTarget)
					}

					if i != 0 {
						targets += " "
					}

					targets += fmt.Sprintf(`%v|"%v"`, rawTarget[0], strings.Join(rawTarget[1:], " "))
				}

				return targets
			}(),
			"PACKAGE_VERSION": *packageVersion,
			"PACKAGE_SUFFIX":  *packageSuffix,
		},
	); err != nil {
		panic(err)
	}
}
