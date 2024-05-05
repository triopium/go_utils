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
		"Source rundown file.", nil, nil)
	add("DateFrom", "df", "", "date",
		"Filter rundowns from date", nil, nil)
}

type commandDummyVars struct {
	SourceDirectory string
	DateFrom        time.Time
}

func RunCommandDummy() {
	cmdVars := commandDummyVars{}
	commandDummyConfigure()
	commandDummyConfig.RunSub(&cmdVars)
	fmt.Printf("effective command vars %+v\n", cmdVars)
}
