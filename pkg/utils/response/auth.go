package response

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type OTPResponse struct {
	OtpID string `json:"otp_id"`
}
