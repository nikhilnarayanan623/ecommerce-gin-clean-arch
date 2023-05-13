package otp

type OtpVerification interface {
	SentOtp(phoneNumber string) (string, error)
	VerifyOtp(phoneNumber string, code string) error
}
