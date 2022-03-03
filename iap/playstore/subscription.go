package playstore

import (
	"context"
	"google.golang.org/api/androidpublisher/v3"
)

// AcknowledgeSubscription acknowledges a subscription purchase.
func (c *Client) AcknowledgeSubscription(ctx context.Context, packageName string, subscriptionID string, token string,
	req *androidpublisher.SubscriptionPurchasesAcknowledgeRequest) error {

	ps := androidpublisher.NewPurchasesSubscriptionsService(c.service)
	err := ps.Acknowledge(packageName, subscriptionID, token, req).Context(ctx).Do()

	return err
}

// VerifySubscription verifies subscription status
func (c *Client) VerifySubscription(ctx context.Context, packageName string, subscriptionID string, token string) (*androidpublisher.SubscriptionPurchase, error) {
	ps := androidpublisher.NewPurchasesSubscriptionsService(c.service)
	result, err := ps.Get(packageName, subscriptionID, token).Context(ctx).Do()

	return result, err
}

// CancelSubscription cancels a user's subscription purchase.
func (c *Client) CancelSubscription(ctx context.Context, packageName string, subscriptionID string, token string) error {
	ps := androidpublisher.NewPurchasesSubscriptionsService(c.service)
	err := ps.Cancel(packageName, subscriptionID, token).Context(ctx).Do()

	return err
}

// RefundSubscription refunds a user's subscription purchase, but the subscription remains valid
// until its expiration time and it will continue to recur.
func (c *Client) RefundSubscription(ctx context.Context, packageName string, subscriptionID string, token string) error {
	ps := androidpublisher.NewPurchasesSubscriptionsService(c.service)
	err := ps.Refund(packageName, subscriptionID, token).Context(ctx).Do()

	return err
}

// RevokeSubscription refunds and immediately revokes a user's subscription purchase.
// Access to the subscription will be terminated immediately and it will stop recurring.
func (c *Client) RevokeSubscription(ctx context.Context, packageName string, subscriptionID string, token string) error {
	ps := androidpublisher.NewPurchasesSubscriptionsService(c.service)
	err := ps.Revoke(packageName, subscriptionID, token).Context(ctx).Do()

	return err
}
