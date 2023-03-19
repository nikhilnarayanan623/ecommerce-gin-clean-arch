package domain

import (
	"time"
)

type OrderStatus struct {
	ID     uint   `json:"id" gorm:"primaryKey;not null"`
	Status string `json:"status" gorm:"unique;not null"`
}

type ShopOrder struct {
	ID              uint `josn:"id" gorm:"primaryKey;not null"`
	UserID          uint `json:"user_id" gorm:"not null"`
	User            User
	OrderDate       time.Time `json:"order_date" gorm:"not null"`
	AddressID       uint      `json:"address_id" gorm:"not null"`
	Address         Address
	OrderTotalPrice uint `json:"order_total_price" gorm:"not null"`
	OrderStatusID   uint `json:"order_status_id" gorm:"not nulll"`
	COD             bool `json:"cod"`
}

type OrderLine struct {
	ID            uint `json:"id" gorm:"primaryKey;not null"`
	ProductItemID uint `json:"proudct_item_id" gorm:"not null"`
	ShopOrderID   uint `json:"shop_order_id" gorm:"not null"`
	ShopOrder     ShopOrder
	Qty           uint `json:"qty" gorm:"not null"`
	Price         uint `json:"price" gorm:"not null"`
}
