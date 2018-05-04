package errors

import (
	"fmt"
	"strings"
)

type MultiError []error

func NewMultiError() MultiError {
	return make(MultiError, 0)
}

func (merr MultiError) Append(err error) MultiError {
	if err != nil {
		merr = append(merr, err)
	}
	return merr
}

func (merr MultiError) Error() error {
	if len(merr) == 0 {
		return nil
	} else {
		var text []string
		for _, err := range merr {
			text = append(text, err.Error())
		}
		return fmt.Errorf("%d errors encountered:\n%s\n", len(merr), strings.Join(text, "\n"))
	}
}
