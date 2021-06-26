package admob

import (
	"encoding/json"
	"github.com/hiyali/go-lib-ssv/admob"
	"github.com/sabriboughanmi/go_utils/utils"
	"net/url"
)

//VerifyURL Verifies Admob CallBack URL Verification.
func VerifyURL(callBackUrl *url.URL) error {
	return admob.Verify(callBackUrl)
}

//GetParameters Returns the SSVCallback model sent by Admob
func GetParameters(requestBody string) (SSVCallback, error) {
	var admobSSVCallback SSVCallback
	err := utils.RequestUrlToStruct(requestBody, &admobSSVCallback, utils.JsonMapper)
	return admobSSVCallback, err
}

//CustomDataAs converts the custom data to a specific type
func (callback *SSVCallback) CustomDataAs(out interface{}) error {
	return json.Unmarshal([]byte(callback.CustomData), out)
}
