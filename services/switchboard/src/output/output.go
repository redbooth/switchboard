package output

import (
	"github.com/redbooth/switchboard/src/errors"
	"github.com/redbooth/switchboard/src/header"
	"io"
)

type Conf interface{}

type Output io.WriteCloser

type Constructor (func(Conf, chan<- error, header.Header) Output)

type Bundle struct {
	Conf        Conf
	Constructor Constructor
}

type MultiOutput struct {
	outputs []Output
	writer  io.Writer
}

func NewMultiOutput(bundles []Bundle, errors chan<- error, h header.Header) MultiOutput {
	outputs := make([]Output, len(bundles))
	writers := make([]io.Writer, len(bundles))
	for i, bundle := range bundles {
		outputs[i] = bundle.Constructor(bundle.Conf, errors, h)
		writers[i] = outputs[i].(io.Writer)
	}
	return MultiOutput{outputs, io.MultiWriter(writers...)}
}

func (outputs MultiOutput) Write(b []byte) (n int, err error) {
	return outputs.writer.Write(b)
}

func (outputs MultiOutput) Close() error {
	merr := errors.NewMultiError()
	for _, output := range outputs.outputs {
		merr.Append(output.Close())
	}
	return merr.Error()
}
