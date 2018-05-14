package input

import (
	"io"
	"sync"
)

type Conf interface{}

type Input interface {
	Read()
}

type Constructor (func(Conf, chan<- error, chan<- io.ReadCloser) Input)

type Bundle struct {
	Conf        Conf
	Constructor Constructor
}

type MultiInput []Input

func NewMultiInput(bundles []Bundle, errors chan<- error, readers chan<- io.ReadCloser) MultiInput {
	multi := make(MultiInput, len(bundles))
	for i, bundle := range bundles {
		multi[i] = bundle.Constructor(bundle.Conf, errors, readers)
	}
	return multi
}

func (inputs MultiInput) Read() {
	wg := &sync.WaitGroup{}
	for _, input := range inputs {
		wg.Add(1)
		go func(input Input) {
			defer wg.Done()
			input.Read()
		}(input)
	}
	wg.Wait()
}
