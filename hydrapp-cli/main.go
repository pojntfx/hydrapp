package main

import "github.com/pojntfx/hydrapp/hydrapp-cli/cmd"

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
