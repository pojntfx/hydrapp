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
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/builders"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/config"
	"github.com/pojntfx/hydrapp/hydrapp-utils/pkg/utils"
)

type BrowserState struct {
	Cmd *exec.Cmd
}

const (
	goosMacOS   = "darwin"
	goosWindows = "windows"
)

var (
	ErrNoEscalationMethodFound = errors.New("no escalation method could be found")

	CommitTimeRFC3339 = ""
	BranchID          = ""
)

type File struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Time string `json:"time"`
}

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

func Update(
	ctx context.Context,

	cfg *config.Root,
	state *BrowserState,
	handlePanic func(appName, msg string, err error),
) {
	if strings.TrimSpace(CommitTimeRFC3339) == "" || strings.TrimSpace(BranchID) == "" || os.Getenv(utils.EnvSelfupdate) == "false" {
		return
	}

	buildtime, err := time.Parse(time.RFC3339, CommitTimeRFC3339)
	if err != nil {
		handlePanic(cfg.App.Name, err.Error(), err)

		return
	}

	baseURL, err := url.Parse(cfg.App.BaseURL)
	if err != nil {
		handlePanic(cfg.App.Name, err.Error(), err)

		return
	}
	baseURL.Path = path.Join(baseURL.Path, cfg.Binaries.Path, BranchID, "index.json")

	switch runtime.GOOS {
	case goosWindows:
		msiPath := ""
		for _, c := range cfg.MSI {
			if c.Architecture == runtime.GOOS {
				msiPath = c.Path

				break
			}
		}

		baseURL.Path = path.Join(baseURL.Path, msiPath, BranchID, "index.json")
	case goosMacOS:
		baseURL.Path = path.Join(baseURL.Path, cfg.DMG.Path, BranchID, "index.json")
	}

	res, err := http.DefaultClient.Get(baseURL.String())
	if err != nil {
		handlePanic(cfg.App.Name, err.Error(), err)

		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		handlePanic(cfg.App.Name, err.Error(), err)

		return
	}

	var index []File
	if err := json.Unmarshal(body, &index); err != nil {
		handlePanic(cfg.App.Name, err.Error(), err)

		return
	}

	binary := builders.GetAppIDForBranch(cfg.App.ID, BranchID) + "." + runtime.GOOS + "-" + utils.GetArchIdentifier(runtime.GOARCH) + getBinIdentifier(runtime.GOOS, runtime.GOARCH)
	switch runtime.GOOS {
	case goosWindows:
		binary = builders.GetAppIDForBranch(cfg.App.ID, BranchID) + "." + runtime.GOOS + "-" + utils.GetArchIdentifier(runtime.GOARCH) + ".msi"
	case goosMacOS:
		binary = builders.GetAppIDForBranch(cfg.App.ID, BranchID) + "." + runtime.GOOS + ".dmg"
	}

	var downloadURL *url.URL
	var releasetime time.Time
	for _, file := range index {
		if file.Name == binary {
			releasetime, err = time.Parse(time.RFC3339, file.Time)
			if err != nil {
				handlePanic(cfg.App.Name, err.Error(), err)

				return
			}

			if buildtime.Before(releasetime) {
				downloadURL, err = url.Parse(cfg.App.BaseURL)
				if err != nil {
					handlePanic(cfg.App.Name, err.Error(), err)

					return
				}
				downloadURL.Path = path.Join(downloadURL.Path, cfg.Binaries.Path, BranchID, file.Name)
			}

			break
		}
	}

	if downloadURL == nil {
		return
	}

	if err := zenity.Question(
		fmt.Sprintf("Do you want to upgrade from version %v to %v now?", buildtime, releasetime),
		zenity.Title("Update available"),
		zenity.OKLabel("Update now"),
		zenity.CancelLabel("Ask me next time"),
	); err != nil {
		if err == zenity.ErrCanceled {
			return
		}

		handlePanic(cfg.App.Name, err.Error(), err)

		return
	}

	updatedExecutable, err := ioutil.TempFile(os.TempDir(), binary)
	if err != nil {
		handlePanic(cfg.App.Name, err.Error(), err)

		return
	}
	defer os.Remove(updatedExecutable.Name())

	{
		res, err := http.Get(downloadURL.String())
		if err != nil {
			handlePanic(cfg.App.Name, err.Error(), err)

			return
		}
		if res.StatusCode != http.StatusOK {
			err := fmt.Errorf("%v", res.Status)

			handlePanic(cfg.App.Name, err.Error(), err)

			return
		}

		totalSize, err := strconv.Atoi(res.Header.Get("Content-Length"))
		if err != nil {
			handlePanic(cfg.App.Name, err.Error(), err)

			return
		}

		dialog, err := zenity.Progress(
			zenity.Title("Downloading update"),
		)
		if err != nil {
			handlePanic(cfg.App.Name, err.Error(), err)

			return
		}

		go func() {
			ticker := time.NewTicker(time.Millisecond * 50)
			defer func() {
				ticker.Stop()

				if err := dialog.Complete(); err != nil {
					handlePanic(cfg.App.Name, "could not open progress dialog", err)
				}

				if err := dialog.Close(); err != nil {
					handlePanic(cfg.App.Name, "could not close progress dialog", err)
				}
			}()

			for {
				select {
				case <-ctx.Done():

					return
				case <-ticker.C:
					stat, err := updatedExecutable.Stat()
					if err != nil {
						handlePanic(cfg.App.Name, "could not get info on updated executable", err)
					}

					downloadedSize := stat.Size()
					if totalSize < 1 {
						downloadedSize = 1
					}

					percentage := int((float64(downloadedSize) / float64(totalSize)) * 100)

					if err := dialog.Value(percentage); err != nil {
						handlePanic(cfg.App.Name, "could not set update progress percentage", err)
					}

					if err := dialog.Text(fmt.Sprintf("%v%% (%v MB/%v MB)", percentage, downloadedSize/(1024*1024), totalSize/(1024*1024))); err != nil {
						handlePanic(cfg.App.Name, "could not set update progress description", err)
					}

					if percentage == 100 {
						return
					}
				}
			}
		}()

		if _, err := io.Copy(updatedExecutable, res.Body); err != nil {
			handlePanic(cfg.App.Name, err.Error(), err)

			return
		}
	}

	oldExecutable, err := os.Executable()
	if err != nil {
		handlePanic(cfg.App.Name, err.Error(), err)

		return
	}

	switch runtime.GOOS {
	case "windows":
		if output, err := exec.Command("cmd.exe", "/C", "start", "/b", updatedExecutable.Name()).CombinedOutput(); err != nil {
			err := fmt.Errorf("could not start update installer with output: %s: %v", output, err)

			handlePanic(cfg.App.Name, err.Error(), err)

			return
		}
	case "darwin":
		mountpoint, err := os.MkdirTemp(os.TempDir(), "update-mountpoint")
		if err != nil {
			handlePanic(cfg.App.Name, err.Error(), err)

			return
		}
		defer os.RemoveAll(mountpoint)

		if output, err := exec.Command("hdiutil", "attach", "-mountpoint", mountpoint, updatedExecutable.Name()).CombinedOutput(); err != nil {
			err := fmt.Errorf("could not attach DMG with output: %s: %v", output, err)

			handlePanic(cfg.App.Name, err.Error(), err)

			return
		}

		appPath, err := filepath.Abs(filepath.Join(oldExecutable, "..", ".."))
		if err != nil {
			handlePanic(cfg.App.Name, err.Error(), err)

			return
		}

		appsPath, err := filepath.Abs(filepath.Join(appPath, ".."))
		if err != nil {
			handlePanic(cfg.App.Name, err.Error(), err)

			return
		}

		if output, err := exec.Command("osascript", "-e", fmt.Sprintf(`do shell script "rm -rf \"%v\" && cp -r \"%v\"/* \"%v\"" with administrator privileges`, appPath, mountpoint, appsPath)).CombinedOutput(); err != nil {
			err := fmt.Errorf("could not replace old app with new app with output: %s: %v", output, err)

			handlePanic(cfg.App.Name, err.Error(), err)

			return
		}

		if output, err := exec.Command("hdiutil", "unmount", mountpoint).CombinedOutput(); err != nil {
			err := fmt.Errorf("could not detach DMG with output: %s: %v", output, err)

			handlePanic(cfg.App.Name, err.Error(), err)

			return
		}

		if err := utils.ForkExec(
			oldExecutable,
			os.Args,
		); err != nil {
			handlePanic(cfg.App.Name, err.Error(), err)

			return
		}
	default:
		if err := os.Chmod(updatedExecutable.Name(), 0755); err != nil {
			handlePanic(cfg.App.Name, err.Error(), err)

			return
		}

		// Escalate using Polkit
		if pkexec, err := exec.LookPath("pkexec"); err == nil {
			if output, err := exec.Command(pkexec, "cp", "-f", updatedExecutable.Name(), oldExecutable).CombinedOutput(); err != nil {
				err := fmt.Errorf("could not install updated executable with output: %s: %v", output, err)

				handlePanic(cfg.App.Name, err.Error(), err)

				return
			}
		} else {
			// Escalate using using terminal emulator
			xterm, err := exec.LookPath("xterm")
			if err != nil {
				err := fmt.Errorf("%v: %w", ErrNoEscalationMethodFound, err)

				handlePanic(cfg.App.Name, err.Error(), err)

				return
			}

			suid, err := exec.LookPath("sudo")
			if err != nil {
				suid, err = exec.LookPath("doas")
				if err != nil {
					err := fmt.Errorf("%v: %w", ErrNoEscalationMethodFound, err)

					handlePanic(cfg.App.Name, err.Error(), err)

					return
				}
			}

			if output, err := exec.Command(
				xterm, "-T", "Authentication Required", "-e", fmt.Sprintf(`echo 'Authentication is needed to apply the update.' && %v cp -f '%v' '%v'`, suid, updatedExecutable.Name(), oldExecutable),
			).CombinedOutput(); err != nil {
				err := fmt.Errorf("could not install updated executable with output: %s: %v", output, err)

				handlePanic(cfg.App.Name, err.Error(), err)

				return
			}
		}

		if err := utils.ForkExec(
			oldExecutable,
			os.Args,
		); err != nil {
			handlePanic(cfg.App.Name, err.Error(), err)

			return
		}
	}

	if state != nil && state.Cmd != nil && state.Cmd.Process != nil {
		// Windows does not support the `SIGTERM` signal
		if runtime.GOOS == "windows" {
			if output, err := exec.Command("taskkill", "/pid", strconv.Itoa(state.Cmd.Process.Pid)).CombinedOutput(); err != nil {
				err := fmt.Errorf("could not close old version: %v: %v", string(output), err)

				handlePanic(cfg.App.Name, err.Error(), err)

				return
			}
		} else {
			// We ignore errors here as the old process might already have finished etc.
			_ = state.Cmd.Process.Signal(syscall.SIGTERM)
		}
	}

	os.Exit(0)
}
