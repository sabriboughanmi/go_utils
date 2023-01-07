package querybatchupdate

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	firestoreV1 "google.golang.org/api/firestore/v1"
	"reflect"
	"time"
)

type firestoreFields map[string]firestoreV1.Value

//Data converts the firestore Values(map[string]firestore.Value) to map[string]interface{}
func (value firestoreFields) Data() (map[string]interface{}, error) {
	m := make(map[string]interface{})
	for k, v := range value {
		val, err := convertValue(v)
		if err != nil {
			return nil, err
		}
		m[k] = val
	}
	return m, nil
}

//DataTo converts Firestore Values to a model.
func (value firestoreFields) DataTo(typePtr interface{}) error {
	//Check if user passed a non Ptr Type
	if reflect.ValueOf(typePtr).Kind() != reflect.Ptr {
		return fmt.Errorf("v must be a pointer")
	}

	//Convert the firestoreFields to a map[string]interface{}
	m := make(map[string]interface{})
	for k, v := range value {
		val, err := convertValue(v)
		if err != nil {
			return err
		}
		m[k] = val
	}

	//Marshal the map in order to convert it to Model
	bytes, err := json.Marshal(m)
	if err != nil {
		return err
	}
	//Convert it to the passed model
	return json.Unmarshal(bytes, typePtr)
}

//valueAs set's a firestore.Value to an a Runtime Type Ptr
func valueAs(firestoreValue firestoreV1.Value, typePtr interface{}) error {

	//Get the typePtr reflect.Value as Ptr
	typePtrValuePtr := reflect.ValueOf(typePtr)

	//Check if typePtr is a Ptr
	if typePtrValuePtr.Kind() != reflect.Ptr {
		return fmt.Errorf("typePtr must be a pointer")
	}

	//Set the Value to be nil if nilled
	if firestoreValue.NullValue != "" {
		typePtrValuePtr.Set(reflect.ValueOf(nil))
		return nil
	}

	//Get the typePtr reflect.Value without the Ptr
	typePtrValue := typePtrValuePtr.Elem()
	typePtrValueKind := typePtrValue.Kind()

	switch {

	case typePtrValueKind == reflect.Bool:
		typePtrValue.SetBool(firestoreValue.BooleanValue)
		return nil

	case typePtrValueKind == reflect.String:
		typePtrValue.SetString(firestoreValue.StringValue)
		return nil

	case typePtrValueKind == reflect.Float32, typePtrValueKind == reflect.Float64:
		typePtrValue.SetFloat(firestoreValue.DoubleValue)
		return nil

	case typePtrValueKind == reflect.Int, typePtrValueKind == reflect.Int8, typePtrValueKind == reflect.Int16, typePtrValueKind == reflect.Int32, typePtrValueKind == reflect.Int64,
		typePtrValueKind == reflect.Uint, typePtrValueKind == reflect.Uint8, typePtrValueKind == reflect.Uint16, typePtrValueKind == reflect.Uint32, typePtrValueKind == reflect.Uint64:
		typePtrValue.SetInt(firestoreValue.IntegerValue)
		return nil

	case typePtrValue.Type() == reflect.TypeOf(firestoreV1.LatLng{}):
		typePtrValue.Set(reflect.ValueOf(*firestoreValue.GeoPointValue))
		return nil

	case typePtrValue.Type() == reflect.TypeOf(time.Time{}):
		t, err := time.Parse(time.RFC3339, firestoreValue.TimestampValue)
		if err != nil {
			return err
		}
		typePtrValue.Set(reflect.ValueOf(t))
		return nil

	case typePtrValue.Type() == reflect.TypeOf([]byte{}):
		//Get the Bytes data
		bytesValue, err := base64.StdEncoding.DecodeString(firestoreValue.BytesValue)
		if err != nil {
			return err
		}
		//Set the final value to the typePtrValue
		typePtrValue.Set(reflect.ValueOf(bytesValue))
		return nil

	case typePtrValueKind == reflect.Struct, typePtrValueKind == reflect.Map:
		return firestoreFields(firestoreValue.MapValue.Fields).DataTo(typePtr)

	case typePtrValueKind == reflect.Array, typePtrValueKind == reflect.Slice:
		//Create and Fill the Slice
		slice := make([]interface{}, len(firestoreValue.ArrayValue.Values))
		for i, v := range firestoreValue.ArrayValue.Values {
			val, err := convertValue(*v)
			if err != nil {
				return err
			}
			slice[i] = val
		}
		//Marshal the Data
		bytes, err := json.Marshal(slice)
		if err != nil {
			return err
		}
		//Unmarshal it
		return json.Unmarshal(bytes, typePtr)

	case typePtrValueKind == reflect.Ptr:
		return valueAs(firestoreValue, typePtrValue.Elem())

	default:
		return fmt.Errorf("default called for FieldName '%s' of Kind '%s' \n", typePtrValue.Type().Name(), typePtrValue.Kind().String())
	}

}

// GetFieldAs converts the firestore Value to the specified type and stores the result in typePtr.
func (value firestoreFields) GetFieldAs(key string, typePtr interface{}) error {
	firestoreValue, ok := value[key]
	if !ok {
		return fmt.Errorf("key %q not found in firestoreFields", key)
	}

	return valueAs(firestoreValue, typePtr)
}

// TryGetFieldAs checks either a Key exists and converts it's firestore Value to the specified type and stores the result in typePtr.
func (value firestoreFields) TryGetFieldAs(key string, typePtr interface{}) (bool, error) {
	firestoreValue, ok := value[key]
	if !ok {
		return ok, nil
	}

	return true, valueAs(firestoreValue, typePtr)
}

//convertValue converts a firestore.Value to it's corresponding value as interface
func convertValue(value firestoreV1.Value) (interface{}, error) {
	switch {
	case value.NullValue != "":
		return nil, nil
	case value.BooleanValue:
		return value.BooleanValue, nil
	case value.IntegerValue != 0:
		return value.IntegerValue, nil
	case value.DoubleValue != 0:
		return value.DoubleValue, nil
	case value.TimestampValue != "":
		t, err := time.Parse(time.RFC3339, value.TimestampValue)
		if err != nil {
			return nil, nil
		}
		return t, nil
	case value.StringValue != "":
		return value.StringValue, nil
	case value.BytesValue != "":
		b, err := base64.StdEncoding.DecodeString(value.BytesValue)
		if err != nil {
			return err, nil
		}
		return b, nil
	case value.ReferenceValue != "":
		return value.ReferenceValue, nil
	case value.GeoPointValue != nil:
		return value.GeoPointValue, nil
	case value.ArrayValue != nil:
		slice := make([]interface{}, 0, len(value.ArrayValue.Values))
		for _, v := range value.ArrayValue.Values {
			val, err := convertValue(*v)
			if err != nil {
				return nil, err
			}
			slice = append(slice, val)
		}
		return slice, nil
	case value.MapValue != nil:
		m := make(map[string]interface{})
		for k, v := range value.MapValue.Fields {
			val, err := convertValue(v)
			if err != nil {
				return nil, err
			}
			m[k] = val
		}
		return m, nil
	default:
		return nil, fmt.Errorf("unknown value type: %v", value)

	}
}
