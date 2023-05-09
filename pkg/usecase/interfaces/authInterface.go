package interfaces

import (
	"context"
	"time"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/token"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
)

//go:generate mockgen -destination=../../mock/mockUseCase/authUseCaseMock.go -package=mockUseCase . AuthUseCase
type AuthUseCase interface {
	UserLogin(ctx context.Context, loginDetails req.Login) (userID uint, err error)
	UserLoginOtpSend(ctx context.Context, loginDetails req.OTPLogin) (otpRes res.OTPResponse, err error)
	LoginOtpVerify(ctx context.Context, otpVeirifyDetails req.OTPVerify) (userID uint, err error)

	GenerateAccessToken(ctx context.Context, tokenParams GenerateTokenParams) (tokenString string, err error)
	GenerateRefreshToken(ctx context.Context, tokenParams GenerateTokenParams) (tokenString string, err error)
	VerifyAndGetRefreshTokenSession(ctx context.Context, refreshToken string, usedFor token.UserType) (domain.RefreshSession, error)
}

type GenerateTokenParams struct {
	UserID     uint
	UserType   token.UserType
	ExpireDate time.Time
}

// type userType string

// const (
// 	User  userType = "user-auth"
// 	Admin userType = "admin-auth"
// )

// type CreateTokenParams struct {
// 	UserType userType
// 	UserID   uint
// 	Duration time.Duration
// }
