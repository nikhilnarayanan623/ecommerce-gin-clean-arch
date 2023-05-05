package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type UserType string

const (
	TokenForAdmin UserType = "admin"
	TokenForUser  UserType = "user"
)

type TokenAuth interface {
	CreateToken(payload *Payload, user UserType) (tokenString string, err error)
	VerifyToken(tokenString string, user UserType) (payload *Payload, err error)
}

type jwtAuth struct {
	adminSecretKey string
	userSecretKey  string
}

func NewJWTAuth(adminSecretKey, userSecretKey string) TokenAuth {
	return &jwtAuth{
		adminSecretKey: adminSecretKey,
		userSecretKey:  userSecretKey,
	}
}

type Payload struct {
	TokenID  uuid.UUID `json:"token_id"`
	UserID   uint      `json:"user_id"`
	ExpireAt time.Time `json:"expire_at"`
}

var (
	errExpiredToken = errors.New("token expired")
	errInvalidToken = errors.New("invalid token")
)

func (c *Payload) Valid() error {
	if time.Since(c.ExpireAt) > 0 {
		return errExpiredToken
	}
	return nil
}

func (c *jwtAuth) CreateToken(payload *Payload, usedFor UserType) (tokenString string, err error) {

	if payload == nil {
		return "", errors.New("payload should not be nil")
	}

	var signInKey []byte

	switch usedFor {
	case TokenForAdmin:
		signInKey = []byte(c.adminSecretKey)
	case TokenForUser:
		signInKey = []byte(c.userSecretKey)
	default:
		return tokenString, errors.New("invalid user_type")
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	tokenString, err = jwtToken.SignedString(signInKey)
	if err != nil {
		return "", errors.New("faild to sign the token with sign key")
	}

	return tokenString, nil
}

func (c *jwtAuth) VerifyToken(tokenString string, usedFor UserType) (payload *Payload, err error) {

	var signInKey []byte

	switch usedFor {
	case TokenForAdmin:
		signInKey = []byte(c.adminSecretKey)
	case TokenForUser:
		signInKey = []byte(c.userSecretKey)
	default:
		return nil, errors.New("invalid user_type")
	}

	jwtToken, err := jwt.ParseWithClaims(tokenString, &Payload{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errInvalidToken
		}
		return signInKey, nil
	})

	if err != nil {
		validationErr, ok := err.(*jwt.ValidationError)

		if ok && errors.Is(validationErr.Inner, errExpiredToken) {
			return nil, errExpiredToken
		}
		return nil, errInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, errInvalidToken
	}

	return
}
