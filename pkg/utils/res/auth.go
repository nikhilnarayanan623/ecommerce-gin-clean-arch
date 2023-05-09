package res

import "github.com/google/uuid"

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type OTPResponse struct {
	OTPID uuid.UUID `json:"otp_id" gorm:"not null"`
}
