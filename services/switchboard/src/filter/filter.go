package filter

import (
	"../header"
	"time"
)

type Conf interface{}

type Filter interface {
	Filter(header.Header) bool
}

type Constructor (func(Conf, chan<- error) Filter)

type Bundle struct {
	Conf        Conf
	Constructor Constructor
}

type MultiFilter []Filter

func NewMultiFilter(bundles []Bundle, errors chan<- error) MultiFilter {
	multi := make(MultiFilter, len(bundles))
	for i, bundle := range bundles {
		multi[i] = bundle.Constructor(bundle.Conf, errors)
	}
	return multi
}

func (filters MultiFilter) Filter(h header.Header) bool {
	if len(filters) == 0 {
		return true
	}

	results := make(chan bool, len(filters))
	timeout := time.After(time.Second)

	// execute all the filters in parallel
	for _, filter := range filters {
		go func(f Filter) {
			results <- f.Filter(h)
		}(filter)
	}

	// immediately return false if any fail
	result := true
	select {
	case result = <-results:
		if !result {
			break
		}
	case <-timeout:
		result = false
		break
	}
	return result
}
