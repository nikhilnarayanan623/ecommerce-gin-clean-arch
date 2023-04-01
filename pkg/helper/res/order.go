package res

import (
	"time"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
)

type ResOrderCheckout struct {
	UserID          uint   `json:"user_id"`
	PaymentMethodID uint   `json:"payment_method_id"`
	PaymentType     string `json:"payment_type"`
	AmountToPay     uint   `json:"amount_to_pay"`
	Discount        uint   `json:"discount"`
	CouponCode      string `json:"coupon_code" `
	AddressID       uint   `json:"address_id"`
}

type ResOrder struct {
	ProductItemID uint   `json:"product_item_id"`
	ProductName   string `json:"product_name"`
	Image         string `json:""`
	Price         uint   `json:"price"`
	Qty           uint   `json:"qty"`
	SubTotal      uint   `json:"sub_total"`
	OrderDate     string `json:"order_date" `
	Status        string `json:"status"`
}

type ResShopOrder struct {
	UserID          uint       `json:"user_id"`
	ShopOrderID     uint       `json:"shop_order_id"`
	OrderDate       time.Time  `json:"order_date"`
	AddressID       uint       `json:"address_id" `
	Address         ResAddress `json:"address"`
	OrderTotalPrice uint       `json:"order_total_price" `
	Discount        uint       `json:"discount"`
	OrderStatusID   uint       `json:"order_status_id"`
	OrderStatus     string     `json:"order_status"`
	PaymentMethodID uint       `json:"payment_method_id" gorm:"primaryKey;not null"`
	PaymentType     string     `json:"" gorm:"unique;not null"`
}

// admin side
type ResShopOrdersPage struct {
	Orders   []ResShopOrder
	Statuses []domain.OrderStatus
}

// checkout
type ResCheckOut struct {
	Addresses    []ResAddress  `json:"addresses"`
	ProductItems []ResCartItem `json:"product_items"`
	TotalPrice   uint          `json:"total_price"`
}

// return
type ResOrderReturn struct {
	OrderReturnID uint      `json:"order_return_id" copier:"ID"`
	RequestDate   time.Time `json:"request_date" `
	ReturnReason  string    `json:"return_reason" `
	RefundAmount  uint      `json:"refund_amount" `

	IsApproved   bool      `json:"is_approved" `
	ReturnDate   time.Time `json:"return_date"`
	ApprovalDate time.Time `json:"approval_date"`
	AdminComment string    `json:"admin_comment"`
}
