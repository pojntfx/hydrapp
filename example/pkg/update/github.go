//go:build selfupdate
// +build selfupdate

package update

import (
	"context"
	"encoding/json"
	"errors"
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
	"syscall"
	"time"

	"github.com/ncruces/zenity"
	"github.com/pojntfx/hydrapp/example/pkg/utils"
)

var (
	ErrNoAssetFound            = errors.New("no asset could be found")
	ErrNoEscalationMethodFound = errors.New("no escalation method could be found")
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

func Update(
	ctx context.Context,

	apiURL,
	owner,
	repo,

	currentVersion,
	appID string,

	state *BrowserState,
	handlePanic func(msg string, err error),
) error {
	var rel release
	{
		u, err := url.Parse(apiURL)
		if err != nil {
			return err
		}

		u.Path = path.Join("repos", owner, repo, "releases", "latest")

		res, err := http.Get(u.String())
		if err != nil {
			return err
		}
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("%v", res.Status)
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(body, &rel); err != nil {
			return err
		}
	}

	if rel.Name == currentVersion {
		return nil
	}

	if err := zenity.Question(
		fmt.Sprintf("Do you want to upgrade from version %v to %v now?", currentVersion, rel.Name),
		zenity.Title("Update available"),
		zenity.OKLabel("Update now"),
		zenity.CancelLabel("Ask me next time"),
	); err != nil {
		if err == zenity.ErrCanceled {
			return nil
		}

		return err
	}

	binary := appID + "." + runtime.GOOS
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
		return ErrNoAssetFound
	}

	updatedExecutable, err := ioutil.TempFile(os.TempDir(), binary)
	if err != nil {
		return err
	}
	defer os.Remove(updatedExecutable.Name())

	{
		u, err := url.Parse(downloadURL)
		if err != nil {
			return err
		}

		res, err := http.Get(u.String())
		if err != nil {
			return err
		}
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("%v", res.Status)
		}

		totalSize, err := strconv.Atoi(res.Header.Get("Content-Length"))
		if err != nil {
			return err
		}

		dialog, err := zenity.Progress(
			zenity.Title("Downloading update"),
		)
		if err != nil {
			return err
		}

		go func() {
			ticker := time.NewTicker(time.Millisecond * 50)
			defer func() {
				ticker.Stop()

				if err := dialog.Complete(); err != nil {
					handlePanic("could not open progress dialog", err)
				}

				if err := dialog.Close(); err != nil {
					handlePanic("could not close progress dialog", err)
				}
			}()

			for {
				select {
				case <-ctx.Done():

					return
				case <-ticker.C:
					stat, err := updatedExecutable.Stat()
					if err != nil {
						handlePanic("could not get info on updated executable", err)
					}

					downloadedSize := stat.Size()
					if totalSize < 1 {
						downloadedSize = 1
					}

					percentage := int((float64(downloadedSize) / float64(totalSize)) * 100)

					if err := dialog.Value(percentage); err != nil {
						handlePanic("could not set update progress percentage", err)
					}

					if err := dialog.Text(fmt.Sprintf("%v%% (%v MB/%v MB)", percentage, downloadedSize/(1024*1024), totalSize/(1024*1024))); err != nil {
						handlePanic("could not set update progress description", err)
					}

					if percentage == 100 {
						return
					}
				}
			}
		}()

		if _, err := io.Copy(updatedExecutable, res.Body); err != nil {
			return err
		}
	}

	oldExecutable, err := os.Executable()
	if err != nil {
		return err
	}

	switch runtime.GOOS {
	case "windows":
		if output, err := exec.Command("cmd.exe", "/C", "start", "/b", updatedExecutable.Name()).CombinedOutput(); err != nil {
			return fmt.Errorf("could not start update installer with output: %s: %v", output, err)
		}
	case "darwin":
		mountpoint, err := os.MkdirTemp(os.TempDir(), "update-mountpoint")
		if err != nil {
			return err
		}
		defer os.RemoveAll(mountpoint)

		if output, err := exec.Command("hdiutil", "attach", "-mountpoint", mountpoint, updatedExecutable.Name()).CombinedOutput(); err != nil {
			return fmt.Errorf("could not attach DMG with output: %s: %v", output, err)
		}

		appPath, err := filepath.Abs(filepath.Join(oldExecutable, "..", ".."))
		if err != nil {
			return err
		}

		appsPath, err := filepath.Abs(filepath.Join(appPath, ".."))
		if err != nil {
			return err
		}

		if output, err := exec.Command("osascript", "-e", fmt.Sprintf(`do shell script "rm -rf \"%v\" && cp -r \"%v\"/* \"%v\"" with administrator privileges`, appPath, mountpoint, appsPath)).CombinedOutput(); err != nil {
			return fmt.Errorf("could not replace old app with new app with output: %s: %v", output, err)
		}

		if output, err := exec.Command("hdiutil", "unmount", mountpoint).CombinedOutput(); err != nil {
			return fmt.Errorf("could not detach DMG with output: %s: %v", output, err)
		}

		if err := utils.ForkExec(
			oldExecutable,
			os.Args,
		); err != nil {
			return err
		}
	default:
		if err := os.Chmod(updatedExecutable.Name(), 0755); err != nil {
			return err
		}

		// Escalate using Polkit
		if pkexec, err := exec.LookPath("pkexec"); err == nil {
			if output, err := exec.Command(pkexec, "cp", "-f", updatedExecutable.Name(), oldExecutable).CombinedOutput(); err != nil {
				return fmt.Errorf("could not install updated executable with output: %s: %v", output, err)
			}
		} else {
			// Escalate using using terminal emulator
			xterm, err := exec.LookPath("xterm")
			if err != nil {
				return fmt.Errorf("%v: %w", ErrNoEscalationMethodFound, err)
			}

			suid, err := exec.LookPath("sudo")
			if err != nil {
				suid, err = exec.LookPath("doas")
				if err != nil {
					return fmt.Errorf("%v: %w", ErrNoEscalationMethodFound, err)
				}
			}

			if output, err := exec.Command(
				xterm, "-T", "Authentication Required", "-e", fmt.Sprintf(`echo 'Authentication is needed to apply the update.' && %v cp -f '%v' '%v'`, suid, updatedExecutable.Name(), oldExecutable),
			).CombinedOutput(); err != nil {
				return fmt.Errorf("could not install updated executable with output: %s: %v", output, err)
			}
		}

		if err := utils.ForkExec(
			oldExecutable,
			os.Args,
		); err != nil {
			return err
		}
	}

	if state.Cmd != nil && state.Cmd.Process != nil {
		// Windows does not support the `SIGTERM` signal
		if runtime.GOOS == "windows" {
			if output, err := exec.Command("taskkill", "/pid", strconv.Itoa(state.Cmd.Process.Pid)).CombinedOutput(); err != nil {
				return fmt.Errorf("could not close old version: %v: %v", string(output), err)
			}
		} else {
			// We ignore errors here as the old process might already have finished etc.
			_ = state.Cmd.Process.Signal(syscall.SIGTERM)
		}
	}

	os.Exit(0)

	return nil
}
