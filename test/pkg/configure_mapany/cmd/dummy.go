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
	add("GirlNamesAll", "gna", "jana,petra", "[]string", "",
		"Specified names", nil, nil)
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
	add("SourceDirectorySpecial", "sds", "", "string", c.NotNil,
		"Source file name", nil, helper.DirectoryExists)
	add("NumberSliceMap", "nsm", "10,12,13", "[]int", "",
		"number slice", nil, nil)
	// add("GirlNamesMap", "gnm", "jana,petra", "map[string]bool", "",
	add("GirlNamesMap", "gnm", "elvira,jaina", "[]string", "",
		// "Specified names", nil, nil)
		"Specified names", []string{"elvira", "jaina"}, nil)
	// "Specified names", nil, nil)
	// "Specified names", [][]string{{"elvira"}}, nil)
	add("GirlNamesSpecific", "gnms", "elvira,jaina", "[]string", "",
		"Specified names", []string{"elvira", "jaina", "renata"}, nil)
	add("GirlNameOne", "gno", "elvira,jaina", "string", "",
		"Specified names", []string{"elvira", "jaina", "renata"}, nil)
}

// add("GirlNamesMap", "gnm", "jana,petra", "map[string]bool", "",
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
	GirlNamesAll           []string
	DateFrom               time.Time
	Resume                 bool
	Count                  int
	NumberSlice            []int
	NumberSliceMap         map[int]bool
	FileName               string
	SourceDirectorySpecial string
	GirlNamesMap           map[string]bool
	GirlNamesSpecific      map[string]bool
	GirlNameOne            string
	// GirlNamesMap []string
	// ChoseVar        string
	// ChoseFunc       string
}

func RunCommandDummy() {
	cmdVars := commandDummyVars{}
	commanderDummyConfigure()
	commanderDummyConfig.SubcommandOptionsParse(&cmdVars)
	fmt.Printf("effective command vars %+v\n", cmdVars)
	// if cmdVars.Multiple == nil {
	// fmt.Println("is nil")
	// }
}
