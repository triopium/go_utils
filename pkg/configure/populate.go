package configure

// import (
// 	"flag"
// 	"fmt"
// 	"os"
// 	"reflect"
// 	"strconv"
// 	"strings"
// )

// // PopulateStruct automatically populates a struct from env vars and flags
// func PopulateStruct(cfg interface{}) error {
// 	v := reflect.ValueOf(cfg).Elem()
// 	t := v.Type()

// 	// Parse flags before populating
// 	flagSet := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
// 	fieldPointers := map[string]interface{}{} // Store field pointers for flag parsing

// 	// First, process struct fields
// 	for i := 0; i < v.NumField(); i++ {
// 		field := v.Field(i)
// 		fieldType := t.Field(i)
// 		fieldName := fieldType.Name

// 		// Convert struct field name to ENV variable format (UPPER_SNAKE_CASE)
// 		envKey := toUpperSnakeCase(fieldName)

// 		// Read from ENV variables if present
// 		if envVal, exists := os.LookupEnv(envKey); exists {
// 			setField(field, envVal)
// 		}

// 		// Set up command-line flags
// 		flagKey := strings.ToLower(fieldName) // Flag format: lowercase (e.g., "port")
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

// 		// Apply default values if field is still zero
// 		if field.IsZero() {
// 			defaultValue := fieldType.Tag.Get("default")
// 			if defaultValue != "" {
// 				setField(field, defaultValue)
// 			}
// 		}
// 	}

// 	// Parse CLI arguments
// 	flagSet.Parse(os.Args[1:])

// 	// Apply CLI flag values (override env/defaults)
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

// // Helper function to set field values dynamically
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
