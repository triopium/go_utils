package helper

import (
	"errors"
	"fmt"
	"log/slog"
)

type Errw struct {
	BaseErr          error
	ErrsDetailsSlice []any
	ErrsDetailsMap   map[string]any
	Count            int
	SubErrsSlice     []*Errw
	SubErrsMap       map[error]*Errw
}

func (er *Errw) addNew(err error, params ...any) {
	eW := new(Errw)
	eW.BaseErr = err
	eW.Count++
	eW.ErrsDetailsSlice = make([]any, 0)
	eW.ErrsDetailsSlice = append(eW.ErrsDetailsSlice, params)
	eW.ErrsDetailsMap = make(map[string]any)
	er.SubErrsSlice = append(er.SubErrsSlice, eW)
	er.SubErrsMap[err] = eW
}

func (er *Errw) Add(err error, params ...any) {
	// init
	if er.SubErrsSlice == nil {
		er.SubErrsSlice = make([]*Errw, 0, 3)
		er.SubErrsMap = make(map[error]*Errw)
		er.addNew(err, params...)
		return
	}

	// Chek if main already added
	_, ok := er.SubErrsMap[err]
	if ok {
		er.SubErrsMap[err].Count++
		er.SubErrsMap[err].ErrsDetailsSlice = append(er.SubErrsMap[err].ErrsDetailsSlice, params)
		return
	}
	er.addNew(err, params...)
}

func (er *Errw) ErrorsReturn() error {
	var res error
	for _, sub := range er.SubErrsSlice {
		if len(sub.ErrsDetailsSlice) == 0 {
			out := fmt.Errorf("%d: %w", sub.Count, sub.BaseErr)
			res = errors.Join(res, out)
			continue
		}
		out := fmt.Errorf("%d: %w: %v", sub.Count, sub.BaseErr, sub.ErrsDetailsSlice)
		res = errors.Join(res, out)
		// err := fmt.Errorf("%w; %w; %w", err, err2, err3)
		// return errors.Join(e.Errors...)
	}
	return res
}

func (er *Errw) ErrorsReturnMap() error {
	var res error
	for _, sub := range er.SubErrsSlice {
		sub.ErrsDetailsMap["kek"] = "jak"
		sub.ErrsDetailsMap["sek"] = "tek"
		if len(sub.ErrsDetailsMap) == 0 {
			out := fmt.Errorf("%d: %w", sub.Count, sub.BaseErr)
			res = errors.Join(res, out)
			continue
		}
		out := fmt.Errorf("%d: %w: %v", sub.Count, sub.BaseErr, sub.ErrsDetailsMap)
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
