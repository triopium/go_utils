package cmd

import (
	"fmt"

	"github.com/triopium/go_utils/pkg/configure"
	"github.com/triopium/go_utils/pkg/helper"
)

var commanderDummyConfig = configure.CommanderConfig{}

func commanderDummyConfigure() {
	add := commanderDummyConfig.AddOption2
	add("SourceDirectory", "srcDir", "/tmp", "string",
		"Source directory", nil, helper.DirectoryExists)
	// "Source directory", []any{"jak", "tak"}, helper.DirectoryExists)
	// opt := configure.Opt[string]{}
	// add("DateFrom", "df", "", "date",
	// "Filter rundowns from date", nil, nil)
	// add("Multiple", "m", "", "[]string",
	// "Multiple choices", nil, nil)
	// add("ChoseVar", "chv", "", "string",
	// "Filter rundowns from date", []any{"kek", "lek"}, nil)
	// add("ChoseFunc", "chf", "", "string",
	// "Filter rundowns from date", nil, ChooseFunction)
}

func ChooseFunction(in any) bool {
	return len(in.(string)) > 2
}

type commandDummyVars struct {
	SourceDirectory string
	// DateFrom        time.Time
	// ChoseVar        string
	// ChoseFunc       string
	// Multiple        []string
}

func RunCommandDummy() {
	cmdVars := commandDummyVars{}
	commanderDummyConfigure()
	commanderDummyConfig.RunSub(&cmdVars)
	fmt.Printf("effective command vars %+v\n", cmdVars)
	// if cmdVars.Multiple == nil {
	// fmt.Println("is nil")
	// }
}
