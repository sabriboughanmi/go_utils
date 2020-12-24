package utils

import (
	"fmt"
	"strings"
)

//ReplaceKeys formats a string using key value pairs
func ReplaceKeys(s string, keyValues map[string]string) string {
	fmt.Println("modif")
	var newString = s
	for key, value := range keyValues {
		newString = strings.Replace(newString, "{"+key+"}", value, 1)
	}
	return newString
}
