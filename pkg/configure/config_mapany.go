package configure

import (
	"flag"
	"fmt"
	"log/slog"
	"reflect"
	"strconv"

	"github.com/triopium/go_utils/pkg/helper"
)

type RootConfig struct {
	SourceDirectory string
	Version         bool
	Verbose         int
	DryRun          bool
	LogType         string
	DebugConfig     bool
}

type OptDesc struct {
	LongFlag   string
	ShortFlag  string
	Default    string
	Type       string
	Descripton string
}

type Opt[T any] struct {
	OptDesc
	Alloved   []any
	FuncMatch func(v T) (bool, error)
}

type CommanderConfig struct {
	OptsMap     map[string][6]interface{}
	Opts        []any
	Subs        Subcommands
	Values      interface{}
	VersionInfo interface{}
}

func (o *Opt[T]) Error(err error) {
	if err == nil {
		return
	}
	errMsgFormat := "cannot parse flag %s as type %s, err %v"
	errMsg := fmt.Errorf(errMsgFormat, o.LongFlag, o.Type, err)
	panic(errMsg)
}

func (cc *CommanderConfig) AddOption(o any) {
	cc.Opts = append(cc.Opts, o)
}

var CommanderRoot = CommanderConfig{
	Opts: []any{
		Opt[bool]{
			OptDesc{"version", "V", "false", "bool",
				"Version of program."}, nil, nil},
		Opt[int]{
			OptDesc{"verbose", "v", "0", "int",
				"Level of verbosity."}, nil, nil},
		Opt[bool]{
			OptDesc{"dryRun", "dr", "false", "bool",
				"Dry run, useful for tests."}, nil, nil},
		// Opt[string]{
		// 	OptDesc{
		// 		"SourceDirectory", "srcDir", "", "string",
		// 		"source directory"}, nil, helper.DirectoryExists,
		// },
	},
}

type Subcommander map[string]func()

func (cc *CommanderConfig) Init() {
	cc.DeclareFlags()
	rcfg := &RootConfig{}
	flag.Parse()
	err := cc.ParseFlags(rcfg)
	if err != nil {
		panic(err)
	}
	// logging.SetLogLevel(strconv.Itoa(rcfg.Verbose), rcfg.LogType)
	// cc.Values = rcfg
	// if flag.NArg() < 1 {
	// 	cc.VersionInfoPrint()
	// 	return
	// }
	// slog.Info("root config", "config", cc.Values)
	slog.Info("root config", "config", rcfg)
}

func (cc *CommanderConfig) VersionInfoAdd(info interface{}) {
	cc.VersionInfo = info
}

func (cc *CommanderConfig) VersionInfoPrint() {
	fmt.Printf("%+v\n", cc.VersionInfo)
}

func DeclareFlagHandle[T any](
	s interface{}, myMap map[string][6]interface{}) [6]interface{} {
	var def, long, short, alloved, funcMatch interface{}
	var optName string
	switch o := s.(type) {
	case Opt[bool]:
		b, err := strconv.ParseBool(o.Default)
		o.Error(err)
		def = &b
		long = flag.Bool(o.LongFlag, b, o.Descripton)
		short = flag.Bool(o.ShortFlag, b, o.Descripton)
		optName = helper.FirstLetterToUppercase(o.LongFlag)
	case Opt[int]:
		b, err := strconv.Atoi(o.Default)
		o.Error(err)
		def = &b
		long = flag.Int(o.LongFlag, b, o.Descripton)
		short = flag.Int(o.ShortFlag, b, o.Descripton)
		alloved = o.Alloved
		funcMatch = o.FuncMatch
		optName = helper.FirstLetterToUppercase(o.LongFlag)
	case Opt[string]:
		def = o.Default
		long = flag.String(o.LongFlag, o.Default, o.Descripton)
		short = flag.String(o.ShortFlag, o.Default, o.Descripton)
		optName = helper.FirstLetterToUppercase(o.LongFlag)
	default:
		panic("no match")
	}
	myMap[optName] = [6]interface{}{def, long, short, "", alloved, funcMatch}
	return [6]interface{}{def, long, short, "", alloved, funcMatch}
}

func (cc *CommanderConfig) DeclareFlags() {
	slog.Info("declaring flags", "count", len(cc.Opts))
	cc.OptsMap = make(map[string][6]interface{})
	for i := range cc.Opts {
		slog.Info("declaring flag", "opt", cc.Opts[i])
		op := cc.Opts[i]
		DeclareFlagHandle[any](op, cc.OptsMap)
	}
	// flag.Usage = Usage
}

func (cc *CommanderConfig) ParseFlags(iface interface{}) error {
	slog.Info("parsing flags")
	vof := reflect.ValueOf(iface)
	if vof.Kind() != reflect.Ptr || vof.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("Invalid input: not a pointer to a struct")
	}

	vofElem := vof.Elem()
	n := vofElem.NumField()
	slog.Info("parsing flags", "count", n)
	for i := 0; i < n; i++ {
		field := vofElem.Type().Field(i)
		optName := helper.FirstLetterToUppercase(field.Name)
		err := cc.ParseFlag(optName, vofElem, i)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cc *CommanderConfig) ParseFlag(
	optName string, vofe reflect.Value, index int) error {
	var ok bool
	var err error
	var value any
	vals, ok := cc.OptsMap[optName]
	if !ok {
		slog.Info(
			"flag not defined for struct field", "field", optName)
		return nil
	}
	def := vals[0]
	long := vals[1]
	short := vals[2]
	// allovedVars := vals[4]
	allovedFunc := vals[5]
	slog.Info("parsing flag", "name", optName)
	v := vofe.Field(index).Interface()
	switch v.(type) {
	case bool:
		vals := []bool{*long.(*bool), *short.(*bool), *def.(*bool)}
		res := GetBoolValuePriority(vals...)
		vofe.Field(index).SetBool(res)
	case string:
		vals := []string{*long.(*string), *short.(*string), def.(string)}
		res := GetStringValuePriority(vals...)
		fun := allovedFunc.(func(string) (bool, error))
		ok, err = fun(res)
		if ok {
			vofe.Field(index).SetString(res)
		}
		value = res
	case int:
		vals := []int{*long.(*int), *short.(*int), *def.(*int)}
		res := GetIntValuePriority(vals...)
		slog.Info("parsing flag", "name", optName, "value", res)
		vofe.Field(index).SetInt(int64(res))
	default:
		return fmt.Errorf("unknow flag type: %T", v)
	}

	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("falg value not alloved: %v", value)
	}
	return nil
}
