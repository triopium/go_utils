package configure

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/triopium/go_utils/pkg/logging"
	"gopkg.in/yaml.v3"
)

var (
	timeType      = reflect.TypeOf(time.Time{})
	durationType  = reflect.TypeOf(time.Duration(0))
	flagValueType = reflect.TypeOf((*flag.Value)(nil)).Elem()
)

// PopulateStruct automatically populates a struct from JSON/YAML, env vars, and flags
func PopulateStruct(cfg any, debugPrint bool) error {
	loadConfigFile(cfg)

	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()
	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	// Enable debug logging if "DEBUG" env is set
	debug = os.Getenv("DEBUG") == "true"
	// fieldPointers := map[string]any{}

	// Step 1: Iterate over struct fields
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		fieldName := fieldType.Name
		envKey := toUpperSnakeCase(fieldName)
		if !field.CanSet() {
			continue
		}

		// Step 2: Read from ENV
		if envVal, exists := os.LookupEnv(envKey); exists {
			logSet(fieldName, envVal, "ENV "+envKey)
			setField(field, envVal)
		}

		// Step 3: Set up flags
		flagName := strings.ToLower(fieldName)
		bindValue(flagSet, flagName, field, fieldType, SetUsage(flagName))
	}

	// Step 5: Parse CLI arguments
	// err := flagSet.Parse(os.Args[2:])
	// flagSet.Parse(os.Args[2:])
	// if err != nil {
	// panic(err)
	// }

	// Step 6: Log parsed config
	if debugPrint {
		logging.ConfigLogger(cfg)
	}
	return nil
}

func bindValue(
	fs *flag.FlagSet,
	flagName string,
	field reflect.Value,
	fieldType reflect.StructField,
	usage string,
) {
	// Step 4: Apply default values
	if field.IsZero() {
		// fmt.Println("fuck zero [22:44:47]", flagKey, field.String())
		defaultValue := fieldType.Tag.Get("envDefault")
		if defaultValue != "" {
			logSet(flagName, defaultValue, "DEFAULT")
			setField(field, defaultValue)
		}
	}

	// 2. Special known types
	switch field.Type() {
	case durationType:
		ptr := field.Addr().Interface().(*time.Duration)

		fmt.Println("fuck dur [13:52:20]", ptr.String())
		fs.DurationVar(ptr, flagName, *ptr, usage)
		return
	}

	switch field.Kind() {
	case reflect.Ptr:
		// 1. Allocate if nil
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		// 2. Recurse on the value it points to
		bindValue(fs, flagName, field.Elem(), fieldType, usage)

	case reflect.Bool:
		ptr := field.Addr().Interface().(*bool)
		fs.BoolVar(ptr, flagName, *ptr, usage)

	case reflect.Int:
		ptr := field.Addr().Interface().(*int)
		fs.IntVar(ptr, flagName, *ptr, usage)

	case reflect.Float64:
		ptr := field.Addr().Interface().(*float64)
		fs.Float64Var(ptr, flagName, *ptr, usage)

	case reflect.String:
		ptr := field.Addr().Interface().(*string)
		fs.StringVar(ptr, flagName, *ptr, usage)

	case reflect.Struct:
		switch field.Type() {
		case timeType:
			fs.Var(
				&timeValue{t: field.Addr().Interface().(*time.Time)},
				flagName, usage)
		}

	default:
		// fmt.Println("fuck [23:07:56]", name)
		// panic("unsupported kind: " + field.Kind().String())
	}
}

type timeValue struct {
	t *time.Time
}

func (v *timeValue) Set(s string) error {
	parsed, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return err
	}
	*v.t = parsed
	return nil
}

func (v *timeValue) String() string {
	if v.t == nil || v.t.IsZero() {
		return ""
	}
	return v.t.Format(time.RFC3339)
}

// Load config file if present (JSON or YAML)
func loadConfigFile(cfg any) {
	files := []string{"config.json", "config.yaml"}
	for _, file := range files {
		if _, err := os.Stat(file); err == nil {
			data, err := os.ReadFile(file)
			if err != nil {
				fmt.Println("Error reading config file:", err)
				return
			}
			switch {
			case strings.HasSuffix(file, ".json"):
				json.Unmarshal(data, cfg)
			case strings.HasSuffix(file, ".yaml"):
				yaml.Unmarshal(data, cfg)
			}
			fmt.Println("[CONFIG] Loaded from", file)
			return
		}
	}
}

// type timeValue struct {
// 	field  reflect.Value
// 	layout string
// }
// func newTimeValue(field reflect.Value) *timeValue {
// 	return &timeValue{
// 		field:  field,
// 		layout: time.RFC3339,
// 	}
// }
