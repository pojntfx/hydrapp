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
	gpgKeyPassword := flag.String("gpg-key-password", "", " base64-encoded password for the GPG key")
	gpgKeyContent := flag.String("gpg-key-content", "", "base64-encoded GPG key contents")
	gpgKeyID := flag.String("gpg-key-id", "", "ID of the GPG key")
	baseURL := flag.String("base-url", "https://pojntfx.github.io/hydrapp/yum", "Base URL where the repo is to be hosted")

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

	dsc := filepath.Join(home, "rpmbuild", "SRPMS", *packageName+"-"+*suffix+".src.rpm")
	if output, err := exec.Command("rpmlint", dsc).CombinedOutput(); err != nil {
		panic(fmt.Errorf("could lint source package with output: %s :%w", output, err))
	}

	gpgKeyPasswordFile, err := os.CreateTemp(os.TempDir(), "hydrapp-gpg-key-password")
	if err != nil {
		panic(err)
	}
	defer os.Remove(gpgKeyPasswordFile.Name())

	if _, err := gpgKeyPasswordFile.WriteString(*gpgKeyPassword); err != nil {
		panic(err)
	}

	gpgDir := filepath.Join(home, ".gnupg")
	if err := os.MkdirAll(gpgDir, os.ModePerm); err != nil {
		panic(err)
	}

	if err := os.WriteFile(
		filepath.Join(gpgDir, "gpg.conf"),
		[]byte(
			fmt.Sprintf(
				`yes
passphrase-file %v
pinentry-mode loopback`,
				gpgKeyPasswordFile.Name(),
			),
		), os.ModePerm); err != nil {
		panic(err)
	}

	gpgKeyContentFile, err := os.CreateTemp(os.TempDir(), "hydrapp-gpg-key-content")
	if err != nil {
		panic(err)
	}
	defer os.Remove(gpgKeyContentFile.Name())

	if _, err := gpgKeyContentFile.WriteString(*gpgKeyContent); err != nil {
		panic(err)
	}

	if output, err := exec.Command("gpg", "--import", gpgKeyContentFile.Name()).CombinedOutput(); err != nil {
		panic(fmt.Errorf("could create source package with output: %s :%w", output, err))
	}

	if err := os.WriteFile(
		filepath.Join(home, ".rpmmacros"),
		[]byte(
			fmt.Sprintf(
				`%%_signature gpg
%%_gpg_name %v`,
				gpgKeyID,
			),
		), os.ModePerm); err != nil {
		panic(err)
	}

	for _, target := range targets {
		if len(target) < 2 {
			panic("could not work with invalid target definition")
		}

		distro := target[0]
		architectures := target[1:]

		for _, architecture := range architectures {
			if output, err := exec.Command("mock", "-r", distro+"-"+architecture, dsc, "--enable-network").CombinedOutput(); err != nil {
				panic(fmt.Errorf("could not create chroot with output: %s :%w", output, err))
			}

			if output, err := exec.Command("rpmlint", filepath.Join("/var", "lib", "mock", distro+"-"+architecture, "result", "*.rpm")).CombinedOutput(); err != nil {
				panic(fmt.Errorf("could not lint output package with output: %s :%w", output, err))
			}

			outDir := filepath.Join("out", "repositories", distro)
			if err := os.MkdirAll(outDir, os.ModePerm); err != nil {
				panic(err)
			}

			if output, err := exec.Command("cp", filepath.Join("/var", "lib", "mock", distro+"-"+architecture, "result", "*.rpm"), outDir).CombinedOutput(); err != nil {
				panic(fmt.Errorf("could not copy packages to output directory with output: %s :%w", output, err))
			}

			if output, err := exec.Command("rpm", "--addsign", filepath.Join(outDir, "*.rpm")).CombinedOutput(); err != nil {
				panic(fmt.Errorf("could not sign package with output: %s :%w", output, err))
			}

			if output, err := exec.Command("createrepo", outDir).CombinedOutput(); err != nil {
				panic(fmt.Errorf("could not create repo in out directory with output: %s :%w", output, err))
			}

			if output, err := exec.Command("gpg", "--detach-sign", "--armor", filepath.Join(outDir, "repodata", "repomd.xml")).CombinedOutput(); err != nil {
				panic(fmt.Errorf("could not sign repo with output: %s :%w", output, err))
			}

			if output, err := exec.Command("gpg", "--output", filepath.Join(outDir, "repodata", "repo.asc"), "--armor", "--export").CombinedOutput(); err != nil {
				panic(fmt.Errorf("could not add signature to repo with output: %s :%w", output, err))
			}

			if err := os.WriteFile(
				filepath.Join(outDir, "repodata", "hydrapp.repo"),
				[]byte(
					fmt.Sprintf(
						`[hydrapp-repo]
name=Hydrapp YUM repo
baseurl=%v/%v
enabled=1
gpgcheck=1
gpgkey=%v/%v/repodata/repo.asc`,
						*baseURL,
						distro,
						*baseURL,
						distro,
					),
				), os.ModePerm); err != nil {
				panic(err)
			}
		}
	}
}
