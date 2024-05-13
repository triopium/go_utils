package main

import (
	"fmt"
	"log/slog"

	"github.com/triopium/go_utils/pkg/helper"
	"github.com/triopium/go_utils/pkg/logging"
)

var err1 = fmt.Errorf("test error")

func Runner() (res error) {
	errList := new(helper.ErrList)
	defer errList.Handle(&res)
	errList.ErrorRaise(nil)
	errList.ErrorAdd(err1)
	errList.ErrorRaise(err1)
	errList.ErrorRaise(nil)
	return errList.ErrorsReturn()
}

func main() {
	logging.SetLogLevel("-4", "json")
	err := Runner()
	if err != nil {
		slog.Info("main error", "error", err.Error())
	}
}
