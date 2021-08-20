package bindings

/*
#include "toast.h"
*/
import "C"
import "golang.org/x/mobile/app"

func ShowToast(msg string) error {
	return app.RunOnJVM(func(vm, env, ctx uintptr) error {
		C.show_toast(C.uintptr_t(vm), C.uintptr_t(env), C.uintptr_t(ctx), C.CString(msg))

		return nil
	})
}
