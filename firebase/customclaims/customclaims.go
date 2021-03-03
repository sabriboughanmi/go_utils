package customclaims

import (
	"context"
	"firebase.google.com/go/auth"
	"github.com/sabriboughanmi/go_utils/utils"
)


//SetCustomUserClaims
func (cc *CustomClaims) SetCustomUserClaims(authClient *auth.Client,ctx context.Context, uid string) error {
	return authClient.SetCustomUserClaims(ctx, uid, *cc.ToMap())
}

//ToMap  Converts the CustomClaims Model to map[string]interface{}
func (cc *CustomClaims) ToMap() *map[string]interface{} {
	var customClaims = make(map[string]interface{})
	for key, value := range *cc {
		customClaims[string(key)] = value
	}
	return &customClaims
}

//ToCustomClaims Converts the Auth CustomClaims (map[string]interface{}) to CustomClaims Model
func ToCustomClaims(customClaimsMap map[string]interface{}) *CustomClaims {
	var customClaims = make(CustomClaims)
	for key, value := range customClaimsMap {
		customClaims[CustomClaimsKey(key)] = value
	}
	return &customClaims
}

//GetCustomClaimsAs returns a CustomClaim as type using a key
func (cc *CustomClaims) GetCustomClaimsAs(customClaim CustomClaimsKey, typeRef interface{}) error {
	var value interface{}
	var ok bool
	if value, ok = (*cc)[customClaim]; !ok {
		return NotFound
	}

	if err := utils.InterfaceAs(value, &typeRef); err != nil {
		return err
	}

	return nil
}

//GetCustomClaim Returns a CustomClaim with key
func (cc *CustomClaims) GetCustomClaim(customClaim CustomClaimsKey) *interface{} {
	if value, ok := (*cc)[customClaim]; ok {
		return &value
	}
	return nil
}

