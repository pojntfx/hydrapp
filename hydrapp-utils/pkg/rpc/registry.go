package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/gorilla/websocket"
)

var (
	ErrCannotExposeNonFunction = errors.New("can not expose non function")
	ErrInvalidReturn           = errors.New("can only return void, a value or a value and an error")

	errorType = reflect.TypeOf((*error)(nil)).Elem()

	upgrader = websocket.Upgrader{}
)

type Callbacks struct {
	OnReceivePong  func()
	OnSendingPing  func()
	OnFunctionCall func(requestID, functionName string, functionArgs []json.RawMessage)
}

type Registry struct {
	functions map[string]interface{}
	callbacks *Callbacks
	heartbeat time.Duration
}

func NewRegistry(heartbeat time.Duration, callbacks *Callbacks) *Registry {
	if callbacks == nil {
		callbacks = &Callbacks{}
	}

	return &Registry{map[string]interface{}{}, callbacks, heartbeat}
}

func (h *Registry) Bind(functionName string, function interface{}) error {
	v := reflect.ValueOf(function)

	if v.Kind() != reflect.Func {
		return ErrCannotExposeNonFunction
	}

	if n := v.Type().NumOut(); n > 2 || (n == 2 && !v.Type().Out(1).Implements(errorType)) {
		return ErrInvalidReturn
	}

	h.functions[functionName] = function

	return nil
}

func (h *Registry) HandlerFunc(w http.ResponseWriter, r *http.Request) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	if err := conn.SetReadDeadline(time.Now().Add(h.heartbeat)); err != nil {
		return err
	}

	conn.SetPongHandler(func(string) error {
		if h.callbacks.OnReceivePong != nil {
			h.callbacks.OnReceivePong()
		}

		return conn.SetReadDeadline(time.Now().Add(h.heartbeat))
	})

	pings := time.NewTicker(h.heartbeat / 2)
	defer pings.Stop()

	errs := make(chan error)
	go func() {
		defer conn.Close()

		functionNames := []string{}
		for functionName := range h.functions {
			functionNames = append(functionNames, functionName)
		}

		if err := conn.WriteJSON(functionNames); err != nil {
			errs <- err

			return
		}

		for {
			var functionRequest []json.RawMessage
			if err := conn.ReadJSON(&functionRequest); err != nil {
				errs <- err

				return
			}

			if len(functionRequest) != 3 {
				errs <- fmt.Errorf("%v", http.StatusUnprocessableEntity)

				return
			}

			var requestID string
			if err := json.Unmarshal(functionRequest[0], &requestID); err != nil {
				errs <- fmt.Errorf("%v", http.StatusUnprocessableEntity)

				return
			}

			var functionName string
			if err := json.Unmarshal(functionRequest[1], &functionName); err != nil {
				errs <- fmt.Errorf("%v", http.StatusUnprocessableEntity)

				return
			}

			var functionArgs []json.RawMessage
			if err := json.Unmarshal(functionRequest[2], &functionArgs); err != nil {
				errs <- fmt.Errorf("%v", http.StatusUnprocessableEntity)

				return
			}

			if h.callbacks.OnFunctionCall != nil {
				h.callbacks.OnFunctionCall(requestID, functionName, functionArgs)
			}

			rawFunctions, ok := h.functions[functionName]
			if !ok {
				errs <- fmt.Errorf("%v", http.StatusNotFound)

				return
			}

			function := reflect.ValueOf(rawFunctions)

			if len(functionArgs) != function.Type().NumIn() {
				errs <- fmt.Errorf("%v", http.StatusUnprocessableEntity)

				return
			}

			args := []reflect.Value{}
			for i := range functionArgs {
				arg := reflect.New(function.Type().In(i))
				if err := json.Unmarshal(functionArgs[i], arg.Interface()); err != nil {
					errs <- err

					return
				}

				args = append(args, arg.Elem())
			}

			go func() {
				res := function.Call(args)
				switch len(res) {
				case 0:
					if err := conn.WriteJSON([]interface{}{requestID, nil, ""}); err != nil {
						errs <- err

						return
					}
				case 1:
					if res[0].Type().Implements(errorType) {
						if err := conn.WriteJSON([]interface{}{requestID, nil, res[0].Interface().(error).Error()}); err != nil {
							errs <- err

							return
						}
					} else {
						v, err := json.Marshal(res[0].Interface())
						if err != nil {
							errs <- err

							return
						}

						if err := conn.WriteJSON([]interface{}{requestID, json.RawMessage(string(v)), ""}); err != nil {
							errs <- err

							return
						}
					}
				case 2:
					v, err := json.Marshal(res[0].Interface())
					if err != nil {
						errs <- err

						return
					}

					if res[1].Interface() == nil {
						if err := conn.WriteJSON([]interface{}{requestID, json.RawMessage(string(v)), ""}); err != nil {
							errs <- err

							return
						}
					} else {
						if err := conn.WriteJSON([]interface{}{requestID, json.RawMessage(string(v)), res[1].Interface().(error).Error()}); err != nil {
							errs <- err

							return
						}
					}
				}
			}()
		}
	}()

	for {
		select {
		case err := <-errs:
			return err
		case <-pings.C:
			if h.callbacks.OnSendingPing != nil {
				h.callbacks.OnSendingPing()
			}

			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return err
			}

			if err := conn.SetWriteDeadline(time.Now().Add(h.heartbeat)); err != nil {
				return err
			}
		}
	}
}

func init() {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
}
