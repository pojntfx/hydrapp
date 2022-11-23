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

	_ "github.com/pojntfx/hydrapp/hydrapp-utils/pkg/fixes"

	frontend "github.com/pojntfx/hydrapp/hydrapp-example-forms/pkg/frontend"
)

//export Java_com_pojtinger_felicitas_hydrapp_example_forms_MainActivity_LaunchBackend
func Java_com_pojtinger_felicitas_hydrapp_example_forms_MainActivity_LaunchBackend(env *C.JNIEnv, activity C.jobject) C.jstring {
	frontendURL, _, err := frontend.StartServer(context.Background(), "", false)
	if err != nil {
		log.Fatalln("could not start frontend:", err)
	}

	log.Println("Frontend URL:", frontendURL)

	return C.get_java_string(env, C.CString(frontendURL))
}

func main() {}
