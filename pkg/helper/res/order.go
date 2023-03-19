package res

import "time"

type ResOrder struct {
	ProductItemID uint   `json:"product_item_id"`
	ProductName   string `json:"product_name"`
	Image         string `json:""`
	Price         uint   `json:"price"`
	Qty           uint   `json:"qty"`
	TotalPrice    uint   `json:"total_price"`
	OrderDate     time.Time
	Status        string `json:"status"`
}

type ResOrderItem struct {
	ResOrder ResOrder
	Address  ResAddress
}
