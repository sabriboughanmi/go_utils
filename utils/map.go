package utils

import (
	"fmt"
	"strings"
)

//GetValueFromSubMap Returns the value of SubMap at path.
// the keyPath need to be separated by . any other symbol is rejected
func GetValueFromSubMap(dict map[string]interface{}, keyPath string) (interface{}, error) {
	keys := strings.Split(keyPath, ".")

	var currentMap = dict
	var done = false

	for i := 0; i < len(keys)-1; i++ {
		key := keys[i]
		currentMap, done = currentMap[key].(map[string]interface{})
		if !done {
			return nil, fmt.Errorf("incorrect keyPath")
		}
	}
	return currentMap[keys[len(keys)-1]], nil
}

