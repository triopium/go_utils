package helper

import (
	"fmt"
	"testing"
)

func TestErrMap_Add(t *testing.T) {
	myErrs := ErrMap{}
	type args struct {
		errMain       error
		errorsPartial []error
	}
	tests := []struct {
		name string
		em   ErrMap
		args args
		want bool
	}{
		{"one", myErrs, args{fmt.Errorf("main err"), nil}, true},
		{"two", myErrs, args{nil, nil}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.em.Add(tt.args.errMain, tt.args.errorsPartial...); got != tt.want {
				t.Errorf("ErrMap.Add() = %v, want %v", got, tt.want)
			}
			fmt.Println("after", tt.em)
		})
	}
}
