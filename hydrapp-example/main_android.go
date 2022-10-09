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

	"github.com/pojntfx/hydrapp/hydrapp-example/pkg/backend"
	_ "github.com/pojntfx/hydrapp/hydrapp-example/pkg/fixes"
	"github.com/pojntfx/hydrapp/hydrapp-example/pkg/frontend"
)

//export Java_com_pojtinger_felicitas_hydrapp_example_MainActivity_LaunchBackend
func Java_com_pojtinger_felicitas_hydrapp_example_MainActivity_LaunchBackend(env *C.JNIEnv, activity C.jobject) C.jstring {
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
