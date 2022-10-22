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

func StartServer(network, addr string, backendURL string) (string, func() error, error) {
	if strings.TrimSpace(network) == "" {
		network = "tcp"
	}

	if strings.TrimSpace(addr) == "" {
		addr = ":0"
	}

	listener, err := net.Listen(network, addr)
	if err != nil {
		return "", nil, err
	}

	root := fs.FS(UI)
	dist, err := fs.Sub(root, "out")
	if err != nil {
		panic(err)
	}

	go func() {
		if err := http.Serve(listener, http.FileServer(http.FS(dist))); err != nil {
			if strings.HasSuffix(err.Error(), "use of closed network connection") {
				return
			}

			panic(err)
		}
	}()

	url, err := url.Parse("http://" + listener.Addr().String())
	if err != nil {
		return "", nil, err
	}

	values := url.Query()

	values.Set("socketURL", backendURL)

	url.RawQuery = values.Encode()

	return utils.Localhostize(url.String()), listener.Close, nil
}
