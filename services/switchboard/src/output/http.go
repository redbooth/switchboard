package output

import (
	"../errors"
	"../header"
	"io"
	"log"
	"net/http"
)

type HTTPConf struct {
	Address string
}

type HTTP struct {
	conf   HTTPConf
	errors chan<- error
	header header.Header
	pr     *io.PipeReader
	pw     *io.PipeWriter
}

func NewHTTP(conf HTTPConf, errors chan<- error, h header.Header) *HTTP {
	pr, pw := io.Pipe()
	go func() {
		// connect to upstream
		client := &http.Client{}
		client.Post(conf.Address, "application/octet-stream", pr)
		// inject header
		bytes, err := header.Bytes(h)
		if err != nil {
			log.Panicf("Unable to convert header into bytes: %v\n", err)
		}
		pw.Write(bytes)
	}()
	return &HTTP{conf, errors, h, pr, pw}
}

func (writer *HTTP) Write(b []byte) (n int, err error) {
	return writer.pw.Write(b)
}

func (writer *HTTP) Close() error {
	merr := errors.NewMultiError()
	merr.Append(writer.pr.Close())
	merr.Append(writer.pw.Close())
	return merr.Error()
}
