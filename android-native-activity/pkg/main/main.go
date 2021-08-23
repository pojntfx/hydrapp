package main

import "C"
import "log"

//export ANativeActivity_onCreate
func ANativeActivity_onCreate() {
	log.Fatalln("Hello, world!")
}

func main() {}
