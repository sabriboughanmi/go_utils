package modelsfixer

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

type Tags []string

//getFirestoreTags returns firestore tags if exists.
func getFirestoreTags(tags reflect.StructTag) Tags {
	tagsString := tags.Get(FireStoreTag)
	if tagsString == "" {
		return nil
	}

	return strings.Split(tagsString, ",")

}

func (tags Tags) ContainsTag(tag ETag) bool {
	if tags == nil {
		return false
	}
	for _, t := range tags {
		if t == string(tag) {
			return true
		}
	}
	return false
}

func (tags Tags) ContainsAny(tagsToCheck ...ETag) ETag {
	if tags == nil {
		return ""
	}
	for _, tagToCheck := range tagsToCheck {
		if tags.ContainsTag(tagToCheck) {
			return tagToCheck
		}
	}
	return ""
}

//requireSafetyConversion returns true if the struct field contains a supported Tag.
func requireSafetyConversion(tags reflect.StructTag) bool {

	firestoreTags := getFirestoreTags(tags)
	if firestoreTags == nil {
		return false
	}

	for _, supportedTag := range supportedTags {
		if firestoreTags.ContainsTag(ETag(supportedTag)) {
			return true
		}
	}
	return false
}

func IsLower(s string) bool {
	for _, r := range s {
		if !unicode.IsLower(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

//getFieldValue returns a Value from a type.FieldName
func getFieldValue(obj interface{}, fieldName string) interface{} {
	pointToStruct := reflect.ValueOf(obj) // addressable
	curStruct := pointToStruct.Elem()
	if curStruct.Kind() != reflect.Struct {
		panic("not struct")
	}
	curField := curStruct.FieldByName(fieldName) // type: reflect.Value
	if !curField.IsValid() {
		panic("not found:" + fieldName)
	}
	return curField.Interface()
}

//skipField defines if a struct.field can be skipped.
func skipField(s string) bool {
	var firstChar = s[0:1]
	if IsLower(firstChar) || (firstChar < "A" || firstChar > "Z") {
		return true
	}

	return false
}

//Int64 casts any number to int64
func Int64(n interface{}) int64 {
	var val = reflect.ValueOf(n)
	switch reflect.TypeOf(n).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(val.Uint())
	case reflect.Float32, reflect.Float64:
		return int64(val.Float())
	default:
		panic(fmt.Sprintf("connot convert '%v' of type '%v' to int64", n, reflect.TypeOf(n).Kind()))
	}
}

//Float64 casts any number to float64
func Float64(n interface{}) float64 {
	var val = reflect.ValueOf(n)
	switch reflect.TypeOf(n).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(val.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(val.Uint())
	case reflect.Float32, reflect.Float64:
		return val.Float()
	default:
		panic(fmt.Sprintf("connot convert '%v' of type '%v' to int64", n, reflect.TypeOf(n).Kind()))
	}

}
