package response

import "time"

type UserCoupon struct {
	CouponID   uint   `json:"coupon_id"`
	CouponCode string `json:"coupon_code" `
	CouponName string `json:"coupon_name"`

	ExpireDate       time.Time `json:"expire_date"`
	Description      string    `json:"description"`
	DiscountRate     uint      `json:"discount_rate"`
	MinimumCartPrice uint      `json:"minimum_cart_price"`
	Image            string    `json:"image" binding:"required"`
	BlockStatus      bool      `json:"block_status"`

	Used   bool      `json:"used"`
	UsedAt time.Time `json:"used_at"`
}
