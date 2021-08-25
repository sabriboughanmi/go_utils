package arrays

//Union returns the sum of 2 arrays (A + B) without duplicates
//EG: a :["1", "2", "3", "4", "5"]  + b: ["2", "3", "5", "7", "11"] =
// result: ["1" ,"2" ,"3", "4", "5", "7" ,"11"].
func Union(a, b []string) []string {
	check := make(map[string]int16)
	d := append(a, b...)
	res := make([]string,0)
	for _, val := range d {
		check[val] = 0
	}
	for key, _ := range check {
		res = append(res,key)
	}
	return res
}
