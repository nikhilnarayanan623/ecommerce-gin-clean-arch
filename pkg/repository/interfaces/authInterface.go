package interfaces

import (
	"context"

	"github.com/google/uuid"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
)

//go:generate mockgen -destination=../../mock/mockRepository/authRepoMock.go -package=mockRepository . AuthRepository
type AuthRepository interface {
	SaveRefreshSession(ctx context.Context, refreshSession domain.RefreshSession) error
	FindRefreshSessionByTokenID(ctx context.Context, tokenID uuid.UUID) (domain.RefreshSession, error)

	SaveOtpSession(ctx context.Context,otpSession domain.OtpSession)error
	FindOtpSession(ctx context.Context,otpID uuid.UUID)(domain.OtpSession,error)
}
