package interfaces

import (
	"context"
	"time"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/token"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
)

//go:generate mockgen -destination=../../mock/mockUseCase/authUseCaseMock.go -package=mockUseCase . AuthUseCase
type AuthUseCase interface {
	UserLogin(ctx context.Context, loginDetails req.Login) (userID uint, err error)
	GenerateAccessToken(ctx context.Context, userID uint, userType token.UserType, expireTimeDuration time.Duration) (tokenString string, err error)
	GenerateRefreshToken(ctx context.Context, userID uint, userType token.UserType, expireTimeDuration time.Duration) (tokenString string, err error)
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
