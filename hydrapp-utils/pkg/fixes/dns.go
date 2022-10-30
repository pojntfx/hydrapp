//go:build androiddnsfix && android
// +build androiddnsfix,android

package fixes

import (
	_ "unsafe"
)

//go:linkname defaultNS net.defaultNS
var defaultNS []string

func SetDefaultNS(ns []string) {
	defaultNS = ns
}

// See https://gist.github.com/cs8425/107e01a0652f1f1f6e033b5b68364b5e
func init() {
	// Use public DNS servers; useful on i.e. old Android versions which don't have `/etc/resolv.conf`, where DNS lookups would be broken
	SetDefaultNS([]string{"8.8.8.8:53", "8.8.4.4:53", "[2001:4860:4860::8888]:53", "[2001:4860:4860::8844]:53", "1.1.1.1:53", "1.0.0.1:53", "[2606:4700:4700::1111]:53", "[2606:4700:4700::1001]:53"})
}
