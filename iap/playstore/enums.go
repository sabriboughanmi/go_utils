package playstore

// EProductStatus : The status of the product, e.g. whether it's active.
type EProductStatus string

const (
	ProductStatus_Unspecified EProductStatus = "statusUnspecified" // Unspecified status.
	ProductStatus_active      EProductStatus = "active"            // The product is published and active in the store.
	ProductStatus_inactive    EProductStatus = "inactive"          // The product is not published and therefore inactive in the store.
)

// ESubscriptionPeriod : specifies the Subscription period
type ESubscriptionPeriod string

const (
	SubscriptionPeriod_Invalid     ESubscriptionPeriod = ""    // Invalid Subscription (maybe a consumable).
	SubscriptionPeriod_OneWeek     ESubscriptionPeriod = "P1W" // one week.
	SubscriptionPeriod_OneMonth    ESubscriptionPeriod = "P1M" // one month.
	SubscriptionPeriod_ThreeMonths ESubscriptionPeriod = "P3M" // three months.
	SubscriptionPeriod_SixMonths   ESubscriptionPeriod = "P6M" // six months.
	SubscriptionPeriod_OneYear     ESubscriptionPeriod = "P1Y" // one year.
)

// EPurchaseType : The type of the product, e.g. a recurring subscription.
type EPurchaseType string

const (
	EPurchaseType_Unspecified  EPurchaseType = "purchaseTypeUnspecified" // Unspecified purchase type.
	EPurchaseType_ManagedUser  EPurchaseType = "managedUser"             // The default product type - Can be purchased Single/Multiple times (Consumable,Non-Consumable).
	EPurchaseType_Subscription EPurchaseType = "subscription"            // In-app product with a recurring period.
)

type EAcknowledgementState int64

const (
	AcknowledgementState_YetToBeAcknowledged EAcknowledgementState = 0
	AcknowledgementState_Acknowledged        EAcknowledgementState = 1
)

type EConsumptionState int64

const (
	ConsumptionState_YetToBeConsumed EConsumptionState = 0
	ConsumptionState_Consumed        EConsumptionState = 1
)

type EPurchaseState int64

const (
	PurchaseState_Purchased EPurchaseState = 0
	PurchaseState_Canceled  EPurchaseState = 1
	PurchaseState_Pending   EPurchaseState = 1
)
