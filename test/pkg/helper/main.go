package main

import (
	"fmt"
	"log/slog"

	"github.com/triopium/go_utils/pkg/helper"
	"github.com/triopium/go_utils/pkg/logging"
)

var err1 = fmt.Errorf("runner error")

func Runner() (res error) {
	errMain := new(helper.ErrMain)
	defer errMain.Handle(&res)
	errMain.ErrorRaise(nil)
	errMain.ErrorAdd(err1)
	errMain.ErrorRaise(err1)
	errMain.ErrorRaise(nil)
	return errMain.ErrorsReturn()
}

func main() {
	logging.SetLogLevel("-4", "json")
	err := Runner()
	if err != nil {
		slog.Info("fuck", "error", err.Error())
	}
}
