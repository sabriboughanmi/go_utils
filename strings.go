package go_utils

import (
	"strings"
)

//ReplaceKeys formats a string using key value pairs
func ReplaceKeys(s string, keyValues map[string]string) string {
	var newString = s
	for key, value := range keyValues {
		newString = strings.Replace(newString, "{"+key+"}", value, 1)
	}
	return newString
}
