//go:build android

package main

/*
#include "main.h"
*/
import "C"
import (
	"log"

	"github.com/pojntfx/multi-browser-electron/android-integrated-webserver/pkg/backend"
)

//export Java_com_pojtinger_felicitas_integratedWebserverExample_MainActivity_LaunchBackend
func Java_com_pojtinger_felicitas_integratedWebserverExample_MainActivity_LaunchBackend(env *C.JNIEnv, activity C.jobject) C.jstring {
	url, err := backend.StartServer()
	if err != nil {
		log.Fatalln("could not start integrated webserver:", err)
	}

	return C.get_java_string(env, C.CString(url))
}

func main() {}
