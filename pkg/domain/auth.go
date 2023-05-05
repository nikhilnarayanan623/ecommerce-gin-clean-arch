package domain

import (
	"time"

	"github.com/google/uuid"
)

type RefreshSession struct {
	TokenID      uuid.UUID `json:"token_id" gorm:"not null"`
	RefreshToken string    `json:"refresh_token" gorm:"not null"`
	ExpireAt     time.Time `json:"expire_at" gorm:"not null"`
	IsBlocked    bool      `json:"is_blocked" gorm:"not null;default:false"`
}
