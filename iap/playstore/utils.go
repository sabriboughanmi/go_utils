package playstore

import (
	"context"
	"google.golang.org/api/androidpublisher/v3"
)

// GetProduct : Gets an in-app product, which can be a managed product or a subscription.
//  - packageName: Package name of the app.
//  - productID: Unique identifier for the in-app product (sku).
func (c *Client) GetProduct(ctx context.Context, packageName string, productID string) (*InAppProduct, error) {
	var iap, err = c.service.Inappproducts.Get(packageName, productID).Context(ctx).Do()
	return &InAppProduct{iap}, err
}

// DefaultPriceToMoney .
func DefaultPriceToMoney(currency string, priceMicros string) *androidpublisher.Money {
	return &androidpublisher.Money{
		CurrencyCode:    currency,
		Nanos:           0,
		Units:           0,
		ForceSendFields: nil,
		NullFields:      nil,
	}
}

//ConvertRegionPrices : Calculates the region prices, using today's exchange rate and country-specific pricing patterns, based on the price in the request for a set of regions.
func (c *Client) ConvertRegionPrices(ctx context.Context, packageName string, price *androidpublisher.Money) (*androidpublisher.ConvertRegionPricesResponse, error) {
	monetizationService := androidpublisher.NewMonetizationService(c.service)
	convertRegionPricesRequest := &androidpublisher.ConvertRegionPricesRequest{
		Price: price,
	}
	return monetizationService.ConvertRegionPrices(packageName, convertRegionPricesRequest).Context(ctx).Do()
}
