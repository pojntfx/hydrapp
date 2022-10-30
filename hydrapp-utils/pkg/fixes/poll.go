//go:build androidacceptfix && arm && android
// +build androidacceptfix,arm,android

package fixes

import (
	"syscall"
	"unsafe"
	_ "unsafe"
)

//go:linkname Accept4Func internal/poll.Accept4Func
var Accept4Func func(int, int) (int, syscall.Sockaddr, error)

//go:linkname AcceptFunc internal/poll.AcceptFunc
var AcceptFunc func(int) (int, syscall.Sockaddr, error)

type _Socklen uint32

//go:linkname errnoErr syscall.errnoErr
func errnoErr(e syscall.Errno) error

//go:linkname accept syscall.accept
func accept(s int, rsa *syscall.RawSockaddrAny, addrlen *_Socklen) (fd int, err error) {
	r0, _, e1 := syscall.Syscall(syscall.SYS_ACCEPT, uintptr(s), uintptr(unsafe.Pointer(rsa)), uintptr(unsafe.Pointer(addrlen)))
	fd = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}

//go:linkname accept4 syscall.accept4
func accept4(s int, rsa *syscall.RawSockaddrAny, addrlen *_Socklen, flags int) (fd int, err error)

//go:linkname anyToSockaddr syscall.anyToSockaddr
func anyToSockaddr(rsa *syscall.RawSockaddrAny) (syscall.Sockaddr, error)

// This restores the old behavior from Go 1.16, which was changed in Go 1.17.
// Enables the use of `accept` on kernels < 2.6.28
// See delta between https://go.googlesource.com/sys/+/refs/changes/90/313690/3/unix/syscall_linux.go#1155 and https://go.googlesource.com/sys/+/refs/heads/master/unix/syscall_linux.go#1239 for when they removed the fallback
func init() {
	Accept4Func = func(fd int, flags int) (nfd int, sa syscall.Sockaddr, err error) {
		var rsa syscall.RawSockaddrAny
		var len _Socklen = syscall.SizeofSockaddrAny
		nfd, err = accept4(fd, &rsa, &len, flags)
		if err != nil {
			return
		}
		if len > syscall.SizeofSockaddrAny {
			panic("RawSockaddrAny too small")
		}
		sa, err = anyToSockaddr(&rsa)
		if err != nil {
			syscall.Close(nfd)
			nfd = 0
		}
		return
	}

	AcceptFunc = func(fd int) (nfd int, sa syscall.Sockaddr, err error) {
		var rsa syscall.RawSockaddrAny
		var len _Socklen = syscall.SizeofSockaddrAny
		// Try accept4 first for Android, then try accept for kernel older than 2.6.28
		nfd, err = accept4(fd, &rsa, &len, 0)
		if err == syscall.ENOSYS {
			nfd, err = accept(fd, &rsa, &len)
		}
		if err != nil {
			return
		}
		sa, err = anyToSockaddr(&rsa)
		if err != nil {
			syscall.Close(nfd)
			nfd = 0
		}
		return
	}
}
