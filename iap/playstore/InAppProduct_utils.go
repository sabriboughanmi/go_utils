package playstore

//GetStatus : The status of the product, e.g. whether it's active.
//	ProductStatus_Unspecified // Unspecified status.
//	ProductStatus_active // The product is published and active in the store.
//	ProductStatus_inactive // The product is not published and therefore inactive in the store.
func (iap InAppProduct) GetStatus() EProductStatus {
	return EProductStatus(iap.AndroidPublisherInAppProduct.Status)
}

//GetSubscriptionPeriod : specifies the Subscription period.
//  SubscriptionPeriod_Invalid : Invalid Subscription (maybe a consumable).
//	SubscriptionPeriod_OneWeek (one week).
//	SubscriptionPeriod_OneMonth (one month).
//	SubscriptionPeriod_ThreeMonths (three months).
//	SubscriptionPeriod_SixMonths (six months).
//	SubscriptionPeriod_OneYear (one year).
func (iap InAppProduct) GetSubscriptionPeriod() ESubscriptionPeriod {
	return ESubscriptionPeriod(iap.AndroidPublisherInAppProduct.SubscriptionPeriod)
}

//GetPurchaseType : The type of the product.
//	EPurchaseType_Unspecified (Unspecified purchase type).
//	EPurchaseType_ManagedUser Can be purchased Single/Multiple times (Consumable,Non-Consumable).
//	EPurchaseType_Subscription (In-app product with a recurring period).
func (iap InAppProduct) GetPurchaseType() EPurchaseType {
	return EPurchaseType(iap.AndroidPublisherInAppProduct.PurchaseType)
}
