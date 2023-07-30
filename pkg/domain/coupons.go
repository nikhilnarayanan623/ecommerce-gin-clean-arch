package domain

import (
	"time"
)

type Coupon struct {
	CouponID   uint   `json:"coupon_id" gorm:"primaryKey;not null"`
	CouponName string `json:"coupon_name" gorm:"unique;not null" binding:"required,min=3,max=25"`
	CouponCode string `json:"coupon_code" gorm:"unique;not null"`

	ExpireDate       time.Time `json:"expire_date" gorm:"not null"`
	Description      string    `json:"description" gorm:"not null" binding:"required,min=6,max=150"`
	DiscountRate     uint      `json:"discount_rate" gorm:"not null" binding:"required,numeric,min=1,max=100"`
	MinimumCartPrice uint      `json:"minimum_cart_price" gorm:"not null" binding:"required,numeric,min=1"`
	Image            string    `json:"image" binding:"required"`
	BlockStatus      bool      `json:"block_status" gorm:"not null"`
	CreatedAt        time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// which is for store the user who are used coupon
type CouponUses struct {
	CouponUsesID uint      `json:"coupon_uses_id" gorm:"primaryKey;not null"`
	CouponID     uint      `json:"coupon_id" gorm:"not null"`
	Coupon       Coupon    `json:"-"`
	UserID       uint      `json:"user_id" gorm:"not null"`
	User         User      `json:"-"`
	UsedAt       time.Time `json:"used_at" gorm:"not null"`
}
