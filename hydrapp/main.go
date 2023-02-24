package main

import "github.com/pojntfx/hydrapp/hydrapp/cmd"

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
