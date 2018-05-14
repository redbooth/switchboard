package output

import (
	"github.com/redbooth/switchboard/src/header"
	"github.com/gorilla/websocket"
	"log"
)

type WebsocketConf struct {
	Address string
}

type WebSocket struct {
	conf   WebsocketConf
	errors chan<- error
	header header.Header
	conn   *websocket.Conn
}

func NewWebSocket(conf WebsocketConf, errors chan<- error, h header.Header) *WebSocket {
	// connect to upstream
	conn, _, err := websocket.DefaultDialer.Dial(conf.Address, nil)
	if err != nil {
		log.Panicf("Unable to open websocket connection to %s: %v", conf.Address, err)
	}
	// inject header
	bytes, err := header.Bytes(h)
	if err != nil {
		log.Panicf("Unable to convert header into bytes: %v\n", err)
	}
	err = conn.WriteMessage(websocket.BinaryMessage, bytes)
	if err != nil {
		log.Panicf("Unable to send message for stream %s: %v", h.Id(), err)
	}
	return &WebSocket{conf, errors, h, conn}
}

func (writer *WebSocket) Write(b []byte) (n int, err error) {
	err = writer.conn.WriteMessage(websocket.BinaryMessage, b)
	if err != nil {
		return 0, err
	}
	return len(b), nil
}

func (writer *WebSocket) Close() error {
	return writer.conn.Close()
}
