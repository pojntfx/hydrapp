//go:build windows

package utils

import (
	"os/exec"
)

func ForkExec(path string, args []string) error {
	return exec.Command(path, args...).Start()
}
