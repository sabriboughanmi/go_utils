package arrays

//Difference returns the Difference between 2 arrays (A - B)
//EG: a :[1, 2, 3, 4, 5]  - b: [2, 3, 5, 7, 11] =
// result: [1 4].
func Difference_Uint16(a, b []uint16) (diff []uint16) {
	m := make(map[uint16]bool)

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return
}

//Difference returns the Difference between 2 arrays (A - B)
//EG: a :["1", "2", "3", "4", "5"]  - b: ["2", "3", "5", "7", "11"] =
// result: ["1", "4"].
func Difference_String(a, b []string) (diff []string) {
	m := make(map[string]bool)

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return
}
