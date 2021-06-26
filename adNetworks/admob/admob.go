package admob

import (
	"github.com/hiyali/go-lib-ssv/admob"
	"net/url"
)

//VerifyURL Verifies Admob CallBack URL Verification.
func VerifyURL(callBackUrl *url.URL) error {
	return admob.Verify(callBackUrl)
}

func GetParameters(requestBody string, out interface{}) error {

	return nil
}
