package utils

import (
	"encoding/json"
	"reflect"
	"strconv"
)

//StringToAnyPrimitiveType converts a string to any Primitive Type.
// Supported Types : int... uint .. float.. string
func StringToAnyPrimitiveType(val string, out interface{}) error {
	dType := reflect.TypeOf(out)
	kind := dType.Kind()
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		bitSize := unsafeGetKindBitSize(kind)
		intVal, err := strconv.ParseInt(val, 10, bitSize)
		if err != nil {
			return err
		}
		out = intVal
		break
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		bitSize := unsafeGetKindBitSize(kind)
		intVal, err := strconv.ParseUint(val, 10, bitSize)
		if err != nil {
			return err
		}
		out = intVal
		break

	case reflect.Float64, reflect.Float32:
		bitSize := 32
		if kind == reflect.Float64 {
			bitSize = 64
		}
		intVal, err := strconv.ParseFloat(val, bitSize)
		if err != nil {
			return err
		}
		out = intVal
		break

	case reflect.String:
		out = val
		break
	}

	return nil
}

func mapStringInterfaceToModel(anything map[string]interface{}, outPtr interface{}, mapperKey StructMapperKey) error {

	dType := reflect.TypeOf(outPtr)
	dhVal := reflect.ValueOf(outPtr)

	for i := 0; i < dType.Elem().NumField(); i++ {

		field := dType.Elem().Field(i)
		key := field.Tag.Get(string(mapperKey))

		kind := field.Type.Kind()

		// Get the value from query params with given key
		val := anything[key]

		//  Get reference of field value provided to input `out`
		result := dhVal.Elem().Field(i)

		switch kind {

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intVal, err := strconv.ParseInt(val.(string), 10, 64)
			if err != nil {
				return err
			}
			if result.CanSet() {
				result.SetInt(intVal)
			}
			break
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			intVal, err := strconv.ParseUint(val.(string), 10, 64)
			if err != nil {
				return err
			}
			if result.CanSet() {
				result.SetUint(intVal)
			}
			break

		case reflect.Float64, reflect.Float32:
			floatVal, err := strconv.ParseFloat(val.(string), 64)
			if err != nil {
				return err
			}
			if result.CanSet() {
				result.SetFloat(floatVal)
			}
			break

		case reflect.String:
			if result.CanSet() {
				result.SetString(val.(string))
			}
			break
		case reflect.Struct:
			defaultStructVal := reflect.New(field.Type)
			bytes, err := json.Marshal(val)
			if err != nil {
				return err
			}

			var newEntry map[string]interface{}
			if err = json.Unmarshal(bytes, &newEntry); err != nil {
				return err
			}

			if err = mapStringInterfaceToModel(newEntry, defaultStructVal.Interface(), JsonMapper); err != nil {
				return err
			}

			result.Set(defaultStructVal.Elem())
			break

			//TODO: add Support for Arrays
		}

	}
	return nil
}
