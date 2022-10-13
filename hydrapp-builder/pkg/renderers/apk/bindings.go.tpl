//go:build android
// +build android

package main

/*
#include "main_android.h"
*/
import "C"
import (
	"log"
	"time"

	_ "github.com/pojntfx/hydrapp/hydrapp-example/pkg/fixes"

	backend "{{ .AppBackendPackage }}"
	frontend "{{ .AppFrontendPackage }}"
)

//export Java_{{ .AppID }}_MainActivity_LaunchBackend
func Java_{{ .AppID }}_MainActivity_LaunchBackend(env *C.JNIEnv, activity C.jobject) C.jstring {
	backendURL, _, err := backend.StartServer("", time.Second*10)
	if err != nil {
		log.Fatalln("could not start backend:", err)
	}

	frontendURL, _, err := frontend.StartServer("", backendURL)
	if err != nil {
		log.Fatalln("could not start frontend:", err)
	}

	return C.get_java_string(env, C.CString(frontendURL))
}

func main() {}
