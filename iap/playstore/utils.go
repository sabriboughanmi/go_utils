package playstore

import "context"

// GetProduct : Gets an in-app product, which can be a managed product or a subscription.
//  - packageName: Package name of the app.
//  - productID: Unique identifier for the in-app product (sku).
func (c *Client) GetProduct(ctx context.Context, packageName string, productID string) (*InAppProduct, error) {
	var iap, err = c.service.Inappproducts.Get(packageName, productID).Context(ctx).Do()
	return &InAppProduct{iap}, err
}
