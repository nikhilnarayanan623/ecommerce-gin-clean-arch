package domain

import (
	"time"
)

type RefreshSession struct {
	TokenID      string    `json:"token_id" gorm:"primaryKey;not null"`
	UserID       uint      `json:"user_id" gorm:"not null"`
	RefreshToken string    `json:"refresh_token" gorm:"not null"`
	ExpireAt     time.Time `json:"expire_at" gorm:"not null"`
	IsBlocked    bool      `json:"is_blocked" gorm:"not null;default:false"`
}

type OtpSession struct {
	ID       uint      `json:"id" gorm:"primaryKey;not null"`
	OtpID    string    `json:"otp_id" gorm:"unique;not null"`
	UserID   uint      `json:"user_id" gorm:"not null"`
	Phone    string    `json:"phone" gorm:"not null"`
	ExpireAt time.Time `json:"expire_at" gorm:"not null"`
}
