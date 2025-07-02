package helper

import (
	"fmt"
	"regexp"
	"testing"
)

func TestRegexp(t *testing.T) {
	pattern := "^13:00-14:00"
	regex := regexp.MustCompile(pattern)
	testString := "13:00-14:00"
	ok := regex.MatchString(testString)
	fmt.Println("matches", ok, pattern, testString)

	testString = "kek 13:00-14:00"
	ok = regex.MatchString(testString)
	fmt.Println("matches", ok, pattern, testString)
}

func TestEscapeCSVdelim(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"one", args{"hello	"}, "hello\\t"},
		{"two", args{"hello		"}, "hello\\t\\t"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EscapeCSVdelim(tt.args.value); got != tt.want {
				t.Errorf("EscapeCSVdelim() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEscapeCSVdelimB(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"one", args{"hello\t"}, "helloTAB"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EscapeCSVdelimB(tt.args.value); got != tt.want {
				t.Errorf("EscapeCSVdelimB() = %v, want %v", got, tt.want)
			}
		})
	}
}
