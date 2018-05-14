package handler

import "github.com/redbooth/switchboard/src/errors"

type Conf interface{}

type Handler interface {
	Handle(error)
	Close() error
}

type Constructor (func(Conf) Handler)

type Bundle struct {
	Conf        Conf
	Constructor Constructor
}

type MultiHandler []Handler

func NewMultiHandler(bundles []Bundle) MultiHandler {
	multi := make(MultiHandler, len(bundles))
	for i, bundle := range bundles {
		multi[i] = bundle.Constructor(bundle.Conf)
	}
	return multi
}

// execute all the handlers in parallel
func (handlers MultiHandler) Handle(err error) {
	for _, handler := range handlers {
		go handler.Handle(err)
	}
}

func (handlers MultiHandler) Close() error {
	merr := errors.NewMultiError()
	for _, handler := range handlers {
		merr.Append(handler.Close())
	}
	return merr.Error()
}
