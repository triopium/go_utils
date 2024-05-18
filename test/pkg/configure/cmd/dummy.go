package cmd

import (
	"fmt"
	"time"

	"github.com/triopium/go_utils/pkg/configure"
)

var commandDummyConfig = configure.CommandConfig{}

func commandDummyConfigure() {
	add := commandDummyConfig.AddOption
	add("SourceDirectory", "srcdir", "", "string",
		// "Source directory must exists.", nil, helper.DirectoryExists)
		"Source directory must exists.", nil, nil)
	add("DateFrom", "df", "", "date",
		"Filter rundowns from date", nil, nil)
	add("Multiple", "m", "", "[]string",
		"Multiple choices", nil, nil)
	add("ChoseVar", "chv", "", "string",
		"Filter rundowns from date", []any{"kek", "lek"}, nil)
	add("ChoseFunc", "chf", "", "string",
		"Filter rundowns from date", nil, ChooseFunction)
}

func ChooseFunction(in any) bool {
	return len(in.(string)) > 2
}

type commandDummyVars struct {
	SourceDirectory string
	DateFrom        time.Time
	ChoseVar        string
	ChoseFunc       string
	Multiple        []string
}

func RunCommandDummy() {
	cmdVars := commandDummyVars{}
	commandDummyConfigure()
	commandDummyConfig.RunSub(&cmdVars)
	fmt.Printf("effective command vars %+v\n", cmdVars)
	if cmdVars.Multiple == nil {
		fmt.Println("is nil")
	}
}
