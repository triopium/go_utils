package configure

// Step 6: Apply CLI flag values (override previous values)
// for key, ptr := range fieldPointers {
// 	field := v.FieldByNameFunc(func(name string) bool {
// 		return strings.ToLower(name) == key
// 	})
// 	if field.IsValid() {
// 		switch p := ptr.(type) {
// 		case *string:
// 			if *p != "" {
// 				logSet(field.Type().Name(), *p, "CLI -"+key)
// 				setField(field, *p)
// 			}
// 		case *int:
// 			if *p != 0 {
// 				logSet(field.Type().Name(), strconv.Itoa(*p), "CLI -"+key)
// 				setField(field, strconv.Itoa(*p))
// 			}
// 		}
// 	}
// }
