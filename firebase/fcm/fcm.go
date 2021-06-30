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

//UsersSendMessage Sends the same Message to Multiple Users.
func (nm *NotificationManager) UsersSendMessage(title, body, url string, data map[string]string, registrationTokens ...string) (*messaging.BatchResponse, error) {
	messages := &messaging.MulticastMessage{
		Tokens: registrationTokens,
		Data:   data,
		Notification: &messaging.Notification{
			Title:    title,
			Body:     body,
			ImageURL: url,
		},
	}
	// Send a message to the device corresponding to the provided registration token.
	return nm.client.SendMulticast(nm.ctx, messages)
}

//UserSendMessage Sends a Message to a single User.
func (nm *NotificationManager) UserSendMessage(title, body, url string, data map[string]string, registrationToken string) (string, error) {

	// See documentation on defining a message payload.
	message := &messaging.Message{
		Data: data,
		Notification: &messaging.Notification{
			Title:    title,
			Body:     body,
			ImageURL: url,
		},
		Token: registrationToken,
	}

	// Send message
	return nm.client.Send(nm.ctx, message)
}

//ConditionSendMessage Sends a Message to all devices for which the Condition returns True.
//
//NOTE! using a raw string as a condition is not safe, it is safer to use the ConditionBuilder.
func (nm *NotificationManager) ConditionSendMessage(title, body, url string, data map[string]string, condition string) (string, error) {
	// See documentation on defining a message payload.
	message := &messaging.Message{
		Data: data,
		Notification: &messaging.Notification{
			Title:    title,
			Body:     body,
			ImageURL: url,
		},
		Condition: condition,
	}

	// Send message
	return nm.client.Send(nm.ctx, message)
}

//TopicSendMessage Sends a Message to all devices subscribed to the given Topic.
func (nm *NotificationManager) TopicSendMessage(title, body, url string, data map[string]string, topic string) (string, error) {
	// See documentation on defining a message payload.
	message := &messaging.Message{
		Data: data,
		Notification: &messaging.Notification{
			Title:    title,
			Body:     body,
			ImageURL: url,
		},
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
