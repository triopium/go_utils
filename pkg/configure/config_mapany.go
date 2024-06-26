package configure

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/triopium/go_utils/pkg/helper"
	"github.com/triopium/go_utils/pkg/logging"
)

const (
	NotNil string = "NotNil"
)

type RootConfig struct {
	Usage       bool
	GeneralHelp bool
	Version     bool
	Verbose     int
	DryRun      bool
	LogType     string
	LogTest     bool
	DebugConfig bool
}

var CommanderRoot = CommanderConfig{
	Opts: []any{
		Opt[bool]{
			OptDesc{"Usage", "U", "false", "bool", "",
				"Print usage manual"}, nil, nil},
		Opt[bool]{
			OptDesc{"generalHelp", "H", "false", "bool", "",
				"Get help on all subcommands"}, nil, nil},
		Opt[bool]{
			OptDesc{"version", "V", "false", "bool", "",
				"Print version of program."}, nil, nil},
		Opt[int]{
			OptDesc{"verbose", "v", "0", "int", "",
				"Level of verbosity. Lower the number the more verbose is the output."}, []int{6, 4, 2, 0, -2, -4}, nil},
		Opt[string]{
			OptDesc{"logType", "logt", "json", "string", "",
				"Type of logs formating."},
			[]string{"json", "plain"}, nil},
		// nil, nil},
		Opt[bool]{
			OptDesc{"dryRun", "dr", "false", "bool", "",
				"Dry run, useful for tests. Avoid any pernament changes to filesystem or any expensive tasks"}, nil, nil},
		Opt[bool]{
			OptDesc{"debugConfig", "dc", "false", "bool", "",
				"Debug/print flag values"},
			nil, nil},
		Opt[bool]{
			OptDesc{"logTest", "logts", "false", "bool", "",
				"Print test logs"},
			nil, nil},
	},
}

type OptDesc struct {
	LongFlag  string
	ShortFlag string
	Default   string
	Type      string
	Spec      string
	Help      string
}

type Opt[T any] struct {
	OptDesc
	AllovedValues []T
	FuncMatch     func(v T) (bool, error)
}

type SubCommands map[string]func()

type CommanderConfig struct {
	OptsMap     map[string][6]interface{}
	Opts        []any
	Subs        SubCommands
	Values      interface{}
	VersionInfo interface{}
	*RootConfig
}

func (o *Opt[T]) Error(err error) {
	if err == nil {
		return
	}
	errMsgFormat := "cannot parse flag %s as type %s, err %v"
	errMsg := fmt.Errorf(errMsgFormat, o.LongFlag, o.Type, err)
	panic(errMsg)
}

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
		// if alloved != nil {
		// TODO
		// opt.AllovedValues = alloved.([][]string)
		// }
		cc.Opts = append(cc.Opts, opt)
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

type Subcommander map[string]func()

func (cc *CommanderConfig) Init() {
	cc.DeclareFlags()
	rcfg := &RootConfig{}
	flag.Parse()
	err := cc.ParseFlags(rcfg)
	if err != nil {
		panic(err)
	}
	logging.SetLogLevel(strconv.Itoa(rcfg.Verbose), rcfg.LogType)
	cc.RootConfig = rcfg
}

func (cc *CommanderConfig) PrintHelpForAllCommands() {
	fmt.Printf("## Root command:\n")
	path, is := IsCurrentExecutableBinary()
	fname := filepath.Base(path)
	cmd := ""
	if !is {
		cmd = fmt.Sprintf("go run %s.go -h", fname)
	}
	if is {
		cmd = fmt.Sprintf("./%s -h", fname)
	}
	fmt.Printf("  running source: `go run %s.go -h`\n", fname)
	fmt.Printf("  running compiled: `%s -h`\n\n", fname)
	cmds := strings.Split(cmd, " ")
	cmdexec := exec.Command(cmds[0], cmds[1:]...)
	resultLog, err := cmdexec.CombinedOutput()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	fmt.Println(string(resultLog))
	cmdsub := ""
	for i := range cc.Subs {
		fmt.Printf("## Subcommand: %s\n", i)
		fmt.Printf(
			"  running source: `go run %s.go %s -h`\n", fname, i)
		fmt.Printf(
			"  running compiled: `%s %s -h`\n\n", fname, i)
		if !is {
			cmdsub = fmt.Sprintf("go run %s.go %s -h", fname, i)
		}
		if is {
			cmdsub = fmt.Sprintf("%s %s -h", fname, i)
		}
		cmdsubs := strings.Split(cmdsub, " ")
		cmdexec := exec.Command(cmdsubs[0], cmdsubs[1:]...)
		resultLog, err := cmdexec.CombinedOutput()
		if err != nil {
			slog.Error(err.Error())
			return
		}
		fmt.Println(string(resultLog))
	}
}

func IsCurrentExecutableBinary() (string, bool) {
	expath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	base := filepath.Base(expath)
	return expath, base != "main"
}

func (cc *CommanderConfig) PrintManual() {
	fmt.Printf("# Help\n\n")
	cc.PrintHelpForAllCommands()
	fmt.Printf("# Usage\n\n")
}

func (cc *CommanderConfig) RootFlagsAct() {
	if cc.LogTest {
		logging.LoggingOutputTest()
	}
	if cc.Version {
		cc.VersionInfoPrint()
	}
	if cc.DebugConfig {
		fmt.Printf("Root config: %+v\n", cc.RootConfig)
	}
	if cc.GeneralHelp {
		cc.PrintHelpForAllCommands()
	}
	if cc.Usage {
		cc.PrintManual()
	}
	if flag.NArg() < 1 {
		os.Exit(0)
	}
}

func (cc *CommanderConfig) RunRoot() {
	cc.RootFlagsAct()
	subCmdName := flag.Arg(0)
	if subCmdName == "" {
		return
	}
	slog.Debug("calling sub", "name", subCmdName)
	subCmdFunc, ok := cc.Subs[subCmdName]
	if !ok {
		panic(fmt.Errorf("unknown subcommand: %s", subCmdName))
	}
	subCmdFunc()
}

func (cc *CommanderConfig) AddSub(subName string, subF func()) {
	slog.Debug("subcommand added", "subname", subName)
	if cc.Subs == nil {
		cc.Subs = make(SubCommands)
	}
	cc.Subs[subName] = subF
}

func (cc *CommanderConfig) RunSub(intf interface{}) {
	subcmd := flag.Arg(0)
	slog.Debug("subcommand called", "subcommand", subcmd)
	// FlagsUsage = fmt.Sprintf("subcommand: %s\n", subcmd)
	FlagsUsage = ""
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

func (cc *CommanderConfig) VersionInfoAdd(info interface{}) {
	cc.VersionInfo = info
}

func (cc *CommanderConfig) VersionInfoPrint() {
	fmt.Printf("%+v\n", cc.VersionInfo)
}

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

func (cc *CommanderConfig) DeclareFlags() {
	slog.Debug("declaring flags", "count", len(cc.Opts))
	cc.OptsMap = make(map[string][6]interface{})
	FlagsUsage += fmt.Sprintf(
		"-h, -help\n\t%s\n\n",
		"display this help and exit")
	for i := range cc.Opts {
		optstr := fmt.Sprintf("%v", cc.Opts[i])
		slog.Debug("declaring flag", "opt", optstr)
		op := cc.Opts[i]
		DeclareFlagHandle[any](op, cc.OptsMap)
	}
	flag.Usage = Usage
}

func (o *Opt[T]) DeclareUsage() {
	slog.Debug("declare usage")
	fd := o.OptDesc
	if o.AllovedValues == nil {
		format := "-%s, -%s=%s\n\t%s\n\n"
		FlagsUsage += fmt.Sprintf(
			format, fd.ShortFlag, fd.LongFlag, fd.Default, fd.Help)
	} else {
		format := "-%s, -%s=%s\n\t%s\n\t%v\n\n"
		FlagsUsage += fmt.Sprintf(
			format, fd.ShortFlag, fd.LongFlag, fd.Default, fd.Help, o.AllovedValues)
	}
}

func (cc *CommanderConfig) ParseFlags(iface interface{}) error {
	slog.Debug("parsing flags")
	vof := reflect.ValueOf(iface)
	if vof.Kind() != reflect.Ptr || vof.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("Invalid input: not a pointer to a struct")
	}

	vofElem := vof.Elem()
	n := vofElem.NumField()
	slog.Debug("parsing flags", "count", n)
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

func CheckAllovedVals(
	flagName string, inp any, alloved interface{}, allovedFunc func(any) (bool, error)) (bool, error) {
	// var match bool
	if allovedFunc != nil {
		return allovedFunc(inp)
	}
	return true, nil
}

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

type CheckerUntyped[T any] struct {
	AllovedVals any
	AllovedFunc any
}

func (ch *CheckerUntyped[T]) CheckAlloved(opt, inp any) {
	if ch.AllovedFunc != nil {
		_, err := ch.AllovedFunc.(func(T) (bool, error))(inp.(T))
		if err != nil {
			panic(err)
		}
	}
}

// type Checker[T any] struct {
type Checker[T comparable] struct {
	AllovedVals any
	AllovedFunc any
}

// func (ch *Checker[T]) CheckAlloved(inp T) {
func (ch *Checker[T]) CheckAlloved(opt string, inp any) {
	if ch.AllovedFunc != nil {
		ok, err := ch.AllovedFunc.(func(T) (bool, error))(inp.(T))
		if err != nil {
			panic(err)
		}
		if !ok {
			funcName := runtime.FuncForPC(
				reflect.ValueOf(ch.AllovedFunc).Pointer()).Name()

			panic(fmt.Errorf(
				"opt: %s, value: %v not alloved by allowFunc: %v",
				opt, inp, funcName))
		}
	}
	if ch.AllovedVals == nil {
		return
	}
	var match bool
	if ch.AllovedVals != nil {
		vals := ch.AllovedVals.([]T)
		for _, v := range vals {
			ipnv := inp.(T)
			if v == ipnv {
				match = true
				break
			}
			// res1 := reflect.DeepEqual(v, ipnv)
		}
	}
	if !match {
		panic(
			fmt.Errorf("value not alloved by allowedVals: %v, aloved: %v",
				inp, ch.AllovedVals))
	}
}
