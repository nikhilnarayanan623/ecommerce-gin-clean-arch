package req

type ReqUpdateOrder struct {
	ShopOrderID   uint `json:"shop_order_id" binding:"required"`
	OrderStatusID uint `json:"order_status_id"`
}

// return request
type ReqReturn struct {
	ShopOrderID  uint   `json:"shop_order_id" binding:"required"`
	ReturnReason string `json:"return_reason" binding:"required,min=6,max=50"`
}

type ReqUpdatReturnReq struct {
	OrderReturnID uint   `json:"order_return_id" binding:"required"`
	OrderStatusID uint   `json:"order_status_id" binding:"required"`
	AdminComment  string `json:"admin_comment" bindin:"requied,min=6,max=50"`
}

type ReqCheckout struct {
	UserID          uint   `json:"-"`
	PaymentMethodID uint   `json:"payment_method_id" binding:"required"`
	CouponCode      string `json:"coupon_code" `
	AddressID       uint   `json:"address_id" binding:"required"`
}

type ReqRazorpayVeification struct {
	RazorpayPaymentID string `json:"razorpay_payment_id"`
	RazorpayOrderID   string `json:"razorpay_order_id"`
	RazorpaySignature string `json:"razorpay_signature"`
	UserID            uint   `json:"user_id"`
	ShopOrderID       uint   `json:"shop_order_id"`
}
