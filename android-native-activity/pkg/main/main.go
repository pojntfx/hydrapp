//go:build android

package main

/*
#include "main.h"
*/
import "C"

var (
	events = make(chan func(*C.ANativeActivity))
)

//export GoLoop
func GoLoop(activity *C.ANativeActivity) {
	go main()

	for event := range events {
		event(activity)
	}
}

func Queue(event func(*C.ANativeActivity)) {
	events <- event
}

func main() {
	Queue(func(aa *C.ANativeActivity) {
		panic("Hey!")
	})
}
