//go:build android
// +build android

package main

/*
#include "hydrapp_android.h"
*/
import "C"
import (
	"context"
	"log"

	_ "github.com/pojntfx/hydrapp/hydrapp/pkg/fixes"

	backend "github.com/pojntfx/hydrapp/hydrapp-example-rest/pkg/backend"
	frontend "github.com/pojntfx/hydrapp/hydrapp-example-rest/pkg/frontend"
)

//export Java_com_pojtinger_felicitas_hydrapp_example_rest_MainActivity_LaunchBackend
func Java_com_pojtinger_felicitas_hydrapp_example_rest_MainActivity_LaunchBackend(env *C.JNIEnv, activity C.jobject) C.jstring {
	backendURL, _, err := backend.StartServer(context.Background(), "", false)
	if err != nil {
		log.Fatalln("could not start backend:", err)
	}

	log.Println("Backend URL:", backendURL)

	frontendURL, _, err := frontend.StartServer(context.Background(), "", backendURL, false)
	if err != nil {
		log.Fatalln("could not start frontend:", err)
	}

	log.Println("Frontend URL:", frontendURL)

	return C.get_java_string(env, C.CString(frontendURL))
}

func main() {}
