//go:build android
// +build android

package main

/*
#include "main_android.h"
*/
import "C"
import (
	"log"

	"github.com/pojntfx/multi-browser-electron/desktop-integrated-webserver-rpm/pkg/backend"
	_ "github.com/pojntfx/multi-browser-electron/desktop-integrated-webserver-rpm/pkg/fixes"
)

//export Java_com_pojtinger_felicitas_integratedWebserverExampleRPM_MainActivity_LaunchBackend
func Java_com_pojtinger_felicitas_integratedWebserverExampleRPM_MainActivity_LaunchBackend(env *C.JNIEnv, activity C.jobject) C.jstring {
	url, _, err := backend.StartServer()
	if err != nil {
		log.Fatalln("could not start integrated webserver:", err)
	}

	return C.get_java_string(env, C.CString(url))
}

func main() {}
