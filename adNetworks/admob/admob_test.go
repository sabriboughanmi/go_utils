package admob

import (
	"fmt"
	"net/url"
	"testing"
)

const sampleReqURL = "/?ad_network=5450213213286189855&ad_unit=1234567890&custom_data=anything&reward_amount=1&reward_item=Battle%20Pass&timestamp=1624708005438&transaction_id=123456789&user_id=wNX0mXZaz5hEPWBE6StRdv4cmPk2&signature=MEQCIE2YlbRspqD_lcWJz1KwI19CV-dsB3r6iDJJDuqpbg-9AiB-B_IZoyMLYYnxU-5DZlPyToYj132X6rJQSEJokiAAQA&key_id=3335741209"



func TestVerifyURL(t *testing.T) {
	sampleRequestURL, _ := url.ParseRequestURI(sampleReqURL)
	if err := VerifyURL(sampleRequestURL); err != nil {
		t.Errorf("Error Verifying Request URL - %v", err)
	} else {
		fmt.Println("func VerifyURL OK!")
	}
}

func TestGetParametersURL(t *testing.T) {
	ssvCallback, err := GetParameters(sampleReqURL)
	if err != nil {
		t.Errorf("Error Getting Parameters from url - %v", err)
	}

	if ssvCallback.AdNetwork == "" {
		t.Errorf("Missing AdNetwork Field")
	}
	if ssvCallback.AdUnit == 0 {
		t.Errorf("Missing AdUnit Field")
	}

	if ssvCallback.CustomData == "" {
		t.Errorf("Missing CustomData Field")
	}

	if ssvCallback.RewardAmount == 0 {
		t.Errorf("Missing RewardAmount Field")
	}

	if ssvCallback.RewardItem == "" {
		t.Errorf("Missing RewardItem Field")
	}

	if ssvCallback.Timestamp == 0 {
		t.Errorf("Missing Timestamp Field")
	}

	if ssvCallback.TransactionID == "" {
		t.Errorf("Missing TransactionID Field")
	}

	if ssvCallback.UserID == "" {
		t.Errorf("Missing UserID Field")
	}

	if ssvCallback.Signature == "" {
		t.Errorf("Missing Signature Field")
	}

	if ssvCallback.KeyID == 0 {
		t.Errorf("Missing KeyID Field")
	}

}
