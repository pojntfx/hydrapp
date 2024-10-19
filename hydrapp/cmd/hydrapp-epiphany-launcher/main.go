package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func main() {
	appID := strings.Replace(uuid.NewString(), "-", "", -1)

	dataHomeDir := os.Getenv("XDG_DATA_HOME")
	if strings.TrimSpace(dataHomeDir) == "" {
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}

		dataHomeDir = filepath.Join(userHomeDir, ".local", "share")
	}

	epiphanyID := "org.gnome.Epiphany.WebApp_" + appID

	applicationsDir := filepath.Join(dataHomeDir, "applications")
	if err := os.MkdirAll(applicationsDir, 0755); err != nil {
		panic(err)
	}

	profileDir := filepath.Join(dataHomeDir, epiphanyID)

	desktopFilePath := filepath.Join(applicationsDir, epiphanyID+".desktop")
	if err := os.WriteFile(desktopFilePath, []byte(fmt.Sprintf(`[Desktop Entry]
Exec=epiphany --new-window --application-mode --profile=%v https://news.ycombinator.com/
StartupNotify=true
Terminal=false
Type=Application
Categories=GNOME;GTK;
StartupWMClass=%v
X-Purism-FormFactor=Workstation;Mobile;
Name=Hacker News
Icon=com.pojtinger.felicitas.connmapper.main
NoDisplay=true`, profileDir, epiphanyID)), 0664); err != nil {
		panic(err)
	}
	defer os.RemoveAll(desktopFilePath)

	if err := os.MkdirAll(filepath.Join(profileDir, ".app"), 0755); err != nil {
		panic(err)
	}
	defer os.RemoveAll(profileDir)

	xdgApplicationsDir := filepath.Join(dataHomeDir, "xdg-desktop-portal", "applications")
	if err := os.MkdirAll(xdgApplicationsDir, 0755); err != nil {
		panic(err)
	}

	xdgApplicationsDesktopFilePath := filepath.Join(xdgApplicationsDir, epiphanyID+".desktop")
	if err := os.Symlink(desktopFilePath, xdgApplicationsDesktopFilePath); err != nil {
		panic(err)
	}
	defer os.RemoveAll(xdgApplicationsDesktopFilePath)

	epiphanyCmd := exec.Command("epiphany", "--new-window", "--application-mode", "--profile="+profileDir, "https://news.ycombinator.com/")
	epiphanyCmd.Stdout = os.Stdout
	epiphanyCmd.Stderr = os.Stderr

	if err := epiphanyCmd.Run(); err != nil {
		panic(err)
	}
}
