package res

import "time"

type ResUserCoupon struct {
	ID             uint      `json:"id" gorm:"primaryKey;not null"`
	CouponCode     string    `json:"coupon_code" gorm:"unique;not null"`
	CouponName     string    `json:"coupon_name" gorm:"unique;not null" binding:"required,min=3,max=25"`
	PercentageUpto uint      `json:"percentage_upto" gorm:"not null" binding:"required,numeric,min=1,max=100"`
	MinimumPrice   uint      `json:"minimum_price" gorm:"not null" binding:"required,numeric,min=1"`
	Description    string    `json:"description" gorm:"not null" binding:"required,min=6,max=50"`
	Image          string    `json:"image" binding:"required"`
	ExpireDate     time.Time `json:"expire_date" gorm:"not null"`
	Used           bool      `json:"used" gorm:"not null"`
}
