package configure

import (
	"time"
)

func (cc *CommanderConfig) AddOptionSimple(o any) {
	cc.Opts = append(cc.Opts, o)
}

func (cc *CommanderConfig) AddOption(
	long, short, defValue, typeValue, spec, descr string,
	alloved, funcMatch any) {
	switch typeValue {
	case "bool":
		opt := Opt[bool]{}
		opt.LongFlag = long
		opt.ShortFlag = short
		opt.Default = defValue
		opt.Type = typeValue
		opt.Help = descr
		opt.Spec = spec
		if funcMatch != nil {
			opt.FuncMatch = funcMatch.(func(bool) (bool, error))
		}
		cc.Opts = append(cc.Opts, opt)
	case "int":
		opt := Opt[int]{}
		opt.LongFlag = long
		opt.ShortFlag = short
		opt.Default = defValue
		opt.Type = typeValue
		opt.Help = descr
		opt.Spec = spec
		if funcMatch != nil {
			opt.FuncMatch = funcMatch.(func(int) (bool, error))
		}
		if alloved != nil {
			opt.AllovedValues = alloved.([]int)
		}
		cc.Opts = append(cc.Opts, opt)
	case "[]int":
		opt := Opt[[]int]{}
		opt.LongFlag = long
		opt.ShortFlag = short
		opt.Default = defValue
		opt.Type = typeValue
		opt.Help = descr
		opt.Spec = spec
		if funcMatch != nil {
			opt.FuncMatch = funcMatch.(func([]int) (bool, error))
		}
		cc.Opts = append(cc.Opts, opt)
	case "string":
		opt := Opt[string]{}
		opt.LongFlag = long
		opt.ShortFlag = short
		opt.Default = defValue
		opt.Type = typeValue
		opt.Help = descr
		opt.Spec = spec
		if funcMatch != nil {
			opt.FuncMatch = funcMatch.(func(string) (bool, error))
		}
		if alloved != nil {
			opt.AllovedValues = alloved.([]string)
		}
		cc.Opts = append(cc.Opts, opt)
	case "[]string":
		opt := Opt[[]string]{}
		opt.LongFlag = long
		opt.ShortFlag = short
		opt.Default = defValue
		opt.Type = typeValue
		opt.Help = descr
		opt.Spec = spec
		if funcMatch != nil {
			opt.FuncMatch = funcMatch.(func([]string) (bool, error))
		}
		if alloved != nil {
			opts := alloved.([]string)
			opt.AllovedValues = [][]string{opts}
		}
		cc.Opts = append(cc.Opts, opt)
	case "map[string]bool":
		opt := Opt[map[string]bool]{}
		opt.LongFlag = long
		opt.ShortFlag = short
		opt.Default = defValue
		opt.Type = typeValue
		opt.Help = descr
		opt.Spec = spec
		if funcMatch != nil {
			opt.FuncMatch = funcMatch.(func(map[string]bool) (bool, error))
		}
	case "date":
		opt := Opt[time.Time]{}
		opt.LongFlag = long
		opt.ShortFlag = short
		opt.Default = defValue
		opt.Type = typeValue
		opt.Help = descr
		opt.Spec = spec
		if funcMatch != nil {
			opt.FuncMatch = funcMatch.(func(time.Time) (bool, error))
		}
		if alloved != nil {
			opt.AllovedValues = alloved.([]time.Time)
		}
		cc.Opts = append(cc.Opts, opt)
	}
}
