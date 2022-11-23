package frontend

import (
	"context"
	_ "embed"
	"html/template"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pojntfx/hydrapp/hydrapp-utils/pkg/utils"
)

var (
	//go:embed index.html
	indexHTML string
)

type todo struct {
	Title string
	Body  string
}

type data struct {
	Todos map[string]todo

	GoVersion,
	GoOS,
	GoArch,
	RenderTime string
}

func StartServer(ctx context.Context, addr string, localhostize bool) (string, func() error, error) {
	if strings.TrimSpace(addr) == "" {
		addr = ":0"
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return "", nil, err
	}

	index, err := template.New("index.html").Parse(indexHTML)
	if err != nil {
		return "", nil, err
	}

	todos := map[string]todo{}
	var todosLock sync.Mutex

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := index.Execute(
			w,
			data{
				Todos: todos,

				GoVersion:  runtime.Version(),
				GoOS:       runtime.GOOS,
				GoArch:     runtime.GOARCH,
				RenderTime: time.Now().Format(time.RFC3339),
			},
		); err != nil {
			panic(err)
		}
	})

	mux.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			panic(err)
		}

		todosLock.Lock()
		defer todosLock.Unlock()

		todos[uuid.NewString()] = todo{r.FormValue("title"), r.FormValue("body")}

		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
	})

	mux.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		todosLock.Lock()
		defer todosLock.Unlock()

		delete(todos, r.URL.Query().Get("id"))

		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
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
