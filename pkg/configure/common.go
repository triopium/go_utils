package configure

import (
	"fmt"

	"github.com/triopium/go_utils/pkg/helper"
)

// VersionInfo
type VersionInfo struct {
	ProgramName string
	Version     string
	GitTag      string
	GitCommit   string
	BuildTime   string
}

var FlagsUsage = "Usage:\n"

// Usage called when help command invoked
func Usage() {
	fmt.Println(FlagsUsage)
}

// OPTION GETERS
// GetBoolValuePriority return value according to priority. Priority is given in desceding. Last value is default value.
func GetBoolValuePriority(vals ...bool) bool {
	count := len(vals) - 1
	res := vals[count]
	for i := count - 1; i >= 0; i-- {
		res = helper.XOR(res, vals[i])
	}
	return res
}

func GetIntValuePriority(vals ...int) int {
	count := len(vals) - 1
	def := vals[count]
	res := def
	for i := count - 1; i >= 0; i-- {
		if vals[i] != def {
			res = vals[i]
		}
	}
	return res
}

func GetStringValuePriority(vals ...string) string {
	count := len(vals) - 1
	def := vals[count]
	res := def
	for i := count - 1; i >= 0; i-- {
		if vals[i] != def {
			res = vals[i]
		}
	}
	return res
}

func GetStringValueByPriority(
	longFlagValue, shortFlagValue, envValue, defaultValue string) string {
	res := defaultValue
	if longFlagValue != defaultValue {
		res = longFlagValue
	}
	if shortFlagValue != defaultValue {
		res = shortFlagValue
	}
	return res
}

func GetBoolValueByPriority(
	longFlagValue, shortFlagValue, envValue, defaultValue bool) bool {
	res := helper.XOR(defaultValue, shortFlagValue)
	res = helper.XOR(res, longFlagValue)
	return res
}

func GetIntValueByPriority(
	longFlagValue, shortFlagValue, envValue, defaultValue int) int {
	res := defaultValue
	if longFlagValue != defaultValue {
		res = longFlagValue
	}
	if shortFlagValue != defaultValue {
		res = shortFlagValue
	}
	return res
}
