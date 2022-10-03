package main

import (
	"context"
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
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/ncruces/zenity"
	"github.com/pojntfx/hydrapp/example/pkg/utils"
)

var (
	errNoAssetFound            = errors.New("no asset could be found")
	errNoEscalationMethodFound = errors.New("no escalation method could be found")
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

			binary := *appID + "." + runtime.GOOS
			switch runtime.GOOS {
			case "windows":
				binary += "-" + getArchIdentifier(runtime.GOARCH) + ".msi"
			case "darwin":
				binary += ".dmg"
			default:
				binary += "-" + getArchIdentifier(runtime.GOARCH)
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

				totalSize, err := strconv.Atoi(res.Header.Get("Content-Length"))
				if err != nil {
					panic(err)
				}

				dialog, err := zenity.Progress(
					zenity.Title("Downloading update"),
				)
				if err != nil {
					panic(err)
				}

				go func() {
					ticker := time.NewTicker(time.Millisecond * 50)
					defer func() {
						ticker.Stop()

						if err := dialog.Complete(); err != nil {
							panic(err)
						}

						if err := dialog.Close(); err != nil {
							panic(err)
						}
					}()

					for {
						select {
						case <-ctx.Done():

							return
						case <-ticker.C:
							stat, err := updatedExecutable.Stat()
							if err != nil {
								panic(err)
							}

							downloadedSize := stat.Size()
							if totalSize < 1 {
								downloadedSize = 1
							}

							percentage := int((float64(downloadedSize) / float64(totalSize)) * 100)

							if err := dialog.Value(percentage); err != nil {
								panic(err)
							}

							if err := dialog.Text(fmt.Sprintf("%v%% (%v MB/%v MB)", percentage, downloadedSize/(1024*1024), totalSize/(1024*1024))); err != nil {
								panic(err)
							}

							if percentage == 100 {
								return
							}
						}
					}
				}()

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
				if output, err := exec.Command("cmd.exe", "/C", "start", "/b", updatedExecutable.Name()).CombinedOutput(); err != nil {
					panic(fmt.Errorf("could not start update installer with output: %s: %v", output, err))
				}
			case "darwin":
				mountpoint, err := os.MkdirTemp(os.TempDir(), "update-mountpoint")
				if err != nil {
					panic(err)
				}
				defer os.RemoveAll(mountpoint)

				if output, err := exec.Command("hdiutil", "attach", "-mountpoint", mountpoint, updatedExecutable.Name()).CombinedOutput(); err != nil {
					panic(fmt.Errorf("could not attach DMG with output: %s: %v", output, err))
				}

				appPath, err := filepath.Abs(filepath.Join(oldExecutable, "..", ".."))
				if err != nil {
					panic(err)
				}

				appsPath, err := filepath.Abs(filepath.Join(appPath, ".."))
				if err != nil {
					panic(err)
				}

				if output, err := exec.Command("osascript", "-e", fmt.Sprintf(`do shell script "rm -rf \"%v\" && cp -r \"%v\"/* \"%v\"" with administrator privileges`, appPath, mountpoint, appsPath)).CombinedOutput(); err != nil {
					panic(fmt.Errorf("could not replace old app with new app with output: %s: %v", output, err))
				}

				if output, err := exec.Command("hdiutil", "unmount", mountpoint).CombinedOutput(); err != nil {
					panic(fmt.Errorf("could not detach DMG with output: %s: %v", output, err))
				}

				if err := utils.ForkExec(
					oldExecutable,
					os.Args,
				); err != nil {
					panic(err)
				}
			default:
				if err := os.Chmod(updatedExecutable.Name(), 0755); err != nil {
					panic(err)
				}

				// Escalate using Polkit
				if pkexec, err := exec.LookPath("pkexec"); err == nil {
					if output, err := exec.Command(pkexec, "cp", "-f", updatedExecutable.Name(), oldExecutable).CombinedOutput(); err != nil {
						panic(fmt.Errorf("could not install updated executable with output: %s: %v", output, err))
					}
				} else {
					// Escalate using using terminal emulator
					xterm, err := exec.LookPath("xterm")
					if err != nil {
						panic(fmt.Errorf("%v: %w", errNoEscalationMethodFound, err))
					}

					suid, err := exec.LookPath("sudo")
					if err != nil {
						suid, err = exec.LookPath("doas")
						if err != nil {
							panic(fmt.Errorf("%v: %w", errNoEscalationMethodFound, err))
						}
					}

					if output, err := exec.Command(
						xterm, "-T", "Authentication Required", "-e", fmt.Sprintf(`echo 'Authentication is needed to apply the update.' && %v cp -f '%v' '%v'`, suid, updatedExecutable.Name(), oldExecutable),
					).CombinedOutput(); err != nil {
						panic(fmt.Errorf("could not install updated executable with output: %s: %v", output, err))
					}
				}

				if err := utils.ForkExec(
					oldExecutable,
					os.Args,
				); err != nil {
					panic(err)
				}
			}

			os.Exit(0)
		}()
	}

	fmt.Println("Actual application logic goes here")
}
