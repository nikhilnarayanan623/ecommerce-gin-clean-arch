package res

import "time"

var ResoposeMap map[string]string

// admin
type ResAdminLogin struct {
	ID       uint   `json:"id" `
	UserName string `json:"user_name"`
	Email    string `json:"email"`
}

// response of category for showing the category
type RespCategory struct {
	ID               uint   `json:"id"`
	CategoryName     string `json:"category_name"`
	CategoryID       uint   `json:"category_id"`
	MainCategoryName string `json:"main_category_name"`
}

// reponse for get all variations with its respective category

type SalesReport struct {
	UserID          uint      `json:"user_id"`
	ShopOrderID     uint      `json:"order_id"`
	OrderDate       time.Time `json:"order_date"`
	OrderTotalPrice uint      `json:"order_total_price"`
	Discount        uint      `json:"discount_price"`
	OrderStatus     string    `json:"order_status"`
	PaymentType     string    `json:"payment_type"`
}