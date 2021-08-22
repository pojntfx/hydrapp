//go:build android

package main

/*
#include "toast.h"
*/
import "C"
import (
	"unsafe"

	"github.com/pojntfx/hydrapp/example/pkg/bindings"
)

//export Java_com_pojtinger_gomobilefreeexperiments_MainActivity_showToast
func Java_com_pojtinger_gomobilefreeexperiments_MainActivity_showToast(env *C.JNIEnv, _ C.jobject, ctx C.jobject, raw_msg C.jstring) {
	msg := C.CGoGetStringUTFChars(env, raw_msg)

	if err := bindings.ShowToast(uintptr(unsafe.Pointer(env)), uintptr(ctx), C.GoString(msg)); err != nil {
		panic(err)
	}
}

func main() {}
