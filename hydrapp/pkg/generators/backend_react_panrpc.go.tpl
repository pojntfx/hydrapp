package backend

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pojntfx/hydrapp/hydrapp/pkg/utils"
	"github.com/pojntfx/panrpc/go/pkg/rpc"
	"nhooyr.io/websocket"
)

type exampleStruct struct {
	Name string `json:"name"`
}

type local struct {
	ForRemotes func(
		cb func(remoteID string, remote remote) error,
	) error
}

func (l *local) ExamplePrintString(ctx context.Context, msg string) error {
	fmt.Println(msg)

	return nil
}

func (l *local) ExamplePrintStruct(ctx context.Context, input exampleStruct) error {
	fmt.Println(input)

	return nil
}

func (l *local) ExampleReturnError(ctx context.Context) error {
	return errors.New("test error")
}

func (l *local) ExampleReturnString(ctx context.Context) (string, error) {
	return "Test string", nil
}

func (l *local) ExampleReturnStruct(ctx context.Context) (exampleStruct, error) {
	return exampleStruct{
		Name: "Alice",
	}, nil
}

func (l *local) ExampleReturnStringAndError(ctx context.Context) (string, error) {
	return "Test string", errors.New("test error")
}

func (l *local) ExampleCallback(ctx context.Context) error {
	var peer *remote

	_ = l.ForRemotes(func(remoteID string, remote remote) error {
		peer = &remote

		return nil
	})

	if peer != nil {
		ticker := time.NewTicker(time.Second)
		i := 0
		for {
			if i >= 3 {
				ticker.Stop()

				return nil
			}

			<-ticker.C

			if err := peer.ExampleNotification(ctx, "Backend time: "+time.Now().Format(time.RFC3339)); err != nil {
				return err
			}

			i++
		}
	}

	return nil
}

func (s *local) ExampleClosure(
	ctx context.Context,
	length int,
	onIteration func(ctx context.Context, i int, b string) (string, error),
) (int, error) {
	for i := 0; i < length; i++ {
		rv, err := onIteration(ctx, i, "This is from the backend")
		if err != nil {
			return -1, err
		}

		log.Println("Closure returned:", rv)
	}

	return length, nil
}

type remote struct {
	ExampleNotification func(ctx context.Context, msg string) error
}

func StartServer(ctx context.Context, addr string, heartbeat time.Duration, localhostize bool) (string, func() error, error) {
	if strings.TrimSpace(addr) == "" {
		addr = ":0"
	}

	service := &local{}

	clients := 0
	registry := rpc.NewRegistry[remote, json.RawMessage](
		service,

		ctx,

		&rpc.RegistryHooks{
			OnClientConnect: func(remoteID string) {
				clients++

				log.Printf("%v clients connected", clients)
			},
			OnClientDisconnect: func(remoteID string) {
				clients--

				log.Printf("%v clients connected", clients)
			},
		},
	)
	service.ForRemotes = registry.ForRemotes

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	log.Println("Listening on", listener.Addr())

	go func() {
		defer listener.Close()

		if err := http.Serve(listener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)

					log.Printf("Client disconnected with error: %v", err)
				}
			}()

			switch r.Method {
			case http.MethodGet:
				c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
					OriginPatterns: []string{"*"},
				})
				if err != nil {
					panic(err)
				}

				pings := time.NewTicker(time.Second / 2)
				defer pings.Stop()

				errs := make(chan error)
				go func() {
					for range pings.C {
						if err := c.Ping(ctx); err != nil {
							errs <- err

							return
						}
					}
				}()

				conn := websocket.NetConn(ctx, c, websocket.MessageText)
				defer conn.Close()

				encoder := json.NewEncoder(conn)
				decoder := json.NewDecoder(conn)

				go func() {
					if err := registry.LinkStream(
						func(v rpc.Message[json.RawMessage]) error {
							return encoder.Encode(v)
						},
						func(v *rpc.Message[json.RawMessage]) error {
							return decoder.Decode(v)
						},

						func(v any) (json.RawMessage, error) {
							b, err := json.Marshal(v)
							if err != nil {
								return nil, err
							}

							return json.RawMessage(b), nil
						},
						func(data json.RawMessage, v any) error {
							return json.Unmarshal([]byte(data), v)
						},

						nil,
					); err != nil {
						errs <- err

						return
					}
				}()

				if err := <-errs; err != nil {
					panic(err)
				}
			default:
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		})); err != nil {
			if strings.HasSuffix(err.Error(), "use of closed network connection") {
				return
			}

			panic(err)
		}
	}()

	url, err := url.Parse("ws://" + listener.Addr().String())
	if err != nil {
		return "", nil, err
	}

	if localhostize {
		return utils.Localhostize(url.String()), listener.Close, nil
	}

	return url.String(), listener.Close, nil
}

