package terminus

import (
	"github.com/redbooth/switchboard/src/header"
)

type Conf interface{}

type Terminus interface {
	Terminate(header.Header)
}

type Constructor (func(Conf, chan<- error) Terminus)

type Bundle struct {
	Conf        Conf
	Constructor Constructor
}

type MultiTerminus []Terminus

func NewMultiTerminus(bundles []Bundle, errors chan<- error) MultiTerminus {
	multi := make(MultiTerminus, len(bundles))
	for i, bundle := range bundles {
		multi[i] = bundle.Constructor(bundle.Conf, errors)
	}
	return multi
}

func (termini MultiTerminus) Terminate(h header.Header) {
	// execute all the termini in parallel
	for _, terminus := range termini {
		go terminus.Terminate(h)
	}
}
