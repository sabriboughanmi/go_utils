package playstore

import (
	"context"
	"google.golang.org/api/androidpublisher/v3"
)

// The IABProduct type is an interface for product service
type IABProduct interface {
	VerifyProduct(context.Context, string, string, string) (*androidpublisher.ProductPurchase, error)
	AcknowledgeProduct(context.Context, string, string, string, string) error
}

// The IABSubscription type is an interface  for subscription service
type IABSubscription interface {
	AcknowledgeSubscription(context.Context, string, string, string, *androidpublisher.SubscriptionPurchasesAcknowledgeRequest) error
	VerifySubscription(context.Context, string, string, string) (*androidpublisher.SubscriptionPurchase, error)
	CancelSubscription(context.Context, string, string, string) error
	RefundSubscription(context.Context, string, string, string) error
	RevokeSubscription(context.Context, string, string, string) error
}

// Client type implements VerifySubscription method
type Client struct {
	service *androidpublisher.Service
}

// The InAppProduct : contains the *androidpublisher.InAppProduct and provides some utils.
type InAppProduct struct {
	AndroidPublisherInAppProduct *androidpublisher.InAppProduct
}
