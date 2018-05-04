package terminus

import (
	"../header"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

type HTTPConf struct {
	Address string
}

type HTTP struct {
	conf   HTTPConf
	errors chan<- error
}

func NewHTTP(conf HTTPConf, errors chan<- error) *HTTP {
	return &HTTP{conf, errors}
}

func (terminus *HTTP) Terminate(h header.Header) {
	resp, err := http.Post(terminus.conf.Address, "text/plain", strings.NewReader(h.String()))
	if err != nil {
		terminus.errors <- err
	} else if resp.StatusCode == 200 {
		return
	} else if body, err := ioutil.ReadAll(resp.Body); err != nil {
		terminus.errors <- err
	} else {
		terminus.errors <- errors.New(string(body))
	}
}
