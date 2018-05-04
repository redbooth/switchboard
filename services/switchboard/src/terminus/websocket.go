package terminus

import (
	"../header"
	"github.com/gorilla/websocket"
	"log"
)

type WebsocketConf struct {
	Address string
}

type WebSocket struct {
	conf   WebsocketConf
	errors chan<- error
	conn   *websocket.Conn
}

func NewWebSocket(conf WebsocketConf, errors chan<- error) *WebSocket {
	// connect to upstream
	conn, _, err := websocket.DefaultDialer.Dial(conf.Address, nil)
	if err != nil {
		log.Panicf("Unable to open websocket connection to %s: %v", conf.Address, err)
	}
	return &WebSocket{conf, errors, conn}
}

func (terminus *WebSocket) Terminate(h header.Header) {
	err := terminus.conn.WriteMessage(websocket.TextMessage, []byte(h.String()))
	if err != nil {
		terminus.errors <- err
	}
}
