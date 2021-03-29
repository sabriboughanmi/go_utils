package arrays

//IndexOf returns the Index of an Element into the List
//Note: a -1 is returned if the Element doesn't Exists
func IndexOf_String(element string, a []string) int {
 	var index =-1
	for i, item := range a {
		if item == element {
			index = i
		}
	}
	return index
}


//IndexOf returns the Index of an Element into the List
//Note: a -1 is returned if the Element doesn't Exists
func IndexOf_Uint(element uint, a []uint) int {
	var index =-1
	for i, item := range a {
		if item == element {
			index = i
		}
	}
	return index
}

//IndexOf returns the Index of an Element into the List
//Note: a -1 is returned if the Element doesn't Exists
func IndexOf_Int(element int, a []int) int {
	var index =-1
	for i, item := range a {
		if item == element {
			index = i
		}
	}
	return index
}

//IndexOf returns the Index of an Element into the List
//Note: a -1 is returned if the Element doesn't Exists
func IndexOf_Int64(element int64, a []int64) int {
	var index =-1
	for i, item := range a {
		if item == element {
			index = i
		}
	}
	return index
}

//IndexOf returns the Index of an Element into the List
//Note: a -1 is returned if the Element doesn't Exists
func IndexOf_Int32(element int32, a []int32) int {
	var index =-1
	for i, item := range a {
		if item == element {
			index = i
		}
	}
	return index
}