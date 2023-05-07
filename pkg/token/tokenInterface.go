package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

//go:generate mockgen -destination=../mock/mockTokenAuth/tokenAuthMock.go -package=mockTokenAuth . TokenAuth
type TokenAuth interface {
	CreateToken(payload *Payload, user UserType) (tokenString string, err error)
	VerifyToken(tokenString string, usedFor UserType) (payload Payload, err error)
}

type UserType string

const (
	TokenForAdmin UserType = "admin"
	TokenForUser  UserType = "user"
)

var (
	errExpiredToken = errors.New("token expired")
	errInvalidToken = errors.New("invalid token")
)

type Payload struct {
	TokenID  uuid.UUID `json:"token_id"`
	UserID   uint      `json:"user_id"`
	ExpireAt time.Time `json:"expire_at"`
}

func (c *Payload) Valid() error {
	if time.Since(c.ExpireAt) > 0 {
		return errExpiredToken
	}
	return nil
}
