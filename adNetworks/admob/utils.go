package admob

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
)

var (
	 WrongTypePassed = errors.New("wrong Type Passed.")
)

// MapIt maps get request into struct
func MapIt(req *http.Request, out interface{}) error {
	dType := reflect.TypeOf(out)

	if dType.Kind() != reflect.Struct {
		fmt.Printf("out must be of 'Struct' type: input type is:  %v \n", dType.Kind())
		return WrongTypePassed
	}


	dhVal := reflect.ValueOf(out)

	for i := 0; i < dType.Elem().NumField(); i++ {

		field := dType.Elem().Field(i)
		key := field.Tag.Get("mapper")

		kind := field.Type.Kind()

		// Get the value from query params with given key
		val := req.URL.Query().Get(key)

		//  Get reference of field value provided to input `out`
		result := dhVal.Elem().Field(i)

		// we only check for string for now so,
		if kind == reflect.String {
			// set the value to string field
			// for other kinds we need to use different functions like; SetInt, Set etc
			result.SetString(val)
		} else {
			return errors.New("only supports string")
		}

	}
	return nil
}
