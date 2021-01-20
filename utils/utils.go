package utils

import "encoding/json"

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

