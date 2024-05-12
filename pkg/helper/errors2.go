package helper

import (
	"errors"
	"fmt"
	"log/slog"
)

type ErrMain struct {
	Errors []error
}

type Errw struct {
	BaseErr      error
	Count        int
	SubErrsSlice []*Errw
	SubErrsMap   map[error]*Errw
}

func (er *Errw) addNew(err error) {
	eW := new(Errw)
	eW.BaseErr = err
	eW.Count++
	er.SubErrsSlice = append(er.SubErrsSlice, eW)
	er.SubErrsMap[err] = eW
}

func (er *Errw) Add(err error) {
	// init
	if er.SubErrsSlice == nil {
		er.SubErrsSlice = make([]*Errw, 0, 3)
		er.SubErrsMap = make(map[error]*Errw)
		er.addNew(err)
		return
	}

	// Chek if main already added
	_, ok := er.SubErrsMap[err]
	if ok {
		er.SubErrsMap[err].Count++
		return
	}
	er.addNew(err)
}

func (er *Errw) ErrorsReturn() error {
	var res error
	for _, sub := range er.SubErrsSlice {
		out := fmt.Errorf("%d: %w", sub.Count, sub.BaseErr)
		res = errors.Join(res, out)
		// err := fmt.Errorf("%w; %w; %w", err, err2, err3)
		// return errors.Join(e.Errors...)
	}
	return res
}

func (er *Errw) Wrap(formatString string, a ...any) (ret error) {
	return fmt.Errorf(formatString, a...)
}

func (er *Errw) MaxUwrap(err error) (ret error) {
	if err == nil {
		slog.Info("is nil")
		return nil
	}
	out := errors.Unwrap(err)
	if out == nil {
		slog.Info("is nil")
		return err
	}
	for i := 0; i < 20; i++ {
		ret = out
		out = errors.Unwrap(out)
		if out == nil {
			return ret
		}
	}
	return ret
}

type ErrorAggregate map[string]int // function,line,count of errors

// type count, variables

func (e *ErrMain) Handle(err *error) {
	if r := recover(); r != nil {
		slog.Error(r.(string))
		// *err = fmt.Errorf("kek")
		*err = e.ErrorsReturn()
	}
}

func (e *ErrMain) ErrorRaise(err error) {
	if err != nil {
		e.ErrorAdd(err)
		panic(err.Error() + "fek")
		// panic(e.ErrorsReturn().Error())
	} else {
		slog.Info("not raised")
	}
}

func (e *ErrMain) ErrorAdd(err error) {
	e.Errors = append(e.Errors, err)
}

func (e *ErrMain) ErrorsReturn() error {
	// err := fmt.Errorf("%w; %w; %w", err, err2, err3)
	if len(e.Errors) != 0 {
		return errors.Join(e.Errors...)
	}
	return nil
}
