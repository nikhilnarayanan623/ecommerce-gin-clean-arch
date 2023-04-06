package req

// offer
type ReqOfferCategory struct {
	OfferID    uint `json:"offer_id" binding:"required"`
	CategoryID uint `json:"category_id" binding:"required"`
}

type ReqOfferProduct struct {
	OfferID   uint `json:"offer_id" binding:"required"`
	ProductID uint `json:"product_id" binding:"required"`
}

type ReqApplyCoupon struct {
	CouponCode string `json:"coupon_code" binding:"required"`
}
