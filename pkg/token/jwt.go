package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/config"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
)

type jwtAuth struct {
	adminSecretKey string
	userSecretKey  string
}

// New TokenAuth
func NewTokenService(cfg config.Config) TokenService {

	return &jwtAuth{
		adminSecretKey: cfg.AdminAuthKey,
		userSecretKey:  cfg.UserAuthKey,
	}
}

var (
	ErrInvalidUserType    = errors.New("invalid user type")
	ErrInvalidToken       = errors.New("invalid token")
	ErrFailedToParseToken = errors.New("failed to parse token to claims")
	ErrExpiredToken       = errors.New("token expired")
)

type jwtClaims struct {
	TokenID   string
	UserID    uint
	ExpiresAt time.Time
	// jwt.RegisteredClaims
}

// Generate a new JWT token string from token request
func (c *jwtAuth) GenerateToken(req GenerateTokenRequest) (GenerateTokenResponse, error) {

	if req.UsedFor != Admin && req.UsedFor != User {

		return GenerateTokenResponse{}, ErrInvalidUserType
	}

	tokenID := utils.GenerateUniqueString()
	claims := &jwtClaims{
		TokenID: tokenID,
		UserID:  req.UserID,
		// RegisteredClaims: jwt.RegisteredClaims{
		// 	ExpiresAt: jwt.NewNumericDate(req.ExpirationDate),
		// },
		ExpiresAt: req.ExpireAt,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	var (
		tokenString string
		err         error
	)
	// sign the token by user type
	if req.UsedFor == Admin {
		tokenString, err = token.SignedString([]byte(c.adminSecretKey))
	} else {
		tokenString, err = token.SignedString([]byte(c.userSecretKey))
	}

	if err != nil {
		return GenerateTokenResponse{}, fmt.Errorf("failed to sign the token \nerror:%w", err)
	}

	response := GenerateTokenResponse{
		TokenID:     tokenID,
		TokenString: tokenString,
	}

	return response, nil
}

// Verify JWT token string and return TokenResponse
func (c *jwtAuth) VerifyToken(req VerifyTokenRequest) (VerifyTokenResponse, error) {

	if req.UsedFor != Admin && req.UsedFor != User {
		return VerifyTokenResponse{}, ErrInvalidUserType
	}

	token, err := jwt.ParseWithClaims(req.TokenString, &jwtClaims{}, func(t *jwt.Token) (interface{}, error) {

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		if req.UsedFor == Admin {
			return []byte(c.adminSecretKey), nil
		}
		return []byte(c.userSecretKey), nil
	})

	if err != nil {
		if errors.Is(err, ErrExpiredToken) {
			return VerifyTokenResponse{}, ErrExpiredToken
		}
		return VerifyTokenResponse{}, ErrInvalidToken
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok {
		return VerifyTokenResponse{}, ErrFailedToParseToken
	}

	response := VerifyTokenResponse{
		TokenID: claims.TokenID,
		UserID:  claims.UserID,
	}
	return response, nil
}

// Validate claims
func (c *jwtClaims) Valid() error {
	if time.Since(c.ExpiresAt) > 0 {
		return ErrExpiredToken
	}
	return nil
}
