package playstore

//GetAcknowledgementState : The acknowledgement state of the inapp product.
//	AcknowledgementState_YetToBeAcknowledged // Yet To Be Acknowledged.
//	AcknowledgementState_Acknowledged // Acknowledged.
func (productPurchase ProductPurchase) GetAcknowledgementState() EAcknowledgementState {
	return EAcknowledgementState(productPurchase.ProductPurchase.AcknowledgementState)
}

//GetConsumptionState : The consumption state of the inapp product.
//	ConsumptionState_YetToBeConsumed // Yet To Be Consumed.
//	ConsumptionState_Consumed // Consumed.
func (productPurchase ProductPurchase) GetConsumptionState() EConsumptionState {
	return EConsumptionState(productPurchase.ProductPurchase.ConsumptionState)
}

//GetPurchaseState : The purchase state of the order.
//	PurchaseState_Purchased // Purchased.
//	PurchaseState_Canceled // Canceled.
//	PurchaseState_Pending // Pending.
func (productPurchase ProductPurchase) GetPurchaseState() EPurchaseState {
	return EPurchaseState(productPurchase.ProductPurchase.PurchaseState)
}
