package helper

import (
	"errors"
	"fmt"
	"log/slog"
	"testing"
)

func TestErrWrapped(t *testing.T) {
}

func TestErrw_Wrap(t *testing.T) {
	err1 := fmt.Errorf("base error")
	baseFormat := "%02d: %w"
	errw := Errw{}
	i := 0
	res := errw.Wrap(baseFormat, i, err1)
	for ; i < 10; i++ {
		res = errw.Wrap(baseFormat, i, res)
		fmt.Println(res)
		if errors.Is(res, err1) {
			fmt.Println(i, "is orig error")
		}
	}
	fmt.Println("kak", res)
	out := errw.MaxUwrap(res)
	fmt.Println("kek", out)
}

func TestErrw_Add(t *testing.T) {
	err1 := fmt.Errorf("foo error")
	err2 := fmt.Errorf("bar error")
	ew := new(Errw)
	ew.Add(err1, "kek", "tek")
	ew.Add(err1)
	ew.Add(err1)
	ew.Add(err1, "toker")
	ew.Add(err2)
	ew.Add(err2)
	err := ew.ErrorsReturn()
	fmt.Println(err)
	slog.Error(err.Error())
	fmt.Println()
	err = ew.ErrorsReturnMap()
	fmt.Println(err)
	slog.Error(err.Error())
}

func TestErrJoin(t *testing.T) {
	err1 := fmt.Errorf("foo error")
	err2 := fmt.Errorf("bar error")
	err := errors.Join(err1, err2)
	fmt.Println(err)
	slog.Info("errors", "errs", err)
	slog.Info("errors", "errs", []string{"ahoj", "mahoj"})
	slog.Info("errors", "errs", map[string]int{"kek": 10})
	fmt.Println("fuck", errors.Is(err1, err2))
	fmt.Println("fuck", errors.Is(err1, err1))
}
