package backend

import (
	"io"
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

	// Start an example server in the background
	http.HandleFunc("/", ExampleServer)

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

func ExampleServer(w http.ResponseWriter, r *http.Request) {
	res, err := http.Get("https://jsonplaceholder.typicode.com/users/1")
	if err != nil {
		if _, err := w.Write([]byte(err.Error())); err != nil {
			panic(err)
		}

		return
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		if _, err := w.Write([]byte(err.Error())); err != nil {
			panic(err)
		}

		return
	}

	if _, err := w.Write(data); err != nil {
		panic(err)
	}
}
