package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

func mapStringInterfaceToMappedModel(anything map[string]interface{}, usedType interface{}, mapperKey StructMapperKey, outPutMap *map[string]interface{}) error {

	dType := reflect.TypeOf(usedType)

	for i := 0; i < dType.Elem().NumField(); i++ {

		field := dType.Elem().Field(i)
		key := field.Tag.Get(string(mapperKey))

		kind := field.Type.Kind()

		// Get the value from query params with given key
		val := anything[key]

		switch kind {

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intVal, err := strconv.ParseInt(val.(string), 10, 64)
			if err != nil {
				return err
			}
			(*outPutMap)[key] = intVal
			break

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			intVal, err := strconv.ParseUint(val.(string), 10, 64)
			if err != nil {
				return err
			}
			(*outPutMap)[key] = intVal
			break

		case reflect.Float64, reflect.Float32:
			floatVal, err := strconv.ParseFloat(val.(string), 64)
			if err != nil {
				return err
			}
			(*outPutMap)[key] = floatVal
			break

		case reflect.String:
			(*outPutMap)[key] = val.(string)
			break

		case reflect.Interface:
			(*outPutMap)[key] = val
			break

		case reflect.Bool:
			valBool, err := strconv.ParseBool(val.(string))
			if err != nil {
				return err
			}
			(*outPutMap)[key] = valBool
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

			var outMap = make(map[string]interface{}, len(newEntry))
			if err = mapStringInterfaceToMappedModel(newEntry, defaultStructVal.Interface(), mapperKey, &outMap); err != nil {
				return err
			}
			(*outPutMap)[key] = outMap
			break

		case reflect.Slice:
			bytes, err := json.Marshal(val)
			if err != nil {
				return err
			}

			var newEntry = make([]interface{}, 0)
			if err = json.Unmarshal(bytes, &newEntry); err != nil {

				//Slice can be Corrupted by the url query (it's transformed to map[string]interface{}, so in this case we simply collect Values)
				var tempEntry = make(map[string]interface{})
				if err = json.Unmarshal(bytes, &tempEntry); err != nil {
					return err
				}

				newEntry = make([]interface{}, len(tempEntry))
				//we are iterating the map in this way to keep original array order
				for index := 0; index < len(tempEntry); index++ {
					sIndex := strconv.Itoa(index)
					newEntry[index] = tempEntry[sIndex]
				}

			}

			var newOutput = make([]interface{}, len(newEntry))
			if err = arrayInterfaceToModel(newEntry, &newOutput, mapperKey, field.Type); err != nil {
				return err
			}
			(*outPutMap)[key] = &newOutput
			break

		default:
			fmt.Printf("Unsupported Kind : %s \n", kind.String())
			break
		}

	}
	return nil
}

func arrayInterfaceToModel(anything []interface{}, outPtrArray *[]interface{}, mapperKey StructMapperKey, arrayType reflect.Type) error {
	elemType := arrayType.Elem()
	elemKind := elemType.Kind()

	if len(anything) == 0 {
		outPtrArray = &anything
		return nil
	}

	for i, val := range anything {
		switch elemKind {

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intVal, err := strconv.ParseInt(val.(string), 10, 64)
			if err != nil {
				return err
			}
			(*outPtrArray)[i] = intVal
			break

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			intVal, err := strconv.ParseUint(val.(string), 10, 64)
			if err != nil {
				return err
			}
			(*outPtrArray)[i] = intVal
			break

		case reflect.Float64, reflect.Float32:
			floatVal, err := strconv.ParseFloat(val.(string), 64)
			if err != nil {
				return err
			}
			(*outPtrArray)[i] = floatVal
			break

		case reflect.String:
			(*outPtrArray)[i] = val.(string)
			break

		case reflect.Interface:
			(*outPtrArray)[i] = val
			break

		case reflect.Bool:
			valBool, err := strconv.ParseBool(val.(string))
			if err != nil {
				return err
			}
			(*outPtrArray)[i] = valBool
			break

		case reflect.Struct:
			defaultStructVal := reflect.New(elemType)
			bytes, err := json.Marshal(val)
			if err != nil {
				return err
			}

			var newEntry map[string]interface{}
			if err = json.Unmarshal(bytes, &newEntry); err != nil {
				return err
			}

			var outMap = make(map[string]interface{}, len(newEntry))
			if err = mapStringInterfaceToMappedModel(newEntry, defaultStructVal.Interface(), mapperKey, &outMap); err != nil {
				return err
			}
			(*outPtrArray)[i] = outMap
			break

		case reflect.Slice:
			bytes, err := json.Marshal(val)
			if err != nil {
				return err
			}
			var newEntry = make([]interface{}, 0)
			if err = json.Unmarshal(bytes, &newEntry); err != nil {

				//Slice can be Corrupted by the url query (it's transformed to map[string]interface{}, so in this case we simply collect Values)
				var tempEntry = make(map[string]interface{})
				if err = json.Unmarshal(bytes, &tempEntry); err != nil {
					return err
				}

				newEntry = make([]interface{}, len(tempEntry))
				//we are iterating the map in this way to keep original array order
				for index := 0; index < len(tempEntry); index++ {
					sIndex := strconv.Itoa(index)
					newEntry[index] = tempEntry[sIndex]
				}
			}

			var newOutPut = make([]interface{}, len(newEntry))
			if err = arrayInterfaceToModel(newEntry, &newOutPut, mapperKey, elemType); err != nil {
				return err
			}
			(*outPtrArray)[i] = newOutPut
			break

		default:
			fmt.Printf("Unsupported Kind b: %s \n", elemKind.String())
			break
		}

	}

	return nil
}
