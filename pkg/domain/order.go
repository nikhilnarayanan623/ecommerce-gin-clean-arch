package domain

import (
	"time"
)

type PaymentMethod struct {
	ID            uint   `json:"id" gorm:"primaryKey;not null"`
	PaymentType   string `json:"" gorm:"unique;not null"`
	BlockStatus   bool   `json:"block_status" gorm:"not null"`
	MaximumAmount uint   `json:"maximum_amount" gorm:"not null"`
}

type OrderStatus struct {
	ID     uint   `json:"id" gorm:"primaryKey;not null"`
	Status string `json:"status" gorm:"unique;not null"`
}
type ShopOrder struct {
	ID              uint          `josn:"shop_order_id" gorm:"primaryKey;not null"`
	UserID          uint          `json:"user_id" gorm:"not null"`
	User            User          `json:"-"`
	OrderDate       time.Time     `json:"order_date" gorm:"not null"`
	AddressID       uint          `json:"address_id" gorm:"not null"`
	Address         Address       `json:"-"`
	OrderTotalPrice uint          `json:"order_total_price" gorm:"not null"`
	Discount        uint          `json:"discount" gorm:"not null"`
	OrderStatusID   uint          `json:"order_status_id" gorm:"not null"`
	OrderStatus     OrderStatus   `json:"-"`
	PaymentMethodID uint          `json:"payment_method_id" gorm:"not null"`
	PaymentMethod   PaymentMethod `json:"-"`
}

type OrderLine struct {
	ID            uint      `json:"id" gorm:"primaryKey;not null"`
	ProductItemID uint      `json:"proudct_item_id" gorm:"not null"`
	ShopOrderID   uint      `json:"shop_order_id" gorm:"not null"`
	ShopOrder     ShopOrder `json:"-"`
	Qty           uint      `json:"qty" gorm:"not null"`
	Price         uint      `json:"price" gorm:"not null"`
}

type OrderReturn struct {
	ID           uint      `json:"id" gorm:"primaryKey;not null"`
	ShopOrderID  uint      `json:"shop_order_id" gorm:"not null;unique"`
	ShopOrder    ShopOrder `json:"-"`
	RequestDate  time.Time `json:"request_date" gorm:"not null"`
	ReturnReason string    `json:"return_reason" gorm:"not null"`
	RefundAmount uint      `json:"refund_amount" gorm:"not null"`

	IsApproved   bool      `json:"is_approved"`
	ReturnDate   time.Time `json:"return_date"`
	ApprovalDate time.Time `json:"approval_date"`
	AdminComment string    `json:"admin_comment"`
}
