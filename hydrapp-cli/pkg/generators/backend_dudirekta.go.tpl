package backend

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pojntfx/dudirekta/pkg/rpc"
	"github.com/pojntfx/hydrapp/hydrapp-utils/pkg/utils"
	"nhooyr.io/websocket"
)

type exampleStruct struct {
	Name string `json:"name"`
}

type local struct {
	Peers func() map[string]remote
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

func (l *local) ExampleNotification(ctx context.Context) error {
	for peerID, peer := range l.Peers() {
		if peerID == rpc.GetRemoteID(ctx) {
			ticker := time.NewTicker(time.Second)
			i := 0
			for {
				if i >= 3 {
					ticker.Stop()

					return nil
				}

				<-ticker.C

				if err := peer.ExampleNotification(ctx, time.Now().Format(time.RFC3339)); err != nil {
					return err
				}

				i++
			}
		}
	}

	return nil
}

type remote struct {
	ExampleNotification func(ctx context.Context, msg string) error
}

func StartServer(ctx context.Context, addr string, heartbeat time.Duration, localhostize bool) (string, func() error, error) {
	if strings.TrimSpace(addr) == "" {
		addr = ":0"
	}

	l := &local{}
	registry := rpc.NewRegistry(
		l,
		remote{},

		time.Second*10,
		ctx,
		nil,
	)
	l.Peers = registry.Peers

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return "", nil, err
	}

	clients := 0
	go func() {
		if err := http.Serve(listener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clients++

			log.Printf("%v clients connected", clients)

			defer func() {
				clients--

				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)

					log.Printf("Client disconnected with error: %v", err)
				}

				log.Printf("%v clients connected", clients)
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

				go func() {
					if err := registry.Link(conn); err != nil {
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
