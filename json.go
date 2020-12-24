package utils

import "encoding/json"

//AnythingToJSON : encode any Type to Json string
func AnythingToJSON(anything interface{}) (string, error) {

	resultBytes, err := json.Marshal(anything)
	if err != nil {
		return "", err
	}
	return string(resultBytes), nil
}

//UnsafeAnythingToJSON : Incode any Type to Json string, can only be used for 100% safe conversion
func UnsafeAnythingToJSON(anything interface{}) []byte {

	resultBytes, err := json.Marshal(anything)
	if err != nil {
		panic(err)
	}
	return resultBytes
}
