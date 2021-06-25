package admob

import (
	"encoding/json"
	"fmt"
	"github.com/hiyali/go-lib-ssv/admob"
	"io"
	"net/url"
)

//VerifyURL Verifies Admob CallBack URL Verification.
func VerifyURL(callBackUrl *url.URL) error {
	return admob.Verify(callBackUrl)
}

func GetParameters(requestBody io.ReadCloser) interface{} {
	var user interface{}
	if err := json.NewDecoder(requestBody).Decode(&user); err != nil {
		fmt.Println("error decoding api payload")
		return err
	}
	return user
}
