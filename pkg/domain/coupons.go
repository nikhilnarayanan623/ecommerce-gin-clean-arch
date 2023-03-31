package domain

import "time"

type Coupon struct {
	ID             uint   `json:"id" gorm:"primaryKey;not null"`
	CouponName     string `json:"coupon_name" gorm:"unique;not null" binding:"required,min=3,max=25"`
	Description    string `json:"description" gorm:"not null" binding:"required,min=6,max=50"`
	PercentageUpto uint   `json:"percentage_upto" gorm:"not null" binding:"required,numeric,min=1,max=100"`
	MinimumPrice   uint   `json:"minimum_price" gorm:"not null" binding:"required,numeric,min=1"`
	Image          string `json:"image" binding:"required"`
	BlockStatus    bool   `json:"block_status" gorm:"not null"`
}

type UserCoupon struct {
	ID             uint      `json:"id" gorm:"primaryKey;not null"`
	CouponCode     string    `json:"coupon_code" gorm:"unique;not null"`
	UserID         uint      `json:"user_id" gorm:"not null"`
	User           User      `json:"-"`
	CouponID       uint      `json:"coupon_id" gorm:"not null"`
	Coupon         Coupon    `json:"-"`
	DiscountAmount uint      `json:"discount_amount"`
	ExpireDate     time.Time `json:"expire_date" gorm:"not null"`
	Used           bool      `json:"used" gorm:"not null"`
	LastApplied    time.Time `json:"last_applied"`
}
