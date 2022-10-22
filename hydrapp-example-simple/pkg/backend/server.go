package backend

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pojntfx/hydrapp/hydrapp-utils/pkg/utils"
)

func StartServer(addr string, localhostize bool) (string, func() error, error) {
	if strings.TrimSpace(addr) == "" {
		addr = ":0"
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return "", nil, err
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/servertime", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Go server time: " + time.Now().Format(time.RFC3339)))
	})

	mux.HandleFunc("/ifconfigio", func(w http.ResponseWriter, r *http.Request) {
		res, err := http.Get("https://ifconfig.io/all.json")
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
	})

	go func() {
		if err := http.Serve(listener, mux); err != nil {
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

	if localhostize {
		return utils.Localhostize(url.String()), listener.Close, nil
	}

	return url.String(), listener.Close, nil
}
