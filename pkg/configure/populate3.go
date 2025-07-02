package configure

// import (
// 	"encoding/json"
// 	"flag"
// 	"fmt"
// 	"os"
// 	"reflect"
// 	"strconv"
// 	"strings"

// 	"gopkg.in/yaml.v3"
// )

// // PopulateStruct automatically populates a struct from JSON/YAML, env vars, and flags
// func PopulateStruct(cfg interface{}) error {
// 	loadConfigFile(cfg)

// 	v := reflect.ValueOf(cfg).Elem()
// 	t := v.Type()

// 	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
// 	fieldPointers := map[string]interface{}{}

// 	for i := 0; i < v.NumField(); i++ {
// 		field := v.Field(i)
// 		fieldType := t.Field(i)
// 		fieldName := fieldType.Name

// 		envKey := toUpperSnakeCase(fieldName)

// 		// Step 1: Read from ENV
// 		if envVal, exists := os.LookupEnv(envKey); exists {
// 			setField(field, envVal)
// 		}

// 		// Step 2: Set up flags
// 		flagKey := strings.ToLower(fieldName)
// 		switch field.Kind() {
// 		case reflect.String:
// 			ptr := new(string)
// 			fieldPointers[flagKey] = ptr
// 			flagSet.StringVar(ptr, flagKey, v.Field(i).String(), fmt.Sprintf("Set %s", flagKey))
// 		case reflect.Ptr:
// 			// Handle *string and *int
// 			if field.Type().Elem().Kind() == reflect.String {
// 				ptr := new(string)
// 				fieldPointers[flagKey] = ptr
// 				flagSet.StringVar(ptr, flagKey, "", fmt.Sprintf("Set %s", flagKey))
// 			} else if field.Type().Elem().Kind() == reflect.Int {
// 				ptr := new(int)
// 				fieldPointers[flagKey] = ptr
// 				flagSet.IntVar(ptr, flagKey, 0, fmt.Sprintf("Set %s", flagKey))
// 			}
// 		case reflect.Int:
// 			ptr := new(int)
// 			fieldPointers[flagKey] = ptr
// 			flagSet.IntVar(ptr, flagKey, int(v.Field(i).Int()), fmt.Sprintf("Set %s", flagKey))
// 		case reflect.Slice:
// 			// Handle []string and []int
// 			ptr := new(string)
// 			fieldPointers[flagKey] = ptr
// 			flagSet.StringVar(ptr, flagKey, "", fmt.Sprintf("Set %s", flagKey))
// 		}

// 		// Step 3: Apply default values
// 		if field.IsZero() {
// 			defaultValue := fieldType.Tag.Get("default")
// 			if defaultValue != "" {
// 				setField(field, defaultValue)
// 			}
// 		}
// 	}

// 	// Step 4: Parse CLI arguments
// 	flagSet.Parse(os.Args[2:])
// 	// fmt.Println(os.Args)

// 	// Step 5: Apply CLI flag values (override env & config)
// 	for key, ptr := range fieldPointers {
// 		field := v.FieldByNameFunc(func(name string) bool {
// 			return strings.ToLower(name) == key
// 		})
// 		if field.IsValid() {
// 			switch p := ptr.(type) {
// 			case *string:
// 				if *p != "" {
// 					setField(field, *p)
// 				}
// 			case *int:
// 				if *p != 0 {
// 					setField(field, strconv.Itoa(*p))
// 				}
// 			}
// 		}
// 	}

// 	return nil
// }

// // Convert "HostName" -> "HOST_NAME" for environment variables
// func toUpperSnakeCase(input string) string {
// 	var result strings.Builder
// 	for i, c := range input {
// 		if i > 0 && c >= 'A' && c <= 'Z' {
// 			result.WriteRune('_')
// 		}
// 		result.WriteRune(c)
// 	}
// 	return strings.ToUpper(result.String())
// }

// // Helper function to set struct field values
// func setField(field reflect.Value, value string) {
// 	switch field.Kind() {
// 	case reflect.String:
// 		field.SetString(value)
// 	case reflect.Ptr:
// 		// Handle *string and *int
// 		if field.Type().Elem().Kind() == reflect.String {
// 			field.Set(reflect.ValueOf(&value))
// 		} else if field.Type().Elem().Kind() == reflect.Int {
// 			if intVal, err := strconv.Atoi(value); err == nil {
// 				field.Set(reflect.ValueOf(&intVal))
// 			}
// 		}
// 	case reflect.Int:
// 		if intVal, err := strconv.Atoi(value); err == nil {
// 			field.SetInt(int64(intVal))
// 		}
// 	case reflect.Slice:
// 		elemType := field.Type().Elem().Kind()
// 		parts := strings.Split(value, ",")
// 		if elemType == reflect.String {
// 			field.Set(reflect.ValueOf(parts))
// 		} else if elemType == reflect.Int {
// 			ints := make([]int, len(parts))
// 			for i, p := range parts {
// 				if v, err := strconv.Atoi(strings.TrimSpace(p)); err == nil {
// 					ints[i] = v
// 				}
// 			}
// 			field.Set(reflect.ValueOf(ints))
// 		}
// 	}
// }

// // Load config file if present (JSON or YAML)
// func loadConfigFile(cfg interface{}) {
// 	files := []string{"config.json", "config.yaml"}
// 	for _, file := range files {
// 		if _, err := os.Stat(file); err == nil { // File exists
// 			data, err := os.ReadFile(file)
// 			if err != nil {
// 				fmt.Println("Error reading config file:", err)
// 				return
// 			}
// 			switch {
// 			case strings.HasSuffix(file, ".json"):
// 				json.Unmarshal(data, cfg)
// 			case strings.HasSuffix(file, ".yaml"):
// 				yaml.Unmarshal(data, cfg)
// 			}
// 			fmt.Println("Loaded config from", file)
// 			return
// 		}
// 	}
// }
