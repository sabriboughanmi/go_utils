package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/derekstavis/go-qs"
	"reflect"
	"strings"
)

type StructMapperKey string

const (
	JsonMapper = "json"
)

var (
	WrongTypePassed = errors.New("")
)

/*
InterfaceToType Converts an interface to any type.

Note: Make sure the destination type is passed as a reference (&) and is compatible with the source interface.
This Func can be Heavy if the Structure is too big
*/
func InterfaceAs(inter interface{}, myType interface{}) error {
	bytes, err := json.Marshal(&inter)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, myType)
}

///unsafeGetKindBitSize returns the bitSize of primitive numeric types.
//Note: Only numeric Types should be used (int, uint, int8, uint8, uint16,int16, uint32,int32, uint64,int64, Float32, Float64)
//0, 8, 16, 32, and 64 correspond to int, int8, int16, int32, and int64
func unsafeGetKindBitSize(kind reflect.Kind) int {
	switch kind {
	case reflect.Int:
		return 0
	case reflect.Int8:
		return 8
	case reflect.Int16:
		return 16
	case reflect.Int32:
		return 32
	case reflect.Int64:
		return 64
	case reflect.Uint:
		return 0
	case reflect.Uint8:
		return 8
	case reflect.Uint16:
		return 16
	case reflect.Uint32:
		return 32
	case reflect.Uint64:
		return 64
	case reflect.Float32:
		return 32
	case reflect.Float64:
		return 64
	default:
		panic("Unknown reflect.Type")
		return 0

	}
}

// RequestUrlToStruct maps a get request into struct.
func RequestUrlToStruct(urlRequest string, out interface{}) error {
	mapQuery, err := qs.Unmarshal(urlRequest)

	if err != nil {
		return fmt.Errorf("%#+v\n", mapQuery)
	}
	
	return mapStringInterfaceToModel(mapQuery, out, JsonMapper)
}

/*
// RequestUrlToStruct maps a get request into struct.
func RequestUrlToStructbb(urlRequest string, mapperKey StructMapperKey, out interface{}) error {

	mapQuery, err := qs.Unmarshal(urlRequest)

	if err != nil {
		return fmt.Errorf("%#+v\n", mapQuery)
	}

	bytes, err := json.Marshal(mapQuery)
	if err != nil {
		return fmt.Errorf("Error converting query to json - %#+v\n", mapQuery)
	}

	dType := reflect.TypeOf(out)

	//return json.Unmarshal(bytes, out)

	if dType.Kind() != reflect.Struct {
		var errorMessage = fmt.Sprintf("out must be of 'Struct' type: input type is:  %v \n", dType.Kind())
		return utils.CreateError(WrongTypePassed, errorMessage)
	}

	dhVal := reflect.ValueOf(out)

	for i := 0; i < dType.Elem().NumField(); i++ {

		field := dType.Elem().Field(i)
		key := field.Tag.Get(string(mapperKey))

		kind := field.Type.Kind()

		// Get the value from query params with given key
		val := urlRequest.Query().Get(key)

		//  Get reference of field value provided to input `out`
		result := dhVal.Elem().Field(i)

		switch kind {

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			bitSize := unsafeGetKindBitSize(kind)
			intVal, err := strconv.ParseInt(val, 10, bitSize)
			if err != nil {
				return err
			}
			if result.CanSet() {
				result.SetInt(intVal)
			}
			break
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			bitSize := unsafeGetKindBitSize(kind)
			intVal, err := strconv.ParseUint(val, 10, bitSize)
			if err != nil {
				return err
			}
			if result.CanSet() {
				result.SetUint(intVal)
			}
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
			if result.CanSet() {
				result.SetFloat(intVal)
			}
			break

		case reflect.String:
			if result.CanSet() {
				result.SetString(val)
			}
			break

		case reflect.Struct:
			if result.CanSet() {
				defaultVal := reflect.New(field.Type)
				err := json.Unmarshal([]byte(val), defaultVal.Interface())
				if err != nil {
					return err
				}
				result.Set(defaultVal.Elem())
			}
			break
		}

	}
	return nil
}
*/
//ValidateEmail checks if an email address is in valid format
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email must be a non-empty string")
	}
	if parts := strings.Split(email, "@"); len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return fmt.Errorf("malformed email string: %q", email)
	}
	return nil
}
