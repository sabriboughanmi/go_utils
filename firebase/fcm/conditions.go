package fcm

import "fmt"

func TopicsCountExceeded(topics ...Topic) {
	if len(topics) > maxTopicsPerCondition {
		panic(fmt.Sprintf("Topics Count Exceeded: Maximum Allowed Topics per Condition is %d ! \n", maxTopicsPerCondition))
	}
}

//"'TopicA' in topics && ('TopicB' in topics || 'TopicC' in topics)"

//Returns an Empty Condition
func InTopic(operand Topic) Condition {
	//Create Condition
	condition := Condition{
		includedTopics:     make([]Topic, 0),
		conditionFragments: make([]conditionFragment, 0),
		hasChanges:         true,
		condition:          "",
	}
	condition.addCondition(none, none, operand)

	return condition
}

//Returns an Empty Condition
func InTopics(operator Operator, operands ...Topic) Condition {
	//Create Condition
	condition := Condition{
		includedTopics:     make([]Topic, 0),
		conditionFragments: make([]conditionFragment, 0),
		hasChanges:         true,
		condition:          "",
	}
	condition.addCondition(none, operator, operands...)
	return condition
}

//Adds a single Operand Condition
func (c *Condition) AndInTopic(operand Topic) {
	//Add Condition Fragment
	c.addCondition(AND, none, operand)
}

//Adds a Multi Operands Condition
func (c *Condition) AndInTopics(operator Operator, operands ...Topic) {
	//Add Condition Fragment
	c.addCondition(AND, operator, operands...)
}

//Adds a single Operand Condition
func (c *Condition) OrInTopic(operand Topic) {
	//Add Condition Fragment
	c.addCondition(OR, none, operand)
}

//Adds a Multi Operands Condition
func (c *Condition) OrInTopics(operator Operator, operands ...Topic) {
	//Add Condition Fragment
	c.addCondition(OR, operator, operands...)
}

//Expends the Existing Condition
func (c *Condition) addCondition(conditionOperator, operandsOperator Operator, operands ...Topic) {

	//Merge Topics
	c.includedTopics = union(c.includedTopics, operands...)

	//Check Topics Count
	TopicsCountExceeded(c.includedTopics...)

	//Mark as Changed
	c.hasChanges = true

	//Expend the condition
	c.conditionFragments = append(c.conditionFragments,
		conditionFragment{
			conditionOperator: conditionOperator,
			operandsOperator:  operandsOperator,
			topics:            operands,
		})
}


//Returns the Constructed Condition
func (c *Condition) GetCondition() string {
	if !c.hasChanges {
		fmt.Printf("Condition Prefetched!")
		return c.condition
	}
	c.hasChanges = false

	c.condition = c.conditionFragments[0].toString()
	for i := 1; i < len(c.conditionFragments); i++ {
		c.condition += " " + c.conditionFragments[i].toString()
	}
	return c.condition
}
