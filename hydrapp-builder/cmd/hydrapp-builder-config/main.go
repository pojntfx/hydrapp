package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v3"
)

type Root struct {
	App      App       `yaml:"app"`
	License  License   `yaml:"license"`
	Releases []Release `yaml:"releases"`
	DEB      []DEB     `yaml:"deb"`
	DMG      DMG       `yaml:"dmg"`
	Flatpak  []Flatpak `yaml:"flatpak"`
	MSI      []MSI     `yaml:"msi"`
	RPM      []RPM     `yaml:"rpm"`
	APK      APK       `yaml:"apk"`
}

type App struct {
	ID          string `yaml:"id"`
	Name        string `yaml:"name"`
	Summary     string `yaml:"summary"`
	Description string `yaml:"description"`
	Homepage    string `yaml:"homepage"`
	Git         string `yaml:"git"`
	BaseURL     string `yaml:"baseurl"`
}

type License struct {
	SPDX string `yaml:"spdx"`
	Text string `yaml:"text"`
}

type Release struct {
	Version     string    `yaml:"version"`
	Date        time.Time `yaml:"date"`
	Description string    `yaml:"description"`
	Author      string    `yaml:"author"`
	Email       string    `yaml:"email"`
}

type DEB struct {
	Suffix          string    `yaml:"suffix"`
	OS              string    `yaml:"os"`
	Distro          string    `yaml:"distro"`
	Mirrorsite      string    `yaml:"mirrorsite"`
	Components      []string  `yaml:"components"`
	Debootstrapopts string    `yaml:"debootstrapopts"`
	Architecture    string    `yaml:"architecture"`
	Packages        []Package `yaml:"packages"`
}

type Package struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type DMG struct {
	Path      string   `yaml:"path"`
	Universal bool     `yaml:"universal"`
	Packages  []string `yaml:"packages"`
}

type Flatpak struct {
	Path         string `yaml:"path"`
	Architecture string `yaml:"architecture"`
}

type MSI struct {
	Path         string   `yaml:"path"`
	Architecture string   `yaml:"architecture"`
	Packages     []string `yaml:"packages"`
}

type RPM struct {
	Path         string    `yaml:"path"`
	Trailer      string    `yaml:"trailer"`
	Distro       string    `yaml:"distro"`
	Architecture string    `yaml:"architecture"`
	Packages     []Package `yaml:"packages"`
}

type APK struct {
	Path string `yaml:"path"`
}

func main() {
	config := flag.String("config", "hydrapp.yaml", "Config file to use")

	flag.Parse()

	content, err := ioutil.ReadFile(*config)
	if err != nil {
		panic(err)
	}

	var root Root
	if err := yaml.Unmarshal(content, &root); err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", root)
}
