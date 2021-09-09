package main

import "github.com/ncruces/zenity"

func main() {
	if err := zenity.Info(
		"Hello, world!",
		zenity.Title("hydrapp example MSI"),
		zenity.InfoIcon,
	); err != nil {
		panic(err)
	}
}
