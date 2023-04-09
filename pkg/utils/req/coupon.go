package req

import "time"

type ReqCoupon struct {
	CouponName  string `json:"coupon_name" binding:"required,min=3,max=25"`
	Description string `json:"description"  binding:"required,min=6,max=150"`

	ExpireDate       time.Time `json:"expire_date" binding:"required"`
	DiscountRate     uint      `json:"discount_rate"  binding:"required,numeric,min=1,max=100"`
	MinimumCartPrice uint      `json:"minimum_cart_price"  binding:"required,numeric,min=1"`
	Image            string    `json:"image" binding:"required"`
	BlockStatus      bool      `json:"block_status"`
}

type ReqEditCoupon struct {
	CouponID    uint   `json:"coupon_id"`
	CouponName  string `json:"coupon_name" binding:"required,min=3,max=25"`
	Description string `json:"description"  binding:"required,min=6,max=150"`

	ExpireDate       time.Time `json:"expire_date" binding:"required"`
	DiscountRate     uint      `json:"discount_rate"  binding:"required,numeric,min=1,max=100"`
	MinimumCartPrice uint      `json:"minimum_cart_price"  binding:"required,numeric,min=1"`
	Image            string    `json:"image" binding:"required"`
	BlockStatus      bool      `json:"block_status"`
}

type ReqApplyCoupon struct {
	CouponCode string `json:"coupon_code" binding:"required"`
}
