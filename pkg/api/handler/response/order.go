package response

import (
	"time"
)

type EmailAndPhone struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type UserOrder struct {
	OrderTotal  uint
	ShopOrderID uint
}

type PlaceOrder struct {
	UserID          uint   `json:"user_id"`
	PaymentMethodID uint   `json:"payment_method_id"`
	PaymentType     string `json:"payment_type"`
	AmountToPay     uint   `json:"amount_to_pay"`
	Discount        uint   `json:"discount"`
	CouponCode      string `json:"coupon_code" `
	AddressID       uint   `json:"address_id"`
}

type OrderItem struct {
	ProductItemID uint   `json:"product_item_id"`
	ProductName   string `json:"product_name"`
	Image         string `json:""`
	Price         uint   `json:"price"`
	Qty           uint   `json:"qty"`
	SubTotal      uint   `json:"sub_total"`
	OrderDate     string `json:"order_date" `
	Status        string `json:"status"`
}

type ShopOrder struct {
	UserID            uint      `json:"user_id"`
	ShopOrderID       uint      `json:"shop_order_id"`
	OrderDate         time.Time `json:"order_date"`
	AddressID         uint      `json:"address_id" `
	Address           Address   `json:"address"`
	OrderTotalPrice   uint      `json:"order_total_price" `
	Discount          uint      `json:"discount"`
	OrderStatusID     uint      `json:"order_status_id"`
	OrderStatus       string    `json:"order_status"`
	PaymentMethodID   uint      `json:"payment_method_id" gorm:"primaryKey;not null"`
	PaymentMethodName string    `json:"payment_method_name" gorm:"unique;not null"`
}

// checkout
type CheckOut struct {
	Addresses    []Address  `json:"addresses"`
	ProductItems []CartItem `json:"product_items"`
	TotalPrice   uint       `json:"total_price"`
}

// return
type OrderReturn struct {
	OrderReturnID uint      `json:"order_return_id" copier:"ID"`
	ShopOrderID   uint      `json:"shop_order_id"`
	RequestDate   time.Time `json:"request_date" `
	ReturnReason  string    `json:"return_reason" `
	RefundAmount  uint      `json:"refund_amount" `

	OrderStatusID uint      `json:"order_status_id"`
	OrderStatus   string    `json:"order_status"`
	IsApproved    bool      `json:"is_approved" `
	ReturnDate    time.Time `json:"return_date"`
	ApprovalDate  time.Time `json:"approval_date"`
	AdminComment  string    `json:"admin_comment"`
}

// razorpay
type RazorpayOrder struct {
	RazorpayKey     string      `json:"razorpay_key"`
	UserID          uint        `json:"user_id"`
	AmountToPay     uint        `json:"amount_to_pay"`
	RazorpayAmount  uint        `json:"razorpay_amount"`
	RazorpayOrderID interface{} `json:"razorpay_order_id"`
	Email           string      `json:"email"`
	Phone           string      `json:"phone"`

	ShopOrderID uint `json:"shop_order_id"`
}

type StripeOrder struct {
	ClientSecret   string `json:"client_secret"`
	PublishableKey string `json:"publishable_key"`
	AmountToPay    uint   `json:"amount_to_pay"`
	ShopOrderID    uint   `json:"shop_order_id"`
}
