package admob

import (
	"fmt"
	"net/url"
	"testing"
)

func TestVerifyURL(t *testing.T) {
	sampleURL := "/?ad_network=5450213213286189855&ad_unit=1234567890&custom_data=%7B%22SuperPowers%22%3A%205%7D&reward_amount=1&reward_item=Battle%20Pass&timestamp=1624634818272&transaction_id=123456789&user_id=wNX0mXZaz5hEPWBE6StRdv4cmPk2&signature=MEQCIDhLm6-M98N4YEhPE2owusFjXS1Z7RLagfGiiEqyGyDBAiARo3sevcsrTZEPV8ZG1xBwIiPpuSt4Ege5ZBwH9hyjbQ&key_id=3335741209"
	sampleRequestURL, _ := url.ParseRequestURI(sampleURL)
	if err := VerifyURL(sampleRequestURL); err != nil {
		t.Errorf("Error Verifying Request URL - %v", err)
	} else {
		fmt.Println("func VerifyURL OK!")
	}
}
