package main

import (
	"github.com/pojntfx/multi-browser-electron/gomobile-webview-jni/pkg/bindings"
	"golang.org/x/mobile/app"
)

func main() {
	app.Main(func(a app.App) {
		bindings.ShowToast()
	})
}
