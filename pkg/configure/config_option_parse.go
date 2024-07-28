package configure

import (
	"fmt"
	"log/slog"
	"reflect"
	"strings"
	"time"

	"github.com/triopium/go_utils/pkg/helper"
)

func (cc *CommanderConfig) ParseFlag(
	optName string, vofe reflect.Value, index int) error {
	var ok bool
	var allovedFunc interface{}
	var allovedVars interface{}
	vals, ok := cc.OptsMap[optName]
	if !ok {
		slog.Debug(
			"flag not defined for struct field", "field", optName)
		return nil
	}
	def := vals[0]
	long := vals[1]
	short := vals[2]
	spec := vals[3].(string)
	if vals[4] == nil {
		allovedVars = nil
	} else {
		allovedVars = vals[4]
	}
	if vals[5] == nil {
		allovedFunc = nil
	} else {
		allovedFunc = vals[5]
	}
	slog.Debug("parsing flag", "name", optName)
	v := vofe.Field(index).Interface()
	switch v.(type) {
	case bool:
		vals := []bool{*long.(*bool), *short.(*bool), *def.(*bool)}
		res := GetBoolValuePriority(vals...)
		if vals[2] && vals[1] {
			res = false
		}
		if vals[2] && vals[0] {
			res = false
		}
		vofe.Field(index).SetBool(res)
	case string:
		valsp := []string{*long.(*string), *short.(*string), def.(string)}
		res := GetStringValuePriority(valsp...)
		if spec == NotNil && res == "" {
			panic(fmt.Errorf("flag: %s value cannot be empty", optName))
		}
		if res == "" {
			return nil
		}
		ch := Checker[string]{allovedVars, allovedFunc}
		ch.CheckAlloved(optName, res)
		vofe.Field(index).SetString(res)
	case []string:
		vals := []string{*long.(*string), *short.(*string), def.(string)}
		res := GetStringValuePriority(vals...)
		strSlice := strings.Split(res, ",")
		ch := CheckerUntyped[[]string]{allovedVars, allovedFunc}
		ch.CheckAlloved(optName, strSlice)
		rv := reflect.ValueOf(strSlice)
		vofe.Field(index).Set(rv)
	case map[string]bool:
		vals := []string{*long.(*string), *short.(*string), def.(string)}
		res := GetStringValuePriority(vals...)
		strValues := strings.Split(res, ",")

		ch := CheckerUntyped[[]string]{allovedVars, allovedFunc}
		switch av := allovedVars.(type) {
		case [][]string:
			if len(av) == 1 {
				ch = CheckerUntyped[[]string]{av[0], allovedFunc}
			}
		case nil:
			ch = CheckerUntyped[[]string]{av, allovedFunc}
		}
		ch.CheckAlloved(optName, strValues)
		strValuesMap := helper.SliceStringToMapString(strValues)
		rv := reflect.ValueOf(strValuesMap)
		vofe.Field(index).Set(rv)
	case int:
		vals := []int{*long.(*int), *short.(*int), *def.(*int)}
		res := GetIntValuePriority(vals...)
		vofe.Field(index).SetInt(int64(res))
	case []int:
		vals := []string{*long.(*string), *short.(*string), def.(string)}
		res := GetStringValuePriority(vals...)
		out, err := helper.StringToIntSlice(res, ",")
		if err != nil {
			panic(fmt.Errorf("%w: %s", err, res))
		}
		ch := CheckerUntyped[[]int]{allovedVars, allovedFunc}
		ch.CheckAlloved(optName, out)
		rv := reflect.ValueOf(out)
		vofe.Field(index).Set(rv)
	case map[int]bool:
		vals := []string{*long.(*string), *short.(*string), def.(string)}
		res := GetStringValuePriority(vals...)
		strSlice := strings.Split(res, ",")
		strMap := helper.SliceStringToMapInt(strSlice)
		rv := reflect.ValueOf(strMap)
		vofe.Field(index).Set(rv)
	case time.Time:
		vals := []string{*long.(*string), *short.(*string), def.(string)}
		res := GetStringValuePriority(vals...)
		ch := Checker[time.Time]{allovedVars, allovedFunc}
		date, err := helper.ParseStringDate(res, time.Local)
		if err != nil {
			panic(fmt.Errorf("%w: %s", err, res))
		}
		ch.CheckAlloved(optName, date)
		vofe.Field(index).Set(reflect.ValueOf(date))
	default:
		return fmt.Errorf("unknow flag type: %T", v)
	}
	return nil
}
