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
	"time"

	_ "github.com/pojntfx/hydrapp/hydrapp/pkg/fixes"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/utils"

	backend "github.com/pojntfx/hydrapp/hydrapp-example-react-panrpc/pkg/backend"
	frontend "github.com/pojntfx/hydrapp/hydrapp-example-react-panrpc/pkg/frontend"
)

//export Java_com_pojtinger_felicitas_hydrapp_example_react_panrpc_MainActivity_LaunchBackend
func Java_com_pojtinger_felicitas_hydrapp_example_react_panrpc_MainActivity_LaunchBackend(env *C.JNIEnv, activity C.jobject, filesDir C.jstring) C.jstring {
	if err := utils.PolyfillEnvironment(C.GoString(C.get_c_string(env, filesDir))); err != nil {
		log.Fatalln("could not polyfill environment:", err)
	}

	backendURL, _, err := backend.StartServer(context.Background(), "", time.Second*10, true)
	if err != nil {
		log.Fatalln("could not start backend:", err)
	}

	log.Println("Backend URL:", backendURL)

	frontendURL, _, err := frontend.StartServer(context.Background(), "", backendURL, true)
	if err != nil {
		log.Fatalln("could not start frontend:", err)
	}

	log.Println("Frontend URL:", frontendURL)

	return C.get_java_string(env, C.CString(frontendURL))
}

func main() {}
