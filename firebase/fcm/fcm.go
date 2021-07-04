package fcm

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
)

//GetNotificationManager constructs a new Notification Manager.
func GetNotificationManager(app *firebase.App, ctx context.Context) (NotificationManager, error) {
	client, err := app.Messaging(ctx)
	return NotificationManager{
		client: client,
		ctx:    ctx,
	}, err
}

//UsersSendNotification Sends the same Notification to Multiple Users.
func (nm *NotificationManager) UsersSendNotification(notificationData NotificationData, registrationTokens ...string) (*messaging.BatchResponse, error) {
	messages := &messaging.MulticastMessage{
		Tokens: registrationTokens,
		Data:   notificationData.Data,
		Notification: &messaging.Notification{
			Title:    notificationData.Title,
			Body:     notificationData.Body,
			ImageURL: notificationData.Url,
		},
		Android:    notificationData.AndroidConfig,
		Webpush:    notificationData.WebPush,
		APNS:       notificationData.APNS,
	}
	// Send a message to the device corresponding to the provided registration token.
	return nm.client.SendMulticast(nm.ctx, messages)
}

//UsersSendMessage Sends the same Message to Multiple Users.
func (nm *NotificationManager) UsersSendMessage(data map[string]string, registrationTokens ...string) (*messaging.BatchResponse, error) {
	messages := &messaging.MulticastMessage{
		Tokens: registrationTokens,
		Data:   data,
	}
	// Send a message to the device corresponding to the provided registration token.
	return nm.client.SendMulticast(nm.ctx, messages)
}

//UserSendNotification Sends a Notification to a single User.
func (nm *NotificationManager) UserSendNotification(notificationData NotificationData, registrationToken string) (string, error) {

	// See documentation on defining a message payload.
	message := &messaging.Message{
		Data: notificationData.Data,
		Notification: &messaging.Notification{
			Title:    notificationData.Title,
			Body:     notificationData.Body,
			ImageURL: notificationData.Url,
		},
		Android:    notificationData.AndroidConfig,
		Webpush:    notificationData.WebPush,
		APNS:       notificationData.APNS,
		FCMOptions: notificationData.FCMOptions,
		Token:      registrationToken,
	}

	// Send message
	return nm.client.Send(nm.ctx, message)
}

//UserSendMessage Sends a Message to a single User.
func (nm *NotificationManager) UserSendMessage(data map[string]string, registrationToken string) (string, error) {
	// See documentation on defining a message payload.
	message := &messaging.Message{
		Data:  data,
		Token: registrationToken,
	}

	// Send message
	return nm.client.Send(nm.ctx, message)
}

//ConditionSendNotification Sends a Message to all devices for which the Condition returns True.
//
//NOTE! using a raw string as a condition is not safe, it is safer to use the ConditionBuilder.
func (nm *NotificationManager) ConditionSendNotification(notificationData NotificationData, condition string) (string, error) {
	// See documentation on defining a message payload.
	message := &messaging.Message{
		Data: notificationData.Data,
		Notification: &messaging.Notification{
			Title:    notificationData.Title,
			Body:     notificationData.Body,
			ImageURL: notificationData.Url,
		},
		Android:    notificationData.AndroidConfig,
		Webpush:    notificationData.WebPush,
		APNS:       notificationData.APNS,
		FCMOptions: notificationData.FCMOptions,
		Condition: condition,
	}
	// Send message
	return nm.client.Send(nm.ctx, message)
}

//TopicSendMessage Sends a Message to all devices subscribed to the given Topic.
func (nm *NotificationManager) TopicSendMessage(data map[string]string, topic string) (string, error) {
	// See documentation on defining a message payload.
	message := &messaging.Message{
		Data:  data,
		Topic: topic,
	}
	// Send message
	return nm.client.Send(nm.ctx, message)
}

//TopicSubscribe Subscribes User/Users to a given Topic.
func (nm *NotificationManager) TopicSubscribe(topic string, registrationToken ...string) (*messaging.TopicManagementResponse, error) {
	// Subscribe the devices corresponding to the registration tokens to the topic.
	return nm.client.SubscribeToTopic(nm.ctx, registrationToken, topic)
}

//TopicUnsubscribe Unsubscribes User/Users from a given Topic.
func (nm *NotificationManager) TopicUnsubscribe(topic string, registrationToken ...string) (*messaging.TopicManagementResponse, error) {
	// Subscribe the devices corresponding to the registration tokens to the topic.
	return nm.client.UnsubscribeFromTopic(nm.ctx, registrationToken, topic)
}
