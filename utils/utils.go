package utils

import (
	"encoding/json"
	"fmt"
	"github.com/derekstavis/go-qs"
	"strings"
)

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

// RequestUrlToStruct maps a get request into struct.
func RequestUrlToStruct(urlRequest string, out interface{}) error {
	mapQuery, err := qs.Unmarshal(urlRequest)

	if err != nil {
		return fmt.Errorf("%#+v\n", mapQuery)
	}

	bytes, err := json.Marshal(mapQuery)
	if err != nil {
		return fmt.Errorf("Error converting query to json - %#+v\n", mapQuery)
	}

	return json.Unmarshal(bytes, out)
}

//ValidateEmail checks if an email address is in valid format
func ValidateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email must be a non-empty string")
	}
	if parts := strings.Split(email, "@"); len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return fmt.Errorf("malformed email string: %q", email)
	}
	return nil
}
