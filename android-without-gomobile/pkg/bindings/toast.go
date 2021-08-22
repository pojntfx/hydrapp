//go:build android

package bindings

/*
#include "toast.h"
*/
import "C"

func ShowToast(env, ctx uintptr, msg string) error {
	C.show_toast(C.uintptr_t(env), C.uintptr_t(ctx), C.CString(msg))

	return nil
}
