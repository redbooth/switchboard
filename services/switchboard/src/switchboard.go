package main

import (
	"./conf"
	"./filter"
	"./handler"
	"./header"
	"./input"
	"./output"
	"./terminus"
	"./transformer"
	"io"
)

func main() {
	config := conf.NewConf()
	errors := make(chan error, 1)
	readers := make(chan io.ReadCloser, 1)

	// define stream processors
	inputs := input.NewMultiInput(config.Inputs, errors, readers)
	filters := filter.NewMultiFilter(config.Filters, errors)
	transformers := transformer.NewMultiTransformer(config.Transformers, errors)
	termini := terminus.NewMultiTerminus(config.Termini, errors)
	handlers := handler.NewMultiHandler(config.Handlers)

	// inputs -> filters -> mappers -> outputs -> termini
	go func() {
		for reader := range readers {
			go func(r io.ReadCloser) {
				defer r.Close()

				h := config.NewHeader()
				header.Extract(h, r)

				if !filters.Filter(h) {
					return
				}

				o := output.NewMultiOutput(config.Outputs, errors, h)
				defer o.Close()

				t := transformers.Transform(r)
				if _, err := io.Copy(o, t); err != nil {
					errors <- err
				} else {
					termini.Terminate(h)
				}
			}(reader)
		}
	}()

	// handle errors
	go func() {
		defer handlers.Close()
		for err := range errors {
			handlers.Handle(err)
		}
	}()

	// start input servers once the pipeline is ready
	// (input.Read blocks until all inputs are exhausted)
	inputs.Read()
}
