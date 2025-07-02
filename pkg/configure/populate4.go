package configure

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

var debug bool // Global debug mode

type CommonOpts struct {
	Verbose int
	Help    bool
}

// PopulateStruct automatically populates a struct from JSON/YAML, env vars, and flags
func PopulateStruct(cfg interface{}) error {
	loadConfigFile(cfg)

	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	fieldPointers := map[string]interface{}{}

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
		case reflect.String:
			ptr := new(string)
			fieldPointers[flagKey] = ptr
			flagSet.StringVar(ptr, flagKey, field.String(), fmt.Sprintf("Set %s", flagKey))
		case reflect.Ptr:
			if field.Type().Elem().Kind() == reflect.String {
				ptr := new(string)
				fieldPointers[flagKey] = ptr
				flagSet.StringVar(ptr, flagKey, "", fmt.Sprintf("Set %s", flagKey))
			} else if field.Type().Elem().Kind() == reflect.Int {
				ptr := new(int)
				fieldPointers[flagKey] = ptr
				flagSet.IntVar(ptr, flagKey, 0, fmt.Sprintf("Set %s", flagKey))
			}
		case reflect.Int:
			ptr := new(int)
			fieldPointers[flagKey] = ptr
			flagSet.IntVar(ptr, flagKey, int(field.Int()), fmt.Sprintf("Set %s", flagKey))
		case reflect.Slice:
			ptr := new(string)
			fieldPointers[flagKey] = ptr
			flagSet.StringVar(ptr, flagKey, "", fmt.Sprintf("Set %s", flagKey))
		}

		// Step 4: Apply default values
		if field.IsZero() {
			defaultValue := fieldType.Tag.Get("envDefault")
			if defaultValue != "" {
				logSet(fieldName, defaultValue, "DEFAULT")
				setField(field, defaultValue)
			}
		}
	}

	// Step 5: Parse CLI arguments
	flagSet.Parse(os.Args[2:])

	// Step 6: Apply CLI flag values (override previous values)
	for key, ptr := range fieldPointers {
		field := v.FieldByNameFunc(func(name string) bool {
			return strings.ToLower(name) == key
		})
		if field.IsValid() {
			switch p := ptr.(type) {
			case *string:
				if *p != "" {
					logSet(field.Type().Name(), *p, "CLI -"+key)
					setField(field, *p)
				}
			case *int:
				if *p != 0 {
					logSet(field.Type().Name(), strconv.Itoa(*p), "CLI -"+key)
					setField(field, strconv.Itoa(*p))
				}
			}
		}
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

// Load config file if present (JSON or YAML)
func loadConfigFile(cfg interface{}) {
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
