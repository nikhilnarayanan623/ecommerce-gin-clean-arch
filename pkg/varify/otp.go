package varify

import (
	"fmt"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/config"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/verify/v2"
)

var (
	AUTHTOKEN  string
	ACCOUNTSID string
	SERVICESID string
)

var client *twilio.RestClient

func SetClient(cfg config.Config) {
	client = twilio.NewRestClientWithParams(twilio.ClientParams{
		Password: cfg.AUTHTOKEN,
		Username: cfg.ACCOUNTSID,
	})
	SERVICESID = cfg.SERVICESID
}

func TwilioSendOTP(phoneNumber string) (string, error) {
	fmt.Println(phoneNumber, AUTHTOKEN, ACCOUNTSID, SERVICESID)
	params := &twilioApi.CreateVerificationParams{}
	params.SetTo(phoneNumber)
	params.SetChannel("sms")

	resp, err := client.VerifyV2.CreateVerification(SERVICESID, params)
	if err != nil {
		return "", err
	}

	return *resp.Sid, nil
}

func TwilioVerifyOTP(phoneNumber string, code string) error {
	params := &twilioApi.CreateVerificationCheckParams{}
	params.SetTo(phoneNumber)
	params.SetCode(code)

	resp, err := client.VerifyV2.CreateVerificationCheck(SERVICESID, params)
	if err != nil {
		return err
	} else if *resp.Status == "approved" {
		return nil
	}

	return nil
}
