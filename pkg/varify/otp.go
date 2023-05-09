package varify

import (
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/config"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/verify/v2"
)

// var (
// 	AUTHTOKEN  string
// 	ACCOUNTSID string
// 	SERVICESID string
// 	//client     *twilio.RestClient
// )

// func SetClient(cfg config.Config) {
// 	client = twilio.NewRestClientWithParams(twilio.ClientParams{
// 		Password: cfg.AUTHTOKEN,
// 		Username: cfg.ACCOUNTSID,
// 	})
// 	SERVICESID = cfg.SERVICESID
// }

func TwilioSendOTP(phoneNumber string) (string, error) {
	//fmt.Println(phoneNumber, AUTHTOKEN, ACCOUNTSID, SERVICESID)

	//create a twilio client with twilio details
	password := config.GetCofig().AUTHTOKEN
	userName := config.GetCofig().ACCOUNTSID
	seviceSid := config.GetCofig().SERVICESID

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Password: password,
		Username: userName,
	})
	params := &twilioApi.CreateVerificationParams{}
	params.SetTo(phoneNumber)
	params.SetChannel("sms")

	resp, err := client.VerifyV2.CreateVerification(seviceSid, params)
	if err != nil {
		return "faild to create otp verification", err
	}

	return *resp.ServiceSid, nil
}

func TwilioVerifyOTP(phoneNumber string, code string) error {
	//create a twilio client with twilio details
	password := config.GetCofig().AUTHTOKEN
	userName := config.GetCofig().ACCOUNTSID
	seviceSid := config.GetCofig().SERVICESID
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Password: password,
		Username: userName,
	})

	params := &twilioApi.CreateVerificationCheckParams{}
	params.SetTo(phoneNumber)
	params.SetCode(code)

	resp, err := client.VerifyV2.CreateVerificationCheck(seviceSid, params)
	if err != nil {
		return err
	} else if *resp.Status == "approved" {
		return nil
	}

	return nil
}
