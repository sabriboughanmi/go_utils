package playstore

import (
	"context"
	"github.com/sabriboughanmi/go_utils/iap/playstore/androidpublisher"
)

// VerifyProduct verifies product status
func (c *Client) VerifyProduct(ctx context.Context, packageName string, productID string, token string) (*androidpublisher.ProductPurchase, error) {
	ps := androidpublisher.NewPurchasesProductsService(c.service)
	result, err := ps.Get(packageName, productID, token).Context(ctx).Do()
	return result, err
}

//AcknowledgeProduct : Acknowledges a purchase of an inapp item.
//	Note! this function must be called on all purchases within the next ~24h that follows a purchase, otherwise the purchase will be automatically voided.
//
//  - packageName: The package name of the application the inapp product was sold in (for example, 'com.some.thing').
// 	- productId: The inapp product SKU (for example, 'com.some.thing.inapp1').
// 	- token: The token provided to the user's device when the inapp product was purchased.
func (c *Client) AcknowledgeProduct(ctx context.Context, packageName, productID, token, developerPayload string) error {
	ps := androidpublisher.NewPurchasesProductsService(c.service)
	acknowledgeRequest := &androidpublisher.ProductPurchasesAcknowledgeRequest{DeveloperPayload: developerPayload}
	err := ps.Acknowledge(packageName, productID, token, acknowledgeRequest).Context(ctx).Do()
	return err
}
