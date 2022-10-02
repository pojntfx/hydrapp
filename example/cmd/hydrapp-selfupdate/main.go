package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"syscall"

	"github.com/ncruces/zenity"
)

var (
	errNoAssetFound = errors.New("no asset could be found")
)

type release struct {
	Name   string  `json:"name"`
	Assets []asset `json:"assets"`
}

type asset struct {
	Name string `json:"name"`
	URL  string `json:"browser_download_url"`
}

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

func main() {
	update := flag.Bool("update", true, "Whether to check for updates before starting")
	apiURL := flag.String("api-url", "https://api.github.com/", "GitHub/Gitea API endpoint to use")
	owner := flag.String("owner", "pojntfx", "Repository owner")
	repo := flag.String("repo", "htorrent", "Repository name")
	current := flag.String("current", "v0.0.1", "Current version")
	appID := flag.String("app-id", "htorrent", "App ID to update")

	flag.Parse()

	if *update {
		func() {
			var rel release
			{
				u, err := url.Parse(*apiURL)
				if err != nil {
					panic(err)
				}

				u.Path = path.Join("repos", *owner, *repo, "releases", "latest")

				res, err := http.Get(u.String())
				if err != nil {
					panic(err)
				}
				if res.StatusCode != http.StatusOK {
					panic(res.Status)
				}

				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					panic(err)
				}

				if err := json.Unmarshal(body, &rel); err != nil {
					panic(err)
				}
			}

			if rel.Name == *current {
				return
			}

			if err := zenity.Question(
				fmt.Sprintf("Do you want to upgrade from version %v to %v now?", *current, rel.Name),
				zenity.Title("Update available"),
				zenity.OKLabel("Update now"),
				zenity.CancelLabel("Ask me next time"),
			); err != nil {
				if err == zenity.ErrCanceled {
					return
				}

				panic(err)
			}

			binary := *appID + "." + runtime.GOOS + "-"
			switch runtime.GOOS {
			case "windows":
				binary += getArchIdentifier(runtime.GOARCH) + ".msi"
			case "darwin":
				binary += getArchIdentifier(runtime.GOARCH) + ".dmg"
			default:
				binary += getArchIdentifier(runtime.GOARCH)
			}

			downloadURL := ""
			for _, asset := range rel.Assets {
				if asset.Name == binary {
					downloadURL = asset.URL

					break
				}
			}

			if strings.TrimSpace(downloadURL) == "" {
				panic(errNoAssetFound)
			}

			updatedExecutable, err := ioutil.TempFile(os.TempDir(), binary)
			if err != nil {
				panic(err)
			}
			defer os.Remove(updatedExecutable.Name())

			{
				u, err := url.Parse(downloadURL)
				if err != nil {
					panic(err)
				}

				res, err := http.Get(u.String())
				if err != nil {
					panic(err)
				}
				if res.StatusCode != http.StatusOK {
					panic(res.Status)
				}

				if _, err := io.Copy(updatedExecutable, res.Body); err != nil {
					panic(err)
				}
			}

			oldExecutable, err := os.Executable()
			if err != nil {
				panic(err)
			}

			switch runtime.GOOS {
			case "windows":
				// TODO: Add Windows support
				// 1. Execute MSI
				// 2. Kill self (MSI launches updated app after installation automatically)
			case "darwin":
				// TODO: Add macOS support
				// 1. Mount DMG
				// 2. sudo rm -rf currentExecutable/../..
				// 3. cp /Volumes/${VOLUMENAME} currentExecutable/../..
				// 4. Unmount DMG
				// 5. forkExec(currentExecutable)
			default:
				if err := os.Chmod(updatedExecutable.Name(), 0755); err != nil {
					panic(err)
				}

				// TODO: Add UNIX support if `pkexec` is not in path by spawning ${TERM} and running the same command using `sudo` or `doas`
				if output, err := exec.Command("pkexec", "cp", "-f", updatedExecutable.Name(), oldExecutable).CombinedOutput(); err != nil {
					panic(fmt.Errorf("could not install updated executable with output: %s: %v", output, err))
				}

				if _, err := syscall.ForkExec(
					oldExecutable,
					os.Args,
					&syscall.ProcAttr{
						Env:   os.Environ(),
						Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
					},
				); err != nil {
					panic(err)
				}
			}

			os.Exit(0)
		}()
	}

	fmt.Println("Actual application logic goes here")
}
