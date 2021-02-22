package fcm

import (
	"context"
	"firebase.google.com/go/messaging"
)

type Operator string

const (
	AND Operator = "&&"
	OR  Operator = "||"
	none  Operator = ""
)



type conditionFragment struct {
	//Is the Operator applied between Topics
	operandsOperator Operator
	topics           []Topic
	//Is the Operator applied between Conditions
	conditionOperator Operator

}

//Condition struct to handle Target Topics
type Condition struct {

	includedTopics []Topic //Used Topics
	conditionFragments     []conditionFragment

	hasChanges bool
	condition string
}

type NotificationManager struct {
	client *messaging.Client
	ctx    context.Context
}
