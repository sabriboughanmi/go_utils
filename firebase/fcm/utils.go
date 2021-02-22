package fcm

import "fmt"


//union returns the sum of 2 arrays (A + B) without duplicates
//EG: a :[1, 2, 3, 4, 5]  + b: [2, 3, 5, 7, 11] =
// result: [1 ,2 ,3, 4, 5, 7 ,11].
func union(a []Topic, b ...Topic) []Topic {
	for i := 0; i < len(b); i++ {
		found := false
		for j := 0; j < len(a); j++ {
			if b[i] == a[j] {
				found = true
				break
			}
		}
		if !found {
			a = append(a, b[i])
		}
	}
	return a
}

func (cf conditionFragment) toString() string {
	if len(cf.topics) == 1 {
		return fmt.Sprintf("%s'%s' in topics", cf.conditionOperator, cf.topics[0])
	}

	condition := string(cf.conditionOperator) + " "
	for i := 0; i < len(cf.topics)-1; i++ {
		condition += fmt.Sprintf("'%s' in topics %s", cf.topics[i], cf.operandsOperator)
	}
	condition += fmt.Sprintf("'%s' in topics", cf.topics[len(cf.topics)-1])
	return condition
}
