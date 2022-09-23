package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func runCommand(cmd *exec.Cmd, description string) error {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("could not %v: %w", description, err)
	}

	return nil
}

func main() {
	gpgKeyPassword := flag.String("gpg-key-password", "", " base64-encoded password for the GPG key")
	gpgKeyContent := flag.String("gpg-key-content", "", "base64-encoded GPG key contents")
	// appID := flag.String("app-id", "com.pojtinger.felicitas.hydrapp.example", "Android app ID to use")

	flag.Parse()

	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
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

	if err := runCommand(exec.Command("gpg", "--import", gpgKeyContentFile.Name()), "import GPG key"); err != nil {
		panic(err)
	}

	tmpDir, err := os.MkdirTemp(os.TempDir(), "hydrapp-tempdir")
	if err != nil {
		panic(err)
	}

	if err := os.MkdirAll(filepath.Join(tmpDir, "drawable"), os.ModePerm); err != nil {
		panic(err)
	}
}
