package tmpl

import "reflect"

// Indirect returns the value of a pointer
func Indirect(value interface{}) interface{} {
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		return reflect.ValueOf(value).Elem().Interface()
	}
	return value
}
