package admob

type AdNetwork string

const (
	AdColony                  = "15586990674969969776"
	AdMob                     = "5450213213286189855"
	Applovin                  = "1063618907739174004"
	Chartboost                = "2873236629771172317"
	Facebook_Audience_Network = "10568273599589928883"
	Fuse                      = "8914788932458531264"
	Fyber                     = "4839637394546996422"
	InMobi                    = "7681903010231960328"
	Maio                      = "7505118203095108657"
	MyTarget                  = "8450873672465271579"
	Nend                      = "9383070032774777750"
	Tapjoy                    = "7295217276740746030"
	Unity_Ads                 = "4970775877303683148"
	Vungle                    = "1953547073528090325"
)

//ad_network=5450213213286189855&ad_unit=1234567890&custom_data=%7B%22SuperPowers%22%3A%205%7D&reward_amount=1&reward_item=Battle%20Pass&timestamp=1624634818272&transaction_id=123456789&user_id=wNX0mXZaz5hEPWBE6StRdv4cmPk2&signature=MEQCIDhLm6-M98N4YEhPE2owusFjXS1Z7RLagfGiiEqyGyDBAiARo3sevcsrTZEPV8ZG1xBwIiPpuSt4Ege5ZBwH9hyjbQ&key_id=3335741209

type AdmobParameters struct {
	AdNetwork     AdNetwork `json:"ad_network"`
	AdUnit        string    `json:"ad_unit"`
	CustomData    string    `json:"custom_data"`
	RewardAmount  int       `json:"reward_amount"`
	RewardItem    string    `json:"reward_item"`
	Timestamp     string    `json:"timestamp"`
	TransactionID string    `json:"transaction_id"`
	UserID        string    `json:"user_id"`
	Signature     string    `json:"signature"`
	KeyID         string    `json:"key_id"`
}
