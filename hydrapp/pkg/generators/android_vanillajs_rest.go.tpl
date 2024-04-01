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
	"github.com/pojntfx/hydrapp/hydrapp/pkg/utils"

	backend "{{ .GoMod }}/pkg/backend"
	frontend "{{ .GoMod }}/pkg/frontend"
)

//export Java_{{ .JNIExport }}_MainActivity_LaunchBackend
func Java_{{ .JNIExport }}_MainActivity_LaunchBackend(env *C.JNIEnv, activity C.jobject) C.jstring {
	if err := utils.PolyfillEnvironment(C.GoString(C.get_c_string(env, filesDir))); err != nil {
		log.Fatalln("could not polyfill environment:", err)
	}
	
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
