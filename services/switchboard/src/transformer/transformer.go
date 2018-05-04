package transformer

import (
	"io"
)

type Conf interface{}

type Transformer interface {
	Transform(io.Reader) io.Reader
}

type Constructor (func(Conf, chan<- error) Transformer)

type Bundle struct {
	Conf        Conf
	Constructor Constructor
}

type MultiTransformer []Transformer

func NewMultiTransformer(bundles []Bundle, errors chan<- error) MultiTransformer {
	multi := make(MultiTransformer, len(bundles))
	for i, bundle := range bundles {
		multi[i] = bundle.Constructor(bundle.Conf, errors)
	}
	return multi
}

// execute transformers in series
func (transformers MultiTransformer) Transform(reader io.Reader) io.Reader {
	for _, transformer := range transformers {
		reader = transformer.Transform(reader)
	}
	return reader
}
