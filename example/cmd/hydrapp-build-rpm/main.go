package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	packageName := flag.String("package", "com.pojtinger.felicitas.hydrapp.example-0.0.1", "Name of the package to build")
	suffix := flag.String("suffix", "1.fc34", "Suffix of the package to build")
	spec := flag.String("spec", "com.pojtinger.felicitas.hydrapp.example.spec", "SPEC file to use as the recipe for the package")
	rawTargets := flag.String("targets", `[["epel-9", "x86_64"], ["fedora", "x86_64"], ["opensuse-tumbleweed", "x86_64"]]`, `Targets to build for (in format [["distro1", "arch1", "arch2"], ["distro2", "arch1"]]`)

	flag.Parse()

	var targets [][]string
	if err := json.Unmarshal([]byte(*rawTargets), &targets); err != nil {
		panic(err)
	}

	if output, err := exec.Command("dnf", "install", "-y", "fedora-packager", "@development-tools", "qemu-user-static", "rpm-sign").CombinedOutput(); err != nil {
		panic(fmt.Errorf("could not install dependencies with output: %s :%w", output, err))
	}

	if output, err := exec.Command("rpmdev-setuptree").CombinedOutput(); err != nil {
		panic(fmt.Errorf("could not setup RPM tree with output: %s :%w", output, err))
	}

	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	if output, err := exec.Command("tar", "-cvzf", filepath.Join(home, "rpmbuild", "SOURCES", *packageName+".tar.gz"), "--exclude", "out", "--transform", "s,^,"+*packageName+"/,", ".").CombinedOutput(); err != nil {
		panic(fmt.Errorf("could create tar archive with output: %s :%w", output, err))
	}

	if output, err := exec.Command("rpmbuild", "-bs", *spec).CombinedOutput(); err != nil {
		panic(fmt.Errorf("could create source package with output: %s :%w", output, err))
	}

	if output, err := exec.Command("rpmlint", filepath.Join(home, "rpmbuild", "SRPMS", *packageName+"-"+*suffix+".src.rpm")).CombinedOutput(); err != nil {
		panic(fmt.Errorf("could lint source package with output: %s :%w", output, err))
	}
}
