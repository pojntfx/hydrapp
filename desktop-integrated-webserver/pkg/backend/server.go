package backend

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"

	"github.com/phayes/freeport"
)

func StartServer() (string, func() error, error) {
	// Get a free port
	port, err := freeport.GetFreePort()
	if err != nil {
		return "", nil, err
	}

	// Bind to localhost
	laddr := net.JoinHostPort("localhost", strconv.Itoa(port))
	lis, err := net.Listen("tcp", laddr)
	if err != nil {
		return "", nil, err
	}

	// Start the server
	srv := &http.Server{}
	http.HandleFunc("/", ExampleHandler)
	go srv.Serve(lis)

	// Construct the URL on which the server is being served
	url := url.URL{
		Scheme: "http",
		Host:   laddr,
	}

	return url.String(), srv.Close, nil
}

func ExampleHandler(w http.ResponseWriter, r *http.Request) {
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
