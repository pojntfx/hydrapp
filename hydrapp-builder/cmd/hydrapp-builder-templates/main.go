package main

import (
	"flag"
	"fmt"

	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers/androidmanifest"
)

func main() {
	appID := flag.String("app-id", "com.pojtinger.felicitas.hydrapp.example", "App ID")
	appName := flag.String("app-name", "Hydrapp Example", "App name")

	flag.Parse()

	if path, content, err := androidmanifest.NewRenderer(
		*appID,
		*appName,
	).Render(); err != nil {
		panic(err)
	} else {
		fmt.Printf("%v\n%v\n", path, content)
	}
}
