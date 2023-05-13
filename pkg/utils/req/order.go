package req

import "time"

type ReqUpdateOrder struct {
	ShopOrderID   uint `json:"shop_order_id" binding:"required"`
	OrderStatusID uint `json:"order_status_id"`
}

// return request
type ReqReturn struct {
	ShopOrderID  uint   `json:"shop_order_id" binding:"required"`
	ReturnReason string `json:"return_reason" binding:"required,min=6,max=150"`
}

type UpdatOrderReturn struct {
	OrderReturnID uint      `json:"order_return_id" binding:"required"`
	OrderStatusID uint      `json:"order_status_id" binding:"required"`
	ReturnDate    time.Time `json:"return_date" binding:"omitempty"`
	AdminComment  string    `json:"admin_comment" bindin:"requied,min=6,max=150"`
}

type ReqPlaceOrder struct {
	PaymentMethodID uint `json:"payment_method_id" binding:"required"`
	AddressID       uint `json:"address_id" binding:"required"`
}

type ReqRazorpayVeification struct {
	RazorpayPaymentID string `json:"razorpay_payment_id"`
	RazorpayOrderID   string `json:"razorpay_order_id"`
	RazorpaySignature string `json:"razorpay_signature"`
	UserID            uint   `json:"user_id"`
	ShopOrderID       uint   `json:"shop_order_id"`
}
