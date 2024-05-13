package helper

import (
	"errors"
	"strings"
)

type ErrList struct {
	Errors []error
}

const (
	handle_error_prefix string = "handle_error"
)

func (e *ErrList) Handle(result *error) {
	if r := recover(); r != nil {
		switch t := r.(type) {
		case string:
			if strings.HasPrefix(t, handle_error_prefix) {
				*result = e.ErrorsReturn()
				return
			}
		}
		// panic if not error raised with ErrorRaise
		panic(r)
	}
}

func (e *ErrList) ErrorRaise(err error) {
	if err != nil {
		e.ErrorAdd(err)
		panic(handle_error_prefix + err.Error())
	}
}

func (e *ErrList) ErrorAdd(err error) {
	e.Errors = append(e.Errors, err)
}

func (e *ErrList) ErrorsReturn() error {
	// err := fmt.Errorf("%w; %w; %w", err, err2, err3)
	if len(e.Errors) != 0 {
		return errors.Join(e.Errors...)
	}
	return nil
}
