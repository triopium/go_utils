package configure

import (
	"flag"
	"log/slog"
	"strconv"
	"time"

	"github.com/triopium/go_utils/pkg/helper"
)

func DeclareFlagHandle[T any](
	s interface{}, myMap map[string][6]interface{}) {
	var def, long, short, alloved, funcMatch interface{}
	var optName string
	var spec string
	switch o := s.(type) {
	case Opt[bool]:
		b, err := strconv.ParseBool(o.Default)
		o.Error(err)
		def = &b
		// long = flag.Bool(o.LongFlag, b, o.Help)
		// short = flag.Bool(o.ShortFlag, b, o.Help)
		long = flag.Bool(o.LongFlag, false, o.Help)
		short = flag.Bool(o.ShortFlag, false, o.Help)
		if o.AllovedValues != nil {
			alloved = o.AllovedValues
		}
		if o.FuncMatch != nil {
			funcMatch = o.FuncMatch
		}
		optName = helper.FirstLetterToUppercase(o.LongFlag)
		o.DeclareUsage()
	case Opt[int]:
		b, err := strconv.Atoi(o.Default)
		o.Error(err)
		def = &b
		long = flag.Int(o.LongFlag, b, o.Help)
		short = flag.Int(o.ShortFlag, b, o.Help)
		alloved = o.AllovedValues
		funcMatch = o.FuncMatch
		if o.AllovedValues != nil {
			alloved = o.AllovedValues
		}
		if o.FuncMatch != nil {
			funcMatch = o.FuncMatch
		}
		optName = helper.FirstLetterToUppercase(o.LongFlag)
		o.DeclareUsage()
	case Opt[[]int]:
		def = o.Default
		long = flag.String(o.LongFlag, o.Default, o.Help)
		short = flag.String(o.ShortFlag, o.Default, o.Help)
		if o.AllovedValues != nil {
			alloved = o.AllovedValues
		}
		if o.FuncMatch != nil {
			funcMatch = o.FuncMatch
		}
		optName = helper.FirstLetterToUppercase(o.LongFlag)
		o.DeclareUsage()
	case Opt[string]:
		slog.Debug("declaring flag", "optname", o.LongFlag)
		def = o.Default
		long = flag.String(o.LongFlag, o.Default, o.Help)
		short = flag.String(o.ShortFlag, o.Default, o.Help)
		if o.AllovedValues != nil {
			alloved = o.AllovedValues
		}
		if o.FuncMatch != nil {
			funcMatch = o.FuncMatch
		}
		spec = o.Spec
		optName = helper.FirstLetterToUppercase(o.LongFlag)
		o.DeclareUsage()
	case Opt[[]string]:
		def = o.Default
		long = flag.String(o.LongFlag, o.Default, o.Help)
		short = flag.String(o.ShortFlag, o.Default, o.Help)
		if o.AllovedValues != nil {
			alloved = o.AllovedValues
		}
		if o.FuncMatch != nil {
			funcMatch = o.FuncMatch
		}
		optName = helper.FirstLetterToUppercase(o.LongFlag)
		o.DeclareUsage()
	case Opt[time.Time]:
		def = o.Default
		long = flag.String(o.LongFlag, o.Default, o.Help)
		short = flag.String(o.ShortFlag, o.Default, o.Help)
		if o.AllovedValues != nil {
			alloved = o.AllovedValues
		}
		if o.FuncMatch != nil {
			funcMatch = o.FuncMatch
		}
		optName = helper.FirstLetterToUppercase(o.LongFlag)
		o.DeclareUsage()
	default:
		panic("no match")
	}
	myMap[optName] = [6]interface{}{
		def, long, short, spec, alloved, funcMatch}
}
