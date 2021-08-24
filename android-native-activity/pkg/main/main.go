//go:build android

package main

/*
#include "main.h"
*/
import "C"

var (
	events = make(chan func(env *C.JNIEnv, activity C.jobject))
)

//export GoLoop
func GoLoop(env *C.JNIEnv, activity C.jobject) {
	go main()

	for event := range events {
		event(env, activity)
	}
}

func Queue(event func(env *C.JNIEnv, activity C.jobject)) {
	events <- event
}

func main() {
	Queue(func(env *C.JNIEnv, activity C.jobject) {
		C.show_toast(env, activity)
	})

	Queue(func(env *C.JNIEnv, activity C.jobject) {
		close(events)
	})
}
