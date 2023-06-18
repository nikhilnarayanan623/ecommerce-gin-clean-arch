package otp

import (
	"fmt"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/config"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/verify/v2"
)

type twilioOtp struct {
	serviceID string
	client    twilio.RestClient
}

func NewOtpAuth(cfg config.Config) OtpAuth {
	client := *twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: cfg.TwilioAccountSID,
		Password: cfg.TwilioAuthToken,
	})

	return &twilioOtp{
		serviceID: cfg.TwilioServiceID,
		client:    client,
	}
}

func (c *twilioOtp) SentOtp(phoneNumber string) (string, error) {

	params := &twilioApi.CreateVerificationParams{}
	params.SetTo(phoneNumber)
	params.SetChannel("sms")

	resp, err := c.client.VerifyV2.CreateVerification(c.serviceID, params)
	if err != nil {
		return "", err
	}

	return *resp.Sid, nil
}

func (c *twilioOtp) VerifyOtp(phoneNumber string, code string) (valid bool, err error) {

	params := &twilioApi.CreateVerificationCheckParams{}
	params.SetTo(phoneNumber)
	params.SetCode(code)

	resp, err := c.client.VerifyV2.CreateVerificationCheck(c.serviceID, params)

	if err != nil {
		return false, err
	}
	if resp != nil && *resp.Status != "approved" {
		return false, fmt.Errorf("invalid otp")
	}

	return true, nil
}
