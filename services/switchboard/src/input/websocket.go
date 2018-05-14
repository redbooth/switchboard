package input

import (
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
)

type WebsocketConf struct {
	Port uint16
}

type WebSocket struct {
	conf     WebsocketConf
	errors   chan<- error
	readers  chan<- io.ReadCloser
	upgrader *websocket.Upgrader
}

func NewWebSocket(conf WebsocketConf, errors chan<- error, readers chan<- io.ReadCloser) *WebSocket {
	upgrader := &websocket.Upgrader{
		CheckOrigin:       func(r *http.Request) bool { return true },
		EnableCompression: true,
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
	}
	return &WebSocket{conf, errors, readers, upgrader}
}

func (input *WebSocket) Read() {
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		conn, err := input.upgrader.Upgrade(res, req, nil)
		if err != nil {
			input.errors <- err
			return
		}
		defer conn.Close()

		pr, pw := io.Pipe()
		defer pw.Close()
		input.readers <- pr

		for {
			_, reader, err := conn.NextReader()
			if err != nil {
				input.errors <- err
				break
			} else {
				io.Copy(pw, reader)
			}
		}
	})
	port := fmt.Sprintf(":%d", input.conf.Port)
	input.errors <- http.ListenAndServe(port, nil)
}
