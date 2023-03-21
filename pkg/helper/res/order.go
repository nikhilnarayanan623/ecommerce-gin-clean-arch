package res

import (
	"time"
)

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
	UserID          uint `json:"user_id"`
	ShopOrderID     uint `json:"shop_order_id"`
	OrderDate       time.Time
	AddressID       uint       `json:"address_id" `
	Address         ResAddress `json:"address"`
	OrderTotalPrice uint       `json:"order_total_price" `
	OrderStatusID   uint       `json:"order_status_id"`
	OrderStatus     string     `json:"order_status"`
	COD             bool       `json:"cod"`
}

type ResCheckOut struct {
	Addresses    []ResAddress       `json:"addresses"`
	ProductItems []ResponseCartItem `json:"product_items"`
	TotalPrice   uint               `json:"total_price"`
}
