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
	"testing"

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
	GoName        string
	BinName       string
	FlagsDeclared bool
}

func (o *Opt[T]) Error(err error) {
	if err == nil {
		return
	}
	errMsgFormat := "cannot parse flag %s as type %s, err %v"
	errMsg := fmt.Errorf(errMsgFormat, o.LongFlag, o.Type, err)
	panic(errMsg)
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

func CommandConstruct(flags string) string {
	path, is := IsCurrentExecutableBinary()
	fname := filepath.Base(path)
	cmd := ""
	if !is {
		cmd = fmt.Sprintf("go run %s.go %s", fname, flags)
	}
	if is {
		cmd = fmt.Sprintf("./%s %s", fname, flags)
	}
	return cmd
}

func PrintCommandExample(goName, binName, flags string) {
	path, is := IsCurrentExecutableBinary()
	fname := filepath.Base(path)
	if testing.Testing() {
		fmt.Printf("running from source:\n\n")
		fmt.Printf("```\ngo run %s.go %s\n```", goName, flags)
		fmt.Printf("\n\n")
		fmt.Printf("running compiled:\n\n")
		fmt.Printf("```\n./%s %s\n```", binName, flags)
		fmt.Printf("\n\n")
		return
	}
	if !is {
		fmt.Println("running from source:")
		fmt.Printf("```\ngo run %s.go %s\n```\n", fname, flags)
		fmt.Printf("\n")
		fmt.Printf("running compiled:\n")
	}
	wd, err := os.Getwd()
	fname = filepath.Base(wd)
	if err != nil {
		panic(err)
	}
	fmt.Printf("```\n./%s %s\n```\n", fname, flags)
	fmt.Printf("\n")
}

func GetCommand(flags string) string {
	path, is := IsCurrentExecutableBinary()
	fname := filepath.Base(path)
	var cmd string
	if is {
		cmd = fmt.Sprintf("./%s %s", fname, flags)
	} else {
		cmd = fmt.Sprintf("go run %s.go -h", fname)
	}
	return cmd
}

func (cc *CommanderConfig) PrintHelpForAllCommands(
	goName, binName string) {
	fmt.Printf("## Root command:\n")
	PrintCommandExample(goName, binName, "-h")
	cmd := ""
	cmd = GetCommand("-h")
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
		PrintCommandExample(goName, binName, i+" -h")
		cmdsub = GetCommand("-h")
		cmdsubs := strings.Split(cmdsub, " ")
		cmdexec := exec.Command(cmdsubs[0], cmdsubs[1:]...)
		resultLog, err := cmdexec.CombinedOutput()
		if err != nil {
			slog.Error(err.Error())
			return
		}
		fmt.Printf("%s\n", string(resultLog))
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
	cc.PrintHelpForAllCommands(cc.GoName, cc.BinName)
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
		cc.PrintHelpForAllCommands(cc.GoName, cc.BinName)
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

func (cc *CommanderConfig) SubcommandOptionsParse(intf interface{}) {
	subcmd := flag.Arg(0)
	slog.Debug("subcommand called", "subcommand", subcmd)
	// FlagsUsage = fmt.Sprintf("subcommand: %s\n", subcmd)
	FlagsUsage = ""
	if !cc.FlagsDeclared {
		cc.DeclareFlags()
		err := flag.CommandLine.Parse(flag.Args()[1:])
		if err != nil {
			panic(err)
		}
		cc.FlagsDeclared = true
	}
	err := cc.ParseFlags(intf)
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
	if ch.AllovedVals != nil {
		switch vals := inp.(type) {
		case []string:
			allVals := ch.AllovedVals.([]string)
			allValsMap := helper.SliceStringToMapString(allVals)
			for _, str := range vals {
				_, ok := allValsMap[str]
				if !ok {
					panic(fmt.Errorf(
						"value not alloved, value: %q, alloved values: %v",
						str, allVals))
				}
			}
		default:
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
