package backend

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"

	"github.com/phayes/freeport"
)

func StartServer() (string, chan error, func() error, error) {
	// Create a server which can be closed
	srv := &http.Server{}

	// Get free local socket
	port, err := freeport.GetFreePort()
	if err != nil {
		return "", nil, nil, err
	}
	laddr := net.JoinHostPort("localhost", strconv.Itoa(port))

	// Start an example server in the background
	http.HandleFunc("/", ExampleServer)

	lis, err := net.Listen("tcp", laddr)
	if err != nil {
		return "", nil, nil, err
	}

	done := make(chan error)
	go func() {
		done <- srv.Serve(lis)
	}()

	// Get the URL for the free local socket
	url := url.URL{
		Scheme: "http",
		Host:   laddr,
	}

	return url.String(), done, srv.Close, nil
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

	data, err := ioutil.ReadAll(res.Body)
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
