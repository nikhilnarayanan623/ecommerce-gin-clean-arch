package request

import "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"

type PaymentMethod struct {
	PaymentType   string `json:"payment_type" binding:"required,min=2,max=20"`
	BlockStatus   bool   `json:"block_status" binding:"omitempty"`
	MaximumAmount uint   `json:"maximum_amount" binding:"required,min=1,max=500000"`
}

type PaymentMethodUpdate struct {
	ID            uint `json:"id" binding:"required"`
	BlockStatus   bool `json:"block_status" binding:"omitempty"`
	MaximumAmount uint `json:"maximum_amount" binding:"required,min=1,max=500000"`
}

type RazorpayVerify struct {
	OrderID   string `json:"razorpay_order_id"`
	PaymentID string `json:"razorpay_payment_id"`
	Signature string `json:"razorpay_signature"`
}

type ApproveOrder struct {
	ShopOrderID uint
	PaymentType domain.PaymentType
}
