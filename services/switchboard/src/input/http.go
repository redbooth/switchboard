package input

import (
	"fmt"
	"io"
	"net/http"
)

type HTTPConf struct {
	Port uint16
}

type HTTP struct {
	conf    HTTPConf
	errors  chan<- error
	readers chan<- io.ReadCloser
}

func NewHTTP(conf HTTPConf, errors chan<- error, readers chan<- io.ReadCloser) *HTTP {
	return &HTTP{conf, errors, readers}
}

func (input *HTTP) Read() {
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case "POST":
			input.readers <- req.Body
		}
	})
	port := fmt.Sprintf(":%d", input.conf.Port)
	input.errors <- http.ListenAndServe(port, nil)
}
