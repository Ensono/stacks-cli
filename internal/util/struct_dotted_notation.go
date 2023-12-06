package util

import (
	"fmt"
	"reflect"
	"strings"
)

// GetValueByDottedPath returns the value of a field in a struct by using a string
// This is useful when checking the values of a struct for a specific value, or if it is null
// The string must be names of the attributes in the struct and not the ones in the mapstructure
func GetValueByDottedPath(data interface{}, path string) (interface{}, error) {

	// get the name of the fields that are being sought
	fields := strings.Split(path, ".")
	val := reflect.ValueOf(data)

	// iterate around the fields
	for _, field := range fields {
		if val.Kind() == reflect.Ptr {
			val = reflect.ValueOf(val).Elem()
		}

		if val.Kind() == reflect.Struct {
			// val = val.FieldByNameFunc(func(n string) bool { return strings.ToLower(n) == field })
			val = caseInsensitiveFieldByName(val, field)
		} else {
			return nil, fmt.Errorf("unable to find field: %s", field)
		}
	}

	if val.IsValid() {
		return val.Interface(), nil
	}

	return nil, fmt.Errorf("field not found: %s", path)
}

func caseInsensitiveFieldByName(v reflect.Value, name string) reflect.Value {
	name = strings.ToLower(name)
	return v.FieldByNameFunc(func(n string) bool { return strings.ToLower(n) == name })
}
