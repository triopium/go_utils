package cmd

import (
	"fmt"
	"time"

	"github.com/triopium/go_utils/pkg/configure"
	"github.com/triopium/go_utils/pkg/helper"
)

var commanderDummyConfig = configure.CommanderConfig{}

func commanderDummyConfigure() {
	add := commanderDummyConfig.AddOption2
	add("SourceDirectory", "srcDir", "/tmp", "string",
		"Source directory", []string{"/tmp", "/home"}, helper.DirectoryExists)
	add("GirlNames", "gn", "jana,petra", "[]string",
		"Specified names", nil, AllovedNames)
	add("DateFrom", "df", "2020", "date", "date from", nil, nil)
	add("Resume", "re", "true", "bool", "should resume?", nil, nil)
	add("Count", "cn", "10", "int", "count", nil, nil)
	add("NumberSlice", "ns", "10,12", "[]int", "number slice", nil, nil)
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

func AllovedNames(input []string) (bool, error) {
	alloved := []string{"jana", "petra", "klara"}
	allovedMap := make(map[string]bool)
	for _, n := range alloved {
		allovedMap[n] = true
	}
	for _, i := range input {
		if !allovedMap[i] {
			err := fmt.Errorf("value not alloved: %s", i)
			return false, err
		}
	}
	return true, nil
}

func ChooseFunction(in any) bool {
	return len(in.(string)) > 2
}

type commandDummyVars struct {
	SourceDirectory string
	GirlNames       []string
	DateFrom        time.Time
	Resume          bool
	Count           int
	NumberSlice     []int
	// GirlNames       []string
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
