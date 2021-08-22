//go:build android

package main

/*
#cgo LDFLAGS: -llog

#include <jni.h>
*/
import "C"
import (
	"unsafe"

	"github.com/pojntfx/hydrapp/example/pkg/bindings"
)

//export Java_com_pojtinger_gomobilefreeexperiments_MainActivity_showToast
func Java_com_pojtinger_gomobilefreeexperiments_MainActivity_showToast(env *C.JNIEnv, _ C.jobject, ctx C.jobject, msg C.jstring) {
	if err := bindings.ShowToast(uintptr(unsafe.Pointer(env)), uintptr(ctx), "Hello, world!"); err != nil {
		panic(err)
	}
}

func main() {}
