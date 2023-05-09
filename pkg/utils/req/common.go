package req

import (
	"time"

	"github.com/google/uuid"
)

// login struct for user and admin
type LoginStruct struct {
	UserName string `json:"user_name" binding:"omitempty,min=3,max=15"`
	Phone    string `json:"phone" binding:"omitempty,min=10,max=10"`
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"required,min=5,max=30"`
}
type OTPLogin struct {
	Email    string `json:"email" binding:"omitempty,email"`
	UserName string `json:"user_name" binding:"omitempty,min=3,max=16"`
	Phone    string `json:"phone" binding:"omitempty,min=10,max=10"`
}

type OTPVerify struct {
	OTP   string    `json:"otp" binding:"required,min=4,max=8"`
	OTPID uuid.UUID `json:"otp_id" gorm:"not null" binding:"required"`
}

type BlockStruct struct {
	UserID uint `json:"user_id" binding:"required,numeric"`
}

type ReqPagination struct {
	Count      uint `json:"count"`
	PageNumber uint `json:"page_number"`
}

type ReqSalesReport struct {
	StartDate  time.Time     `json:"start_date"`
	EndDate    time.Time     `json:"end_date"`
	Pagination ReqPagination `json:"pagination"`
}

// stock
type ReqUpdateStock struct {
	SKU      string `json:"sku"`
	QtyToAdd uint   `json:"qty_to_add"`
}
