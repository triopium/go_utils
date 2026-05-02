package configure

import (
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/triopium/go_utils/pkg/logging"
)

var debug bool // Global debug mode

type CommonOpts struct {
	Verbose int
	Help    bool
}

func PopulateStructFatal(cfg any, debugPrint bool) {
	err := PopulateStruct(cfg, debugPrint)
	if err != nil {
		log.Fatal(err)
	}
}

func SetUsage(name string) string {
	usage := fmt.Sprintf("Set %s", name)
	return usage
}

// PopulateStruct automatically populates a struct from JSON/YAML, env vars, and flags
func PopulateStructWorks(cfg any, debugPrint bool) error {
	loadConfigFile(cfg)

	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	fieldPointers := map[string]any{}

	// Enable debug logging if "DEBUG" env is set
	debug = os.Getenv("DEBUG") == "true"

	// Step 1: Iterate over struct fields
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		fieldName := fieldType.Name
		envKey := toUpperSnakeCase(fieldName)

		// Step 2: Read from ENV
		if envVal, exists := os.LookupEnv(envKey); exists {
			logSet(fieldName, envVal, "ENV "+envKey)
			setField(field, envVal)
		}

		// Step 3: Set up flags
		flagKey := strings.ToLower(fieldName)
		switch field.Kind() {
		case reflect.Bool:
			ptr := field.Addr().Interface().(*bool)
			fieldPointers[flagKey] = ptr
			flagSet.BoolVar(ptr, flagKey, field.Bool(), SetUsage(flagKey))

		case reflect.String:
			ptr := field.Addr().Interface().(*string)
			fieldPointers[flagKey] = ptr
			flagSet.StringVar(ptr, flagKey, field.String(), SetUsage(flagKey))

		case reflect.Ptr:
			if field.IsNil() {
				field.Set(reflect.New(field.Type().Elem()))
			}
			bindValue(flagSet, flagKey, field.Elem(), fieldType, SetUsage(flagKey))

		case reflect.Int:
			ptr := field.Addr().Interface().(*int)
			fieldPointers[flagKey] = ptr
			flagSet.IntVar(ptr, flagKey, int(field.Int()), SetUsage(flagKey))
		case reflect.Slice:
			fmt.Println("fuck slice [22:39:51]")
			ptr := new(string)
			fieldPointers[flagKey] = ptr
			flagSet.StringVar(ptr, flagKey, "", SetUsage(flagKey))
		}

		// Step 4: Apply default values
		if field.IsZero() {

			fmt.Println("fuck zero [22:44:47]", flagKey, field.String())
			defaultValue := fieldType.Tag.Get("envDefault")
			if defaultValue != "" {
				logSet(fieldName, defaultValue, "DEFAULT")
				setField(field, defaultValue)
			}
		}
	}

	// Step 5: Parse CLI arguments
	flagSet.Parse(os.Args[2:])

	if debugPrint {
		logging.ConfigLogger(cfg)
	}
	return nil
}

// Logs configuration sources
func logSet(field, value, source string) {
	if debug {
		fmt.Printf("[CONFIG] %s = %q (%s)\n", field, value, source)
	}
}

// Convert "HostName" -> "HOST_NAME" for environment variables
func toUpperSnakeCase(input string) string {
	var result strings.Builder
	for i, c := range input {
		if i > 0 && c >= 'A' && c <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(c)
	}
	return strings.ToUpper(result.String())
}

// Helper function to set struct field values
func setField(field reflect.Value, value string) {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Ptr:
		if field.Type().Elem().Kind() == reflect.String {
			field.Set(reflect.ValueOf(&value))
		} else if field.Type().Elem().Kind() == reflect.Int {
			if intVal, err := strconv.Atoi(value); err == nil {
				field.Set(reflect.ValueOf(&intVal))
			}
		}
	case reflect.Int:
		if intVal, err := strconv.Atoi(value); err == nil {
			field.SetInt(int64(intVal))
		}
	case reflect.Slice:
		elemType := field.Type().Elem().Kind()
		parts := strings.Split(value, ",")
		if elemType == reflect.String {
			field.Set(reflect.ValueOf(parts))
		} else if elemType == reflect.Int {
			ints := make([]int, len(parts))
			for i, p := range parts {
				if v, err := strconv.Atoi(strings.TrimSpace(p)); err == nil {
					ints[i] = v
				}
			}
			field.Set(reflect.ValueOf(ints))
		}
	}
}
