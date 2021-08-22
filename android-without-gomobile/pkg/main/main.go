//go:build android

package main

import "github.com/pojntfx/hydrapp/example/pkg/bindings"

//export show_toast
func ShowToast(msg string) error {
	return bindings.ShowToast(msg)
}

func main() {}
