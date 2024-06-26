package configure

import (
	"flag"
	"fmt"
	"log/slog"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/triopium/go_utils/pkg/helper"
	"github.com/triopium/go_utils/pkg/logging"
)

// TODO: print usage for all commands at once using run main.go and all subcmds with help
// TODO: implement validate flag value by function (check alloved values)

type RootCfg struct {
	Version     bool
	Verbose     int
	DryRun      bool
	LogType     string
	DebugConfig bool
}

var CommandRoot = CommandConfig{
	Opts: []FlagOption{
		{FlagDescription{"version", "V", "false", "bool",
			"Version of program."},
			nil, nil},

		{FlagDescription{"verbose", "v", "0", "int",
			"Level of verbosity."},
			[]any{-4, -3, -2, -1, 0, 1, 2, 3, 4}, nil},

		{FlagDescription{"dryRun", "dr", "false", "bool",
			"Dry run, useful for tests."},
			nil, nil},

		{FlagDescription{"logType", "logt", "json", "string",
			"Type of logs formating."},
			[]any{"json", "plain"}, nil},

		{FlagDescription{"debugConfig", "dc", "false", "bool",
			"Debug/print flag values"},
			[]any{"json", "plain"}, nil},
	},
}

type CommandConfig struct {
	OptsMap     map[string][6]interface{}
	Opts        []FlagOption
	Subs        Subcommands
	Values      interface{}
	VersionInfo interface{}
}

// type functor[V any](bool,error)
// type Func func[V int|string](v V) bool

type FlagOption struct {
	FlagDescription
	AllovedValues []any
	// FuncMatch     func[V int|string](v V) bool
	FuncMatch func(any) bool
	// FuncMatch     func(any) (bool, error) `json:"-"`
	// FuncMatch     ConfigVariableFunctionCheck `json:"-"`
}

type FlagDescription struct {
	LongFlag   string
	ShortFlag  string
	Default    string
	Type       string
	Descripton string
}

type Subcommands map[string]func()

func (cc *CommandConfig) Init() {
	cc.DeclareFlags()
	rcfg := &RootCfg{}
	flag.Parse()
	err := cc.ParseFlags(rcfg)
	if err != nil {
		panic(err)
	}
	logging.SetLogLevel(strconv.Itoa(rcfg.Verbose), rcfg.LogType)
	cc.Values = rcfg
	if flag.NArg() < 1 {
		cc.VersionInfoPrint()
		return
	}
	slog.Info("root config", "config", cc.Values)
}

func (cc *CommandConfig) VersionInfoAdd(info interface{}) {
	cc.VersionInfo = info
}

func (cc *CommandConfig) VersionInfoPrint() {
	fmt.Printf("%+v\n", cc.VersionInfo)
}

func (cc *CommandConfig) RunRoot() {
	subCmdName := flag.Arg(0)
	if subCmdName == "" {
		return
	}
	subCmdFunc, ok := cc.Subs[subCmdName]
	if !ok {
		panic(fmt.Errorf("unknown subcommand: %s", subCmdName))
	}
	subCmdFunc()
}

func (cc *CommandConfig) RunSub(intf interface{}) {
	subcmd := flag.Arg(0)
	slog.Info("subcommand called", "subcommand", subcmd)
	FlagsUsage = fmt.Sprintf("subcommand: %s\n", subcmd)
	cc.DeclareFlags()
	err := flag.CommandLine.Parse(flag.Args()[1:])
	if err != nil {
		panic(err)
	}
	err = cc.ParseFlags(intf)
	if err != nil {
		panic(err)
	}
}

func (cc *CommandConfig) AddSub(subName string, subF func()) {
	if cc.Subs == nil {
		cc.Subs = make(Subcommands)
	}
	cc.Subs[subName] = subF
}

func (cc *CommandConfig) AddOption(
	long, short, defValue, typeValue, descr string,
	alloved []any, cfgVarFunc func(any) bool) {
	opt := FlagOption{
		FlagDescription: FlagDescription{
			long, short, defValue, typeValue, descr,
		}, AllovedValues: alloved, FuncMatch: cfgVarFunc}
	cc.Opts = append(cc.Opts, opt)
}

func (opt FlagDescription) Error(err error) {
	if err == nil {
		return
	}
	errMsgFormat := "cannot parse flag %s as type %s, err %v"
	errMsg := fmt.Errorf(errMsgFormat, opt.LongFlag, opt.Type, err)
	panic(errMsg)
}

func CheckAllovedValues(
	flagName string, inp any, alloved interface{}, allovedFunc func(any) bool) {
	var match bool
	// var res any = inp
	if allovedFunc != nil {
		if !allovedFunc(inp) {
			panic(inp)
		}
	}

	switch t := alloved.(type) {
	case []interface{}:
		if len(t) == 0 {
			return
		}
		for _, i := range alloved.([]interface{}) {
			if inp == i {
				match = true
				break
			}
		}
	default:
		err := fmt.Errorf("unknow type of alloved definition: %v", t)
		panic(err)
	}
	if !match {
		err := fmt.Errorf("flag '-%s=%v' not alloved, alloved values: %v",
			flagName, inp, alloved)
		panic(err)
	}
}

func (opt *FlagOption) DeclareUsage() {
	fd := opt.FlagDescription
	if opt.AllovedValues == nil {
		format := "-%s, -%s\n\t%s\n\n"
		FlagsUsage += fmt.Sprintf(format, fd.ShortFlag, fd.LongFlag, fd.Descripton)
	} else {
		format := "-%s, -%s\n\t%s\n\t%v\n\n"
		FlagsUsage += fmt.Sprintf(format, fd.ShortFlag, fd.LongFlag, fd.Descripton, opt.AllovedValues)
	}
}

func (cc *CommandConfig) DeclareFlags() {
	slog.Debug("declaring flags")
	cc.OptsMap = make(map[string][6]interface{})
	for i := range cc.Opts {
		slog.Debug("declaring flag", "opt", cc.Opts[i])
		res := cc.Opts[i].DeclareFlag()
		name := cc.Opts[i].LongFlag
		name = helper.FirstLetterToUppercase(name)
		cc.OptsMap[name] = res
	}
	flag.Usage = Usage
}

func (opt *FlagOption) DeclareFlag() [6]interface{} {
	var def, long, short interface{}
	opt.DeclareUsage()
	switch opt.FlagDescription.Type {
	case "bool":
		b, err := strconv.ParseBool(opt.Default)
		opt.Error(err)
		def = &b
		long = flag.Bool(opt.LongFlag, b, opt.Descripton)
		short = flag.Bool(opt.ShortFlag, b, opt.Descripton)
	case "int":
		b, err := strconv.Atoi(opt.Default)
		opt.Error(err)
		def = &b
		long = flag.Int(opt.LongFlag, b, opt.Descripton)
		short = flag.Int(opt.ShortFlag, b, opt.Descripton)
	case "string":
		def = opt.Default
		long = flag.String(opt.LongFlag, opt.Default, opt.Descripton)
		short = flag.String(opt.ShortFlag, opt.Default, opt.Descripton)
	case "[]string":
		def = opt.Default
		long = flag.String(opt.LongFlag, opt.Default, opt.Descripton)
		short = flag.String(opt.ShortFlag, opt.Default, opt.Descripton)
	case "date":
		def = opt.Default
		long = flag.String(opt.LongFlag, opt.Default, opt.Descripton)
		short = flag.String(opt.ShortFlag, opt.Default, opt.Descripton)
	default:
		err := fmt.Errorf("unknow flag type: %s", opt.FlagDescription.Type)
		opt.Error(err)
	}
	return [6]interface{}{def, long, short, "", opt.AllovedValues, opt.FuncMatch}
}

func (cc *CommandConfig) ParseFlags(iface interface{}) error {
	slog.Info("parsing flags")
	vof := reflect.ValueOf(iface)
	if vof.Kind() != reflect.Ptr || vof.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("Invalid input: not a pointer to a struct")
	}

	vofe := vof.Elem()
	n := vofe.NumField()
	for i := 0; i < n; i++ {
		field := vofe.Type().Field(i)
		optName := helper.FirstLetterToUppercase(field.Name)
		vals, ok := cc.OptsMap[optName]
		if !ok {
			slog.Debug("flag not defined for struct field", "field", optName)
			continue
		}
		slog.Debug("parsing flag", "name", field.Name, "type", field.Type.Name())
		def := vals[0]
		long := vals[1]
		short := vals[2]
		allovedVars := vals[4]
		allovedFunc := vals[5]
		// switch field.Type.Name() {
		// v := field.Interface()
		v := vofe.Field(i).Interface()
		switch v.(type) {
		case bool:
			vals := []bool{*long.(*bool), *short.(*bool), *def.(*bool)}
			res := GetBoolValuePriority(vals...)
			vofe.Field(i).SetBool(res)
		case int:
			vals := []int{*long.(*int), *short.(*int), *def.(*int)}
			res := GetIntValuePriority(vals...)
			CheckAllovedValues(
				optName, res, allovedVars, allovedFunc.(func(any) bool))
			vofe.Field(i).SetInt(int64(res))
		case string:
			vals := []string{*long.(*string), *short.(*string), def.(string)}
			res := GetStringValuePriority(vals...)
			CheckAllovedValues(
				optName, res, allovedVars, allovedFunc.(func(any) bool))
			vofe.Field(i).SetString(res)
		case []string:
			vals := []string{*long.(*string), *short.(*string), def.(string)}
			res := GetStringValuePriority(vals...)
			if res == "" {
				continue
			}
			strSlice := strings.Split(res, ",")
			rv := reflect.ValueOf(strSlice)
			vofe.Field(i).Set(rv)
			// v := slice.Index(0)
			// v.SetString("kek")
			// slice := reflect.MakeSlice(reflect.TypeOf([]string{}), 2, 2)
			// reflect.Copy(slice, rv)
		case time.Time:
			vals := []string{*long.(*string), *short.(*string), def.(string)}
			res := GetStringValuePriority(vals...)
			if res == "" {
				continue
			}
			date, err := helper.ParseStringDate(res, time.Local)
			if err != nil {
				panic(err)
			}
			vofe.Field(i).Set(reflect.ValueOf(date))
		default:
			panic(fmt.Errorf("flag type not implemented: %s", field.Type))
		}
	}
	return nil
}
