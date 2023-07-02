package token

import "time"

type TokenService interface {
	GenerateToken(req GenerateTokenRequest) (GenerateTokenResponse, error)
	VerifyToken(req VerifyTokenRequest) (VerifyTokenResponse, error)
}

type UserType string

const (
	Admin UserType = "admin"
	User  UserType = "user"
)

type GenerateTokenRequest struct {
	UserID   uint
	UsedFor  UserType
	ExpireAt time.Time
}

type GenerateTokenResponse struct {
	TokenID     string
	TokenString string
}

type VerifyTokenRequest struct {
	TokenString string
	UsedFor     UserType
}

type VerifyTokenResponse struct {
	TokenID string
	UserID  uint
}
