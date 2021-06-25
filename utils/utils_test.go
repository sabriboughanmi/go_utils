package utils

import (
	"fmt"
	"net/url"
	"testing"
)

func TestRequestUrlToStruct(t *testing.T) {
	var sampleUrl = "data[name]=Dipesh Dulal&data[id]=12&data[json][name]=dipesh&data[json][id]=64&data[json][fl]=5.05&message=succesfully bind to getRequestInput"
	sampleRequestURL, _ := url.ParseRequestURI(sampleUrl)

	type JSON struct {
		Name string `json:"name"`
		ID   int    `json:"id"`
	}
	type Data struct {
		Name string `json:"name"`
		ID   uint8  `json:"id"`
		JSON JSON   `json:"json"`
	}

	type Main struct {
		Data    Data   `json:"data"`
		Message string `json:"message"`
	}

	var main Main

	if err := RequestUrlToStruct(sampleRequestURL, JsonMapper, main); err != nil {
		t.Errorf("deserialization Error %v", err)
	} else {
		fmt.Println("func RequestUrlToStruct OK!")
	}
}
