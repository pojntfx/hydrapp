package backend

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pojntfx/hydrapp/hydrapp-utils/pkg/rpc"
	"github.com/pojntfx/hydrapp/hydrapp-utils/pkg/utils"
)

type exampleStruct struct {
	Name string `json:"name"`
}

func StartServer(addr string, heartbeat time.Duration, localhostize bool) (string, func() error, error) {
	if strings.TrimSpace(addr) == "" {
		addr = ":0"
	}

	registry := rpc.NewRegistry(heartbeat, &rpc.Callbacks{
		OnReceivePong: func() {
			log.Println("Received pong from client")
		},
		OnSendingPing: func() {
			log.Println("Sending ping to client")
		},
		OnFunctionCall: func(requestID, functionName string, functionArgs []json.RawMessage) {
			log.Printf("Got request ID %v for function %v with args %v", requestID, functionName, functionArgs)
		},
	})

	if err := registry.Bind("examplePrintString", func(msg string) {
		fmt.Println(msg)
	}); err != nil {
		panic(err)
	}

	if err := registry.Bind("examplePrintStruct", func(
		input exampleStruct,
	) {
		fmt.Println(input)
	}); err != nil {
		panic(err)
	}

	if err := registry.Bind("exampleReturnError", func() error {
		return errors.New("test error")
	}); err != nil {
		panic(err)
	}

	if err := registry.Bind("exampleReturnString", func() string {
		return "Test string"
	}); err != nil {
		panic(err)
	}

	if err := registry.Bind("exampleReturnStruct", func() exampleStruct {
		return exampleStruct{
			Name: "Alice",
		}
	}); err != nil {
		panic(err)
	}

	if err := registry.Bind("exampleReturnStringAndError", func() (string, error) {
		return "Test string", errors.New("test error")
	}); err != nil {
		panic(err)
	}

	if err := registry.Bind("exampleReturnStringAndNil", func() (string, error) {
		return "Test string", nil
	}); err != nil {
		panic(err)
	}

	var notificationChan chan string
	if err := registry.Bind("exampleNotification", func() (string, error) {
		if notificationChan == nil {
			notificationChan = make(chan string)

			ticker := time.NewTicker(time.Second * 2)
			i := 0
			go func() {
				for {
					<-ticker.C

					if i >= 3 {
						notificationChan <- ""

						ticker.Stop()

						notificationChan = nil

						return
					}

					notificationChan <- "Go server time: " + time.Now().Format(time.RFC3339)

					i++
				}
			}()
		}

		return <-notificationChan, nil
	}); err != nil {
		panic(err)
	}

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
				if err := registry.HandlerFunc(w, r); err != nil {
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
