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
// 	// Step 1: Load JSON/YAML config file (if present)
// 	loadConfigFile(cfg)

// 	v := reflect.ValueOf(cfg).Elem()
// 	t := v.Type()

// 	// Parse flags before populating
// 	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
// 	fieldPointers := map[string]interface{}{}

// 	// Iterate over struct fields
// 	for i := 0; i < v.NumField(); i++ {
// 		field := v.Field(i)
// 		fieldType := t.Field(i)
// 		fieldName := fieldType.Name

// 		// Convert field name to ENV format (UPPER_SNAKE_CASE)
// 		envKey := toUpperSnakeCase(fieldName)

// 		// Step 2: Read from ENV variables
// 		if envVal, exists := os.LookupEnv(envKey); exists {
// 			setField(field, envVal)
// 		}

// 		// Step 3: Set up command-line flags
// 		flagKey := strings.ToLower(fieldName)
// 		switch field.Kind() {
// 		case reflect.String:
// 			ptr := new(string)
// 			fieldPointers[flagKey] = ptr
// 			flagSet.StringVar(ptr, flagKey, v.Field(i).String(), fmt.Sprintf("Set %s", flagKey))
// 		case reflect.Int:
// 			ptr := new(int)
// 			fieldPointers[flagKey] = ptr
// 			flagSet.IntVar(ptr, flagKey, int(v.Field(i).Int()), fmt.Sprintf("Set %s", flagKey))
// 		}

// 		// Step 4: Apply default values if still zero
// 		if field.IsZero() {
// 			defaultValue := fieldType.Tag.Get("default")
// 			if defaultValue != "" {
// 				setField(field, defaultValue)
// 			}
// 		}
// 	}

// 	// Step 5: Parse CLI arguments
// 	flagSet.Parse(os.Args[1:])

// 	// Step 6: Apply CLI flag values (override env & file)
// 	for key, ptr := range fieldPointers {
// 		field := v.FieldByNameFunc(func(name string) bool {
// 			return strings.ToLower(name) == key
// 		})
// 		if field.IsValid() {
// 			switch p := ptr.(type) {
// 			case *string:
// 				setField(field, *p)
// 			case *int:
// 				setField(field, strconv.Itoa(*p))
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
// 	case reflect.Int:
// 		if intVal, err := strconv.Atoi(value); err == nil {
// 			field.SetInt(int64(intVal))
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
