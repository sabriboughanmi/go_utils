package utils

import (
	"fmt"
	"testing"
)

func TestRequestUrlToStruct(t *testing.T) {
	var sampleUrl = "employee[name]=sonoo&employee[salary]=56000&employee[married]=true&employee[employee][name]=sonoo&employee[employee][salary]=100000&employee[employee][married]=true&Sexe=male"

	type Employee struct {
		Name     string    `json:"name"`
		Salary   int64     `json:"salary"`
		Married  bool      `json:"married"`
		Employee *Employee `json:"employee,omitempty"`
	}
	type Person struct {
		Employee Employee `json:"employee"`
		Sexe     string   `json:"Sexe"`
	}

	var main Person

	if err := RequestUrlToStruct(sampleUrl, &main); err != nil {
		t.Errorf("deserialization Error %v", err)
	} else {
		fmt.Println("func RequestUrlToStruct OK!")
	}

	//fmt.Println(string(UnsafeAnythingToJSON(main)))
}
