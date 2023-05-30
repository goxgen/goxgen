package goxgen

import "reflect"

func Indirect(value interface{}) interface{} {
	if reflect.TypeOf(value).Kind() == reflect.Ptr {
		return reflect.ValueOf(value).Elem().Interface()
	}
	return value
}
