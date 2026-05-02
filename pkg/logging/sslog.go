package logging

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

type Where struct {
	File  string
	Fname string
	Line  int
}

func (w Where) String() string {
	return fmt.Sprintf(
		"%s:%s:%d",
		w.File,
		w.Fname,
		w.Line,
	)
}
func (w Where) StringPrettyPrint() {
	b, err1 := json.MarshalIndent(w, "", "  ")
	if err1 != nil {
		log.Fatal(err1)
	}
	fmt.Println(string(b))
}

func WhereLevel(level int) (Where, error) {
	pc, file, line, ok := runtime.Caller(level)
	// pc, file, line, ok := runtime.Caller(3)
	if !ok {
		return Where{
			File:  file,
			Fname: "???",
			Line:  line,
		}, fmt.Errorf("cannot stat")
	}

	fn := runtime.FuncForPC(pc)
	file = constructLogFileNamePath(file, fn.Name())
	fnName := strings.Join(strings.Split(filepath.Base(fn.Name()), ".")[1:], ".")

	return Where{
		File:  file,
		Fname: fnName,
		Line:  line,
	}, nil
}

func constructLogFileNamePath(full, fn string) string {
	res := filepath.Base(fn)
	ress := strings.Split(res, ".")
	// childPkg := ress[0]
	// fmt.Println("full", full)
	fname := filepath.Base(full)
	// fmt.Println("fname", full)
	// fmt.Println("childPkg", childPkg)
	fnpkg := filepath.Dir(fn)
	// fmt.Println("fnpkg", fnpkg)
	fnpkg = StripDirLevelsFromLeft(fnpkg, 2)
	// fmt.Println("fnpkgs", fnpkg)
	return filepath.Join("./", fnpkg, ress[0], fname)

}

func StripDirLevelsFromLeft(full string, levels int) string {
	full = path.Clean(full)
	full = filepath.ToSlash(full)
	parts := strings.Split(full, "/")

	// remove empty parts (leading "/")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p != "" {
			out = append(out, p)
		}
	}

	if levels >= len(out) {
		return "/"
	}

	out = out[levels:]

	return "/" + strings.Join(out, "/")
}

func WhereLevelPrint(level int) {
	pc, file, line, ok := runtime.Caller(level)
	if !ok {
		return
	}
	fn := runtime.FuncForPC(pc)

	fmt.Printf("%s:%d %s\n", file, line, fn.Name())
}

// ConfigLogger debug print parsed config, halt after print bool?
func ConfigLogger(input any, halt ...bool) {
	w, _ := WhereLevel(2)
	outputFile := os.Stderr
	name := reflect.TypeOf(input)
	_, err := fmt.Fprintf(outputFile, "\n%s:%s: %+v\n", w, name, input)
	if err != nil {
		log.Fatal(err)
	}
	b, err1 := json.MarshalIndent(input, "", "  ")
	if err1 != nil {
		log.Fatal(err1)
	}
	_, err2 := fmt.Fprintf(
		outputFile, "\n%s %s %s\n", w, name, string(b),
	)
	if err2 != nil {
		log.Fatal(err2)
	}
	if len(halt) > 0 {
		os.Exit(0)
	}
}
