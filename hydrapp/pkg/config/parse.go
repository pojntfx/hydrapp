package config

import (
	"io"

	"github.com/pojntfx/hydrapp/hydrapp/pkg/renderers"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/renderers/rpm"
	"gopkg.in/yaml.v3"
)

type Root struct {
	App      App                 `yaml:"app"`
	Go       Go                  `yaml:"go"`
	Releases []renderers.Release `yaml:"releases"`
	DEB      []DEB               `yaml:"deb"`
	DMG      DMG                 `yaml:"dmg"`
	Flatpak  []Flatpak           `yaml:"flatpak"`
	MSI      []MSI               `yaml:"msi"`
	RPM      []RPM               `yaml:"rpm"`
	APK      APK                 `yaml:"apk"`
	Binaries Binaries            `yaml:"binaries"`
	Docs     Docs                `yaml:"docs"`
}

type App struct {
	ID          string `yaml:"id"`
	Name        string `yaml:"name"`
	Summary     string `yaml:"summary"`
	Description string `yaml:"description"`
	License     string `yaml:"license"`
	Homepage    string `yaml:"homepage"`
	Git         string `yaml:"git"`
	BaseURL     string `yaml:"baseurl"`
}

type Go struct {
	Main     string `yaml:"main"`
	Flags    string `yaml:"flags"`
	Generate string `yaml:"generate"`
	Tests    string `yaml:"tests"`
	Image    string `yaml:"img"`
}

type DEB struct {
	Path            string        `yaml:"path"`
	OS              string        `yaml:"os"`
	Distro          string        `yaml:"distro"`
	Mirrorsite      string        `yaml:"mirrorsite"`
	Components      []string      `yaml:"components"`
	Debootstrapopts string        `yaml:"debootstrapopts"`
	Architecture    string        `yaml:"architecture"`
	Packages        []rpm.Package `yaml:"packages"`
}

type DMG struct {
	Path     string   `yaml:"path"`
	Packages []string `yaml:"packages"`
}

type Flatpak struct {
	Path         string        `yaml:"path"`
	Architecture string        `yaml:"architecture"`
	Packages     []rpm.Package `yaml:"packages"`
}

type MSI struct {
	Path         string   `yaml:"path"`
	Architecture string   `yaml:"architecture"`
	Include      string   `yaml:"include"`
	Packages     []string `yaml:"packages"`
}

type RPM struct {
	Path         string        `yaml:"path"`
	Trailer      string        `yaml:"trailer"`
	Distro       string        `yaml:"distro"`
	Architecture string        `yaml:"architecture"`
	Packages     []rpm.Package `yaml:"packages"`
}

type APK struct {
	Path string `yaml:"path"`
}

type Binaries struct {
	Path     string   `yaml:"path"`
	Exclude  string   `yaml:"exclude"`
	Packages []string `yaml:"packages"`
}

type Docs struct {
	Path string `yaml:"path"`
}

func Parse(r io.Reader) (*Root, error) {
	var root Root
	if err := yaml.NewDecoder(r).Decode(&root); err != nil {
		return nil, err
	}

	return &root, nil
}
