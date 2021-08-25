package fcm

import "fmt"

func TopicsCountExceeded(topics ...Topic) {
	if len(topics) > maxTopicsPerCondition {
		panic(fmt.Sprintf("Topics Count Exceeded: Maximum Allowed Topics per ConditionBuilder is %d ! \n", maxTopicsPerCondition))
	}
}

//"'TopicA' in topics && ('TopicB' in topics || 'TopicC' in topics)"

//CreateConditionBuilder Creates a new ConditionBuilder.
func CreateConditionBuilder(operand Topic) ConditionBuilder {
	//Create ConditionBuilder
	condition := ConditionBuilder{
		includedTopics:     make([]Topic, 0),
		conditionFragments: make([]conditionFragment, 0),
		hasChanges:         true,
		condition:          "",
	}
	condition.addCondition(none, none, operand)

	return condition
}

//CreateConditionBuilderMultiTopics Creates a new ConditionBuilder using multiple Topics and an operator.
func CreateConditionBuilderMultiTopics(operator Operator, operands ...Topic) ConditionBuilder {
	//Create ConditionBuilder
	condition := ConditionBuilder{
		includedTopics:     make([]Topic, 0),
		conditionFragments: make([]conditionFragment, 0),
		hasChanges:         true,
		condition:          "",
	}
	condition.addCondition(none, operator, operands...)
	return condition
}

//AndInTopic extends the ConditionBuilder with a new condition using the AND operator.
func (c *ConditionBuilder) AndInTopic(operand Topic) {
	//Add ConditionBuilder Fragment
	c.addCondition(AND, none, operand)
}


//AndInTopics extends the ConditionBuilder with a multiple new Topics using the AND operator.
func (c *ConditionBuilder) AndInTopics(operator Operator, operands ...Topic) {
	//Add ConditionBuilder Fragment
	c.addCondition(AND, operator, operands...)
}

//OrInTopic extends the ConditionBuilder with a new condition using the OR operator.
func (c *ConditionBuilder) OrInTopic(operand Topic) {
	//Add ConditionBuilder Fragment
	c.addCondition(OR, none, operand)
}

//OrInTopics extends the ConditionBuilder with a multiple new Topics using the OR operator.
func (c *ConditionBuilder) OrInTopics(operator Operator, operands ...Topic) {
	//Add ConditionBuilder Fragment
	c.addCondition(OR, operator, operands...)
}

//Expends the Existing ConditionBuilder
func (c *ConditionBuilder) addCondition(conditionOperator, operandsOperator Operator, operands ...Topic) {

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

//BuildCondition Returns the Constructed string Condition.
func (c *ConditionBuilder) BuildCondition() string {
	if !c.hasChanges {
		return c.condition
	}
	c.hasChanges = false

	c.condition = c.conditionFragments[0].toString()
	for i := 1; i < len(c.conditionFragments); i++ {
		c.condition += " " + c.conditionFragments[i].toString()
	}
	return c.condition
}
