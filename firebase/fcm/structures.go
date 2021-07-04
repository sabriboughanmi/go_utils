package fcm

import (
	"context"
	"firebase.google.com/go/messaging"
)

type Operator string

const (
	AND  Operator = "&&"
	OR   Operator = "||"
	none Operator = ""
)

type conditionFragment struct {
	//Is the Operator applied between Topics
	operandsOperator Operator
	topics           []Topic
	//Is the Operator applied between Conditions
	conditionOperator Operator
}

//ConditionBuilder handles the creation of SAFE Conditions to query devices under specific Topics.
type ConditionBuilder struct {
	includedTopics     []Topic //Used Topics
	conditionFragments []conditionFragment

	hasChanges bool
	condition  string
}

type NotificationManager struct {
	client *messaging.Client
	ctx    context.Context
}

//NotificationData Used to represents a notification all platforms config
type NotificationData struct {
	Notification  *messaging.Notification
	Data          map[string]string
	AndroidConfig *messaging.AndroidConfig
	WebPush       *messaging.WebpushConfig
	APNS          *messaging.APNSConfig
	FCMOptions    *messaging.FCMOptions
}
