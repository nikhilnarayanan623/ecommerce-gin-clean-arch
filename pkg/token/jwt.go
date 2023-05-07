package token

import (
	"errors"

	"github.com/golang-jwt/jwt"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/config"
)

type jwtAuth struct {
	adminSecretKey string
	userSecretKey  string
}

func NewJWTAuth(cfg config.Config) TokenAuth {
	return &jwtAuth{
		adminSecretKey: cfg.JWTAdmin,
		userSecretKey:  cfg.JWTUser,
	}
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

func (c *jwtAuth) VerifyToken(tokenString string, usedFor UserType) (payload Payload, err error) {

	var signInKey []byte

	switch usedFor {
	case TokenForAdmin:
		signInKey = []byte(c.adminSecretKey)
	case TokenForUser:
		signInKey = []byte(c.userSecretKey)
	default:
		return payload, errors.New("invalid user_type")
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
			return payload, errExpiredToken
		}
		return payload, errInvalidToken
	}

	convertedPayload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return payload, errInvalidToken
	}

	payload = *convertedPayload

	return
}
