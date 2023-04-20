package req

type ReqPaymentMethod struct {
	PaymentType   string `json:"payment_type" binding:"required,min=2,max=20"`
	BlockStatus   bool   `json:"block_status" binding:"omitempty"`
	MaximumAmount uint   `json:"maximum_amount" binding:"required,min=1,max=500000"`
}

type ReqPaymentMethodUpdate struct {
	ID            uint   `json:"id" binding:"required"`
	PaymentType   string `json:"payment_type" binding:"required,min=2,max=20"`
	BlockStatus   bool   `json:"block_status" binding:"omitempty"`
	MaximumAmount uint   `json:"maximum_amount" binding:"required,min=1,max=500000"`
}
