package request

import (
	"time"
)

type OTPLogin struct {
	Email    string `json:"email" binding:"omitempty,email"`
	UserName string `json:"user_name" binding:"omitempty,min=3,max=16"`
	Phone    string `json:"phone" binding:"omitempty,min=10,max=10"`
}

type OTPVerify struct {
	Otp   string `json:"otp" binding:"required,min=4,max=8"`
	OtpID string `json:"otp_id" `
}

type BlockUser struct {
	UserID uint `json:"user_id" binding:"required,numeric"`
	Block  bool `json:"block"`
}

type SalesReport struct {
	StartDate  time.Time  `json:"start_date"`
	EndDate    time.Time  `json:"end_date"`
	Pagination Pagination `json:"pagination"`
}

// stock
type UpdateStock struct {
	SKU      string `json:"sku"`
	QtyToAdd uint   `json:"qty_to_add"`
}
