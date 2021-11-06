package utils

/*
func TestRequestUrlToStruct(t *testing.T) {
	var sampleUrl = "int=88&float=32.793&bool=false&string=Marsiella&array_string[]=Jobi&array_string[]=Betta&array_string[]=Lilith&array_string[]=Donetta&array_string[]=Cristabel&array_objects[0][index]=0&array_objects[0][indexString]=5&array_objects[1][index]=1&array_objects[1][indexString]=6&array_objects[2][index]=2&array_objects[2][indexString]=7&subMap[uint]=71&subMap[string]=Sabri&subMap[7ala]=true&subMap[subMap_array_int8][]=0&subMap[subMap_array_int8][]=10&subMap[subMap_array_int8][]=50&subMap[subMap_array_int8][]=60&subMap[subMap_array_int8][]=127&array_2D[0][]=0&array_2D[0][]=1&array_2D[0][]=2&array_2D[1][]=3&array_2D[1][]=4&array_2D[1][]=5&array_2D[1][]=6&array_3D[0][0][]=0&array_3D[0][0][]=1&array_3D[0][0][]=2&array_3D[0][1][]=0&array_3D[0][1][]=1&array_3D[0][1][]=2&array_3D[0][2][]=0&array_3D[0][2][]=1&array_3D[0][2][]=2&array_3D[1][0][]=0&array_3D[1][0][]=1&array_3D[1][0][]=2&array_3D[1][1][]=0&array_3D[1][1][]=1&array_3D[1][1][]=2&array_3D[1][2][]=0&array_3D[1][2][]=1&array_3D[1][2][]=2&array_3D[2][0][]=0&array_3D[2][0][]=1&array_3D[2][0][]=2&array_3D[2][1][]=0&array_3D[2][1][]=1&array_3D[2][1][]=2&array_3D[2][2][]=0&array_3D[2][2][]=1&array_3D[2][2][]=2"

	type Object struct {
		Index       int    `json:"index"`
		IndexString string `json:"indexString"`
	}

	type SubMap struct {
		Uint            uint   `json:"uint"`
		String          string `json:"string"`
		The7Ala         bool   `json:"7ala"`
		SubMapArrayInt8 []int8 `json:"subMap_array_int8"`
	}

	type ENUM int
	type Main struct {
		Int          ENUM        `json:"int"`
		InvalidKey   string      `json:"InvalidKeyExample,omitempty"`
		Float        float64     `json:"float"`
		Bool         bool        `json:"bool"`
		String       string      `json:"string"`
		ArrayString  []string    `json:"array_string"`
		ArrayObjects []Object    `json:"array_objects"`
		SubMap       SubMap      `json:"subMap"`
		Array2D      [][]int64   `json:"array_2D"`
		Array3D      [][][]int64 `json:"array_3D"`
	}

	var main Main

	if err := RequestUrlToStruct(sampleUrl, &main, JsonMapper); err != nil {
		t.Errorf("Error - %v", err)
	} else {
		fmt.Println("func RequestUrlToStruct OK!")
	}

	//fmt.Println(string(UnsafeAnythingToJSON(main)))
}
*/
