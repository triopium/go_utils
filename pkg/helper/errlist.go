package helper

import (
	"errors"
)

type ErrList struct {
	Errors []error
}

func (e *ErrList) Handle(result *error) {
	if r := recover(); r != nil {
		*result = e.ErrorsReturn()
	}
}

func (e *ErrList) ErrorRaise(err error) {
	if err != nil {
		e.ErrorAdd(err)
		panic(err.Error())
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
