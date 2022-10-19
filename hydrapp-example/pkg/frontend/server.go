package frontend

//go:generate npm i
//go:generate npm run build

import (
	"embed"
	"io/fs"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/pojntfx/hydrapp/hydrapp-utils/pkg/utils"
)

var (
	//go:embed out
	UI embed.FS
)

func StartServer(addr string, backendURL string) (string, func() error, error) {
	var listener net.Listener
	if strings.TrimSpace(addr) != "" {
		var err error
		listener, err = net.Listen("tcp", addr)
		if err != nil {
			return "", nil, err
		}
	} else {
		var err error
		listener, err = net.Listen("tcp", ":0")
		if err != nil {
			return "", nil, err
		}
	}

	root := fs.FS(UI)
	dist, err := fs.Sub(root, "out")
	if err != nil {
		panic(err)
	}

	go func() {
		if err := http.Serve(listener, http.FileServer(http.FS(dist))); err != nil {
			panic(err)
		}
	}()

	laddr := listener.Addr().String()

	laddr = strings.Replace(laddr, "127.0.0.1", "localhost", 1)
	laddr = strings.Replace(laddr, "[::]", "localhost", 1)

	url, err := url.Parse("http://" + laddr)
	if err != nil {
		return "", nil, err
	}

	values := url.Query()

	values.Set("socketURL", backendURL)

	url.RawQuery = values.Encode()

	return utils.Localhostize(url.String()), listener.Close, nil
}
