//go:build tlscertembed

package fixes

import (
	"crypto/tls"
	"net/http"
)

func init() {
	// Disable TLS certificate validation
	// TODO: Embed certificates instead
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}
