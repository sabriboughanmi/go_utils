package userprivileges


//Difference returns the Difference between 2 arrays (A - B)
//EG: a :[1, 2, 3, 4, 5]  - b: [2, 3, 5, 7, 11] =
// result: [1 4].
func difference(a, b []Privilege) (diff []Privilege) {
	m := make(map[Privilege]bool)

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



//Intersection returns the intersection between 2 arrays
//EG: a :[1, 2, 3, 4, 5], b: [2, 3, 5, 7, 11]
// result: [2 3 5].
func Intersection(a, b []Privilege) (c []Privilege) {
	m := make(map[Privilege]bool)

	for _, item := range a {
		m[item] = true
	}

	for _, item := range b {
		if _, ok := m[item]; ok {
			c = append(c, item)
		}
	}
	return
}