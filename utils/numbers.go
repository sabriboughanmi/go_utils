package utils

import (
	"reflect"
)

//IsInteger returns true if the interface value is of type Integer
func IsInteger(n interface{}) bool {
	switch reflect.TypeOf(n).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	default:
		return false
	}
}

//IsFloatingPointNumber returns true if the interface value is of type Floating Point
func IsFloatingPointNumber(n interface{}) bool {
	switch reflect.TypeOf(n).Kind() {
	case reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

//ToInt64 casts any number to int64
func ToInt64(n interface{}) (int64, bool) {
	var val = reflect.ValueOf(n)
	switch reflect.TypeOf(n).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int(), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(val.Uint()), true
	case reflect.Float32, reflect.Float64:
		return int64(val.Float()), true
	default:
		return 0, false
	}
}

//ToFloat64 casts any number to float64
func ToFloat64(n interface{}) (float64, bool) {
	var val = reflect.ValueOf(n)
	switch reflect.TypeOf(n).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(val.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(val.Uint()), true
	case reflect.Float32, reflect.Float64:
		return val.Float(), true
	default:
		return 0, false
	}
}
