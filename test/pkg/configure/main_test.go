package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestMainCommand(t *testing.T) {
	curDir, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(curDir)
	binar := filepath.Join(curDir, "main.go")
	cmd := exec.Command(binar)
	err = cmd.Start()
	if err != nil {
		t.Error(err)
	}
}
