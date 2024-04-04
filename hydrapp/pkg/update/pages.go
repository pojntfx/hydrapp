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
	"sync"
	"syscall"
	"time"

	"github.com/ncruces/zenity"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/builders"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/config"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/utils"
)

type BrowserState struct {
	Cmd *exec.Cmd
}

var (
	ErrNoEscalationMethodFound = errors.New("no escalation method could be found")

	CommitTimeRFC3339 = ""
	BranchID          = ""
	PackageType       = ""
)

type File struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Time string `json:"time"`
}

type downloadConfiguration struct {
	description string
	url         string
	dst         *os.File
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

	currentBinaryBuildTime, err := time.Parse(time.RFC3339, CommitTimeRFC3339)
	if err != nil {
		handlePanic(cfg.App.Name, err.Error(), err)

		return
	}

	baseURL, err := url.Parse(cfg.App.BaseURL)
	if err != nil {
		handlePanic(cfg.App.Name, err.Error(), err)

		return
	}

	switch PackageType {
	case "dmg":
		baseURL.Path = path.Join(baseURL.Path, cfg.DMG.Path, BranchID)

	case "msi":
		for _, msiCfg := range cfg.MSI {
			if msiCfg.Architecture == runtime.GOARCH {
				baseURL.Path = path.Join(baseURL.Path, msiCfg.Path, BranchID)

				break
			}
		}

	default:
		baseURL.Path = path.Join(baseURL.Path, cfg.Binaries.Path, BranchID)
	}

	indexURL, err := url.JoinPath(baseURL.String(), "index.json")
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Get(indexURL)
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

	updatedBinaryName := ""
	switch PackageType {
	case "dmg":
		updatedBinaryName = builders.GetAppIDForBranch(cfg.App.ID, BranchID) + "." + runtime.GOOS + ".dmg"

	case "msi":
		updatedBinaryName = builders.GetAppIDForBranch(cfg.App.ID, BranchID) + "." + runtime.GOOS + "-" + utils.GetArchIdentifier(runtime.GOARCH) + ".msi"

	default:
		updatedBinaryName = builders.GetAppIDForBranch(cfg.App.ID, BranchID) + "." + runtime.GOOS + "-" + utils.GetArchIdentifier(runtime.GOARCH) + getBinIdentifier(runtime.GOOS, runtime.GOARCH)
	}

	var (
		updatedBinaryURL         = ""
		updatedBinaryReleaseTime time.Time

		updatedRepoKeyURL = ""

		updatedSignatureURL = ""
	)
	for _, file := range index {
		if file.Name == updatedBinaryName {
			updatedBinaryReleaseTime, err = time.Parse(time.RFC3339, file.Time)
			if err != nil {
				handlePanic(cfg.App.Name, err.Error(), err)

				return
			}

			if currentBinaryBuildTime.Before(updatedBinaryReleaseTime) {
				updatedBinaryURL, err = url.JoinPath(baseURL.String(), updatedBinaryName)
				if err != nil {
					handlePanic(cfg.App.Name, err.Error(), err)

					return
				}

				updatedRepoKeyURL, err = url.JoinPath(baseURL.String(), updatedBinaryName+".asc")
				if err != nil {
					handlePanic(cfg.App.Name, err.Error(), err)

					return
				}

				updatedSignatureURL, err = url.JoinPath(baseURL.String(), "repo.asc")
				if err != nil {
					handlePanic(cfg.App.Name, err.Error(), err)

					return
				}
			}

			break
		}
	}

	if strings.TrimSpace(updatedBinaryURL) == "" {
		return
	}

	if err := zenity.Question(
		fmt.Sprintf("Do you want to upgrade from version %v to %v now?", currentBinaryBuildTime, updatedBinaryReleaseTime),
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

	updatedBinaryFile, err := os.CreateTemp(os.TempDir(), updatedBinaryName)
	if err != nil {
		handlePanic(cfg.App.Name, err.Error(), err)

		return
	}
	defer os.Remove(updatedBinaryFile.Name())

	updatedSignatureFile, err := os.CreateTemp(os.TempDir(), updatedBinaryName+".asc")
	if err != nil {
		handlePanic(cfg.App.Name, err.Error(), err)

		return
	}
	defer os.Remove(updatedSignatureFile.Name())

	updatedRepoKeyFile, err := os.CreateTemp(os.TempDir(), "repo.asc")
	if err != nil {
		handlePanic(cfg.App.Name, err.Error(), err)

		return
	}
	defer os.Remove(updatedRepoKeyFile.Name())

	downloadConfigurations := []downloadConfiguration{
		{
			description: "Downloading binary",
			url:         updatedBinaryURL,
			dst:         updatedBinaryFile,
		},
		{
			description: "Downloading signature",
			url:         updatedSignatureURL,
			dst:         updatedSignatureFile,
		},
		{
			description: "Downloading repo key",
			url:         updatedRepoKeyURL,
			dst:         updatedRepoKeyFile,
		},
	}

	for _, downloadConfiguration := range downloadConfigurations {
		res, err := http.Get(downloadConfiguration.url)
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
			zenity.Title(downloadConfiguration.description),
		)
		if err != nil {
			handlePanic(cfg.App.Name, err.Error(), err)

			return
		}

		var dialogWg sync.WaitGroup
		dialogWg.Add(1)
		go func() {
			ticker := time.NewTicker(time.Millisecond * 50)
			defer func() {
				defer dialogWg.Done()

				ticker.Stop()

				if err := dialog.Complete(); err != nil {
					handlePanic(cfg.App.Name, "could not open download progress dialog", err)
				}

				if err := dialog.Close(); err != nil {
					handlePanic(cfg.App.Name, "could not close download progress dialog", err)
				}
			}()

			for {
				select {
				case <-ctx.Done():

					return
				case <-ticker.C:
					stat, err := downloadConfiguration.dst.Stat()
					if err != nil {
						handlePanic(cfg.App.Name, "could not get info on updated executable", err)
					}

					downloadedSize := stat.Size()
					if totalSize < 1 {
						downloadedSize = 1
					}

					percentage := int((float64(downloadedSize) / float64(totalSize)) * 100)

					if err := dialog.Value(percentage); err != nil {
						handlePanic(cfg.App.Name, "could not set update download progress percentage", err)
					}

					if err := dialog.Text(fmt.Sprintf("%v%% (%v MB/%v MB)", percentage, downloadedSize/(1024*1024), totalSize/(1024*1024))); err != nil {
						handlePanic(cfg.App.Name, "could not set update download progress description", err)
					}

					if percentage == 100 {
						return
					}
				}
			}
		}()

		if _, err := io.Copy(downloadConfiguration.dst, res.Body); err != nil {
			handlePanic(cfg.App.Name, err.Error(), err)

			return
		}

		dialogWg.Wait()
	}

	dialog, err := zenity.Progress(
		zenity.Title("Validating update"),
		zenity.Pulsate(),
	)
	if err != nil {
		handlePanic(cfg.App.Name, err.Error(), err)

		return
	}

	if err := dialog.Text("Reading repo key"); err != nil {
		handlePanic(cfg.App.Name, "could not set update validation progress description", err)
	}

	// TODO: Read repo public key

	if err := dialog.Text("Validating binary with signature and key"); err != nil {
		handlePanic(cfg.App.Name, "could not set update validation progress description", err)
	}

	// TODO: Validate binary with signature and key

	if err := dialog.Complete(); err != nil {
		handlePanic(cfg.App.Name, "could not open validation progress dialog", err)
	}

	if err := dialog.Close(); err != nil {
		handlePanic(cfg.App.Name, "could not close validation progress dialog", err)
	}

	oldExecutable, err := os.Executable()
	if err != nil {
		handlePanic(cfg.App.Name, err.Error(), err)

		return
	}

	switch PackageType {
	case "msi":
		if output, err := exec.Command("cmd.exe", "/C", "start", "/b", updatedBinaryFile.Name()).CombinedOutput(); err != nil {
			err := fmt.Errorf("could not start update installer with output: %s: %v", output, err)

			handlePanic(cfg.App.Name, err.Error(), err)

			return
		}

	case "dmg":
		mountpoint, err := os.MkdirTemp(os.TempDir(), "update-mountpoint")
		if err != nil {
			handlePanic(cfg.App.Name, err.Error(), err)

			return
		}
		defer os.RemoveAll(mountpoint)

		if output, err := exec.Command("hdiutil", "attach", "-mountpoint", mountpoint, updatedBinaryFile.Name()).CombinedOutput(); err != nil {
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
		switch runtime.GOOS {
		// TODO: Implement "windows" and "darwin" updaters

		default:
			if err := os.Chmod(updatedBinaryFile.Name(), 0755); err != nil {
				handlePanic(cfg.App.Name, err.Error(), err)

				return
			}

			// Escalate using Polkit
			if pkexec, err := exec.LookPath("pkexec"); err == nil {
				if output, err := exec.Command(pkexec, "cp", "-f", updatedBinaryFile.Name(), oldExecutable).CombinedOutput(); err != nil {
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
					xterm, "-T", "Authentication Required", "-e", fmt.Sprintf(`echo 'Authentication is needed to apply the update.' && %v cp -f '%v' '%v'`, suid, updatedBinaryFile.Name(), oldExecutable),
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
