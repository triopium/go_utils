package cmd

import (
	"fmt"
	"time"

	c "github.com/triopium/go_utils/pkg/configure"
	"github.com/triopium/go_utils/pkg/helper"
)

var commanderDummyConfig = c.CommanderConfig{}

func commanderDummyConfigure() {
	add := commanderDummyConfig.AddOption
	add("SourceDirectory", "srcDir", "/tmp", "string", "",
		"Source directory", []string{"/tmp", "/home"}, helper.DirectoryExists)
	add("GirlNames", "gn", "jana,petra", "[]string", "",
		"Specified names", nil, AllovedNames)
	add("DateFrom", "df", "2020", "date", "",
		"date from", nil, nil)
	add("Resume", "re", "true", "bool", "",
		"should resume?", nil, nil)
	add("Count", "cn", "10", "int", "",
		"count", nil, nil)
	add("NumberSlice", "ns", "10,12", "[]int", "",
		"number slice", nil, nil)
	add("FileName", "fn", "", "string", "",
		"Source file name", nil, helper.FileExists)
	add("SourceDirectorySpecial", "sds", "", "string", c.NotNill,
		"Source file name", nil, helper.DirectoryExists)
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
	SourceDirectory        string
	GirlNames              []string
	DateFrom               time.Time
	Resume                 bool
	Count                  int
	NumberSlice            []int
	FileName               string
	SourceDirectorySpecial string
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
