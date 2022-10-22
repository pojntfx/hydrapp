//go:build android
// +build android

package main

/*
#include "hydrapp_android.h"
*/
import "C"
import (
	"log"
	"time"

	_ "github.com/pojntfx/hydrapp/hydrapp-utils/pkg/fixes"

	backend "github.com/pojntfx/hydrapp/hydrapp-example/pkg/backend"
	frontend "github.com/pojntfx/hydrapp/hydrapp-example/pkg/frontend"
)

//export Java_com_pojtinger_felicitas_hydrapp_example_MainActivity_LaunchBackend
func Java_com_pojtinger_felicitas_hydrapp_example_MainActivity_LaunchBackend(env *C.JNIEnv, activity C.jobject) C.jstring {
	backendURL, _, err := backend.StartServer("", time.Second*10, false)
	if err != nil {
		log.Fatalln("could not start backend:", err)
	}

	log.Println("Backend URL:", backendURL)

	frontendURL, _, err := frontend.StartServer("", backendURL, false)
	if err != nil {
		log.Fatalln("could not start frontend:", err)
	}

	log.Println("Frontend URL:", frontendURL)

	return C.get_java_string(env, C.CString(frontendURL))
}

func main() {}
