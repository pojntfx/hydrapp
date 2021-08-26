package backend

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"

	"github.com/phayes/freeport"
)

func StartServer() (string, error) {
	// Get free local socket
	port, err := freeport.GetFreePort()
	if err != nil {
		return "", err
	}
	laddr := net.JoinHostPort("localhost", strconv.Itoa(port))

	// Start a Hello, world server in the background
	http.HandleFunc("/", HelloWorldServer)

	lis, err := net.Listen("tcp", laddr)
	if err != nil {
		return "", err
	}

	go http.Serve(lis, nil)

	// Get the URL for the free local socket
	url := url.URL{
		Scheme: "http",
		Host:   laddr,
	}

	return url.String(), nil
}

func HelloWorldServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, world!")
}
