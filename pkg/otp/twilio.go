package otp

import (
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/config"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/verify/v2"
)

type twilioOtp struct {
	cfg config.Config
}

func NewTwiloOtp(cfg config.Config) OtpVerification {
	return &twilioOtp{
		cfg: cfg,
	}
}

func (c *twilioOtp) SentOtp(phoneNumber string) (string, error) {

	client := c.getNewTwiloClient()
	params := &twilioApi.CreateVerificationParams{}
	params.SetTo(phoneNumber)
	params.SetChannel("sms")

	seviceSid := c.cfg.SERVICESID
	resp, err := client.VerifyV2.CreateVerification(seviceSid, params)
	if err != nil {
		return "", err
	}

	return *resp.Sid, nil
}

func (c *twilioOtp) VerifyOtp(phoneNumber string, code string) error {

	client := c.getNewTwiloClient()

	params := &twilioApi.CreateVerificationCheckParams{}
	params.SetTo(phoneNumber)
	params.SetCode(code)

	seviceSid := c.cfg.SERVICESID
	resp, err := client.VerifyV2.CreateVerificationCheck(seviceSid, params)
	if err != nil {
		return err
	} else if *resp.Status == "approved" {
		return nil
	}

	return nil
}

func (c *twilioOtp) getNewTwiloClient() twilio.RestClient {
	password := c.cfg.AUTHTOKEN
	userName := c.cfg.ACCOUNTSID

	return *twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: userName,
		Password: password,
	})
}
