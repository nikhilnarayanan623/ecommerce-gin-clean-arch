package response

import "time"

var ResoposeMap map[string]string

// admin
type AdminLogin struct {
	ID       uint   `json:"id" `
	UserName string `json:"user_name"`
	Email    string `json:"email"`
}

// reponse for get all variations with its respective category

type SalesReport struct {
	UserID          uint      `json:"user_id"`
	FirstName       string    `json:"first_name"`
	Email           string    `json:"email"`
	ShopOrderID     uint      `json:"order_id"`
	OrderDate       time.Time `json:"order_date"`
	OrderTotalPrice uint      `json:"order_total_price"`
	Discount        uint      `json:"discount_price"`
	OrderStatus     string    `json:"order_status"`
	PaymentType     string    `json:"payment_type"`
}

type Stock struct {
	ProductItemID    uint              `json:"product_item_id"`
	ProductName      string            `json:"product_name"`
	Price            uint              `json:"price"`
	SKU              string            `json:"sku"`
	QtyInStock       uint              `json:"qty_in_stock"`
	VariationOptions []VariationOption `gorm:"-"`
}
