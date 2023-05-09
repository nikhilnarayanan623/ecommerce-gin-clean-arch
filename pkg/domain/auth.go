package domain

import (
	"time"

	"github.com/google/uuid"
)

type RefreshSession struct {
	TokenID      uuid.UUID `json:"token_id" gorm:"not null"`
	UserID       uint      `json:"user_id" gorm:"not null"`
	RefreshToken string    `json:"refresh_token" gorm:"not null"`
	ExpireAt     time.Time `json:"expire_at" gorm:"not null"`
	IsBlocked    bool      `json:"is_blocked" gorm:"not null;default:false"`
}

type OtpSession struct {
	ID       uint      `json:"id" gorm:"primaryKey;not null"`
	OTPID    uuid.UUID `json:"otp_id" gorm:"not null"`
	UserID   uint      `json:"user_id" gorm:"not null"`
	Phone    string    `json:"phone" gorm:"not null"`
	ExpireAt time.Time `json:"expire_at" gorm:"not null"`
}
