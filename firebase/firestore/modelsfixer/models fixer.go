package modelsfixer

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

//TODO: make it Dynamic, so people can create their own safety Tags and add their handlers functions.

//processField process models fields.
func processField(val interface{}, tags reflect.StructTag) (interface{}, bool) {
	valT := reflect.TypeOf(val)
	valV := reflect.ValueOf(val)

	firestoreTags := getFirestoreTags(tags)

	switch valT.Kind() {

	case reflect.Bool:
		//Check for firestore tags
		if firestoreTags != nil {
			//Get field Value
			var value = valV.Bool()

			//Check for Casts
			var tag = firestoreTags.ContainsAny(Tags_String, Tags_Int)
			if tag != "" {
				switch tag {
				case Tags_String:
					return strconv.FormatBool(value), true
				case Tags_Int:
					return func() int64 {
						if value {
							return 1
						}
						return 0
					}, true
				default:
					panic(fmt.Sprintf("unsupported CastType '%s' for type '%s'", tag, valT.Kind().String()))

				}
			}
		}
		return valV.Interface(), true
	case reflect.String:
		//Check for firestore tags
		if firestoreTags != nil {
			//Get field Value
			var value = valV.String()

			//Check Omitempty
			if firestoreTags.ContainsTag(Tags_Omitempty) && value == "" {
				return nil, false
			}

			//Check for Casts
			var tag = firestoreTags.ContainsAny(Tags_String, Tags_Float, Tags_Int)
			if tag != "" {
				switch tag {
				case Tags_String:
					return valV.String(), true

				case Tags_Float:
					val, err := strconv.ParseFloat(valV.String(), 10)
					if err != nil {
						panic(err)
					}
					return val, true

				case Tags_Int:
					val, err := strconv.ParseFloat(valV.String(), 10)
					if err != nil {
						panic(err)
					}
					return int64(val), true

				default:
					panic(fmt.Sprintf("unsupported CastType '%s' for type '%s'", tag, valT.Kind().String()))

				}
			}
		}

		return valV.Interface(), true

	case reflect.Float32, reflect.Float64,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:

		//Check for firestore tags
		if firestoreTags != nil {
			tag := firestoreTags.ContainsAny(Tags_Float, Tags_Int, Tags_String)
			if tag != "" {

				//Check for Casts
				switch tag {
				case Tags_Float:
					return Float64(valV.Interface()), true

				case Tags_Int:
					return Int64(valV.Interface()), true

				case Tags_String:
					return fmt.Sprintf("%v", valV.Interface()), true

				default:
					panic(fmt.Sprintf("unsupported CastType '%s' for type '%s'", tag, valT.Kind().String()))
				}
			}

		}

		return valV.Interface(), true

	case reflect.Struct:

		//Skip parsing if requested
		if firestoreTags.ContainsTag(Tags_SkipParsing) {
			return val, true
		}

		var safeType = make(map[string]interface{})

		//Iterate over fields
		for i := 0; i < valT.NumField(); i++ {

			var fieldName = valT.Field(i).Name

			//check if the field is private
			if skipField(fieldName) {
				continue
			}

			//Process Field
			var val, mustBeAdded = processField(valV.Field(i).Interface(), valT.Field(i).Tag)
			if !mustBeAdded {
				continue
			}

			fieldTags := getFirestoreTags(valT.Field(i).Tag)

			if len(fieldTags) > 0 {
				//get firestore field name
				fieldName = fieldTags[0]
			}

			//Set the safe data
			safeType[fieldName] = val

		}

		//if safeType length is 0 and this field is tagged with omitempty then return a null
		if firestoreTags.ContainsTag(Tags_Omitempty) && len(safeType) == 0 {
			return nil, false
		}

		return safeType, true

	case reflect.Array, reflect.Slice:
		//Check Omitempty
		if firestoreTags.ContainsTag(Tags_Omitempty) && valV.IsNil() {
			return nil, false
		}

		//Skip parsing if requested
		if firestoreTags.ContainsTag(Tags_SkipParsing) {
			return val, true
		}

		var finalSlice = make([]interface{}, 0)

		for i := 0; i < valV.Len(); i++ {
			item := valV.Index(i)
			elementVal, mustBeTaken := processField(item.Interface(), tags)
			if !mustBeTaken {
				continue
			}
			finalSlice = append(finalSlice, elementVal)
		}

		return finalSlice, true

	case reflect.Map:

		//Check Omitempty
		if firestoreTags.ContainsTag(Tags_Omitempty) && valV.IsNil() {
			return nil, false
		}

		//Skip parsing if requested
		if firestoreTags.ContainsTag(Tags_SkipParsing) {
			return val, true
		}

		var safeType = make(map[string]interface{})

		//Process map fields
		iter := valV.MapRange()
		for iter.Next() {
			k := iter.Key()
			var keyVal string
			if k.Kind() == reflect.String {
				keyVal = k.String()
			} else {

				val, mustBeTaken := processField(k.Interface(), "")
				if !mustBeTaken {
					continue
				}
				keybytes, err := json.Marshal(val)
				if err != nil {
					panic(err)
				}
				keyVal = string(keybytes)
			}

			//Process the Value
			val, mustBeTaken := processField(iter.Value().Interface(), tags)

			//Collect the keys and values
			if mustBeTaken {
				safeType[keyVal] = val
			}

		}

		//if safeType length is 0 and this field is tagged with omitempty then return a null
		if firestoreTags.ContainsTag(Tags_Omitempty) && len(safeType) == 0 {
			return nil, false
		}

		return safeType, true

	case reflect.Ptr:
		//check if null
		if valV.IsNil() {
			//Check Omitempty
			if firestoreTags.ContainsTag(Tags_Omitempty) {
				return nil, false
			}
			//Must be set as null
			return nil, true
		}

		//Process Field
		var v, mustBeAdded = processField(valV.Elem().Interface(), tags)
		if !mustBeAdded {
			return nil, false
		}

		return v, true

	default:
		fmt.Printf("default called for FieldName '%s' of Kind '%s' \n", valT.Name(), valT.Kind().String())

		return nil, false

	}
	return nil, false
}

//GetSafeVersion returns a safe version of your model so firestore CRUD don't throw errors.
func GetSafeVersion(model interface{}) interface{} {
	var safeModel, _ = processField(model, "")
	return safeModel
}
