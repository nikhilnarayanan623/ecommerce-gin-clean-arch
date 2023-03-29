package req

// for a new product
type ReqProduct struct {
	ID          uint   `json:"id"`
	ProductName string `json:"product_name" gorm:"not null" binding:"required,min=3,max=50"`
	Description string `json:"description" gorm:"not null" binding:"required,min=10,max=100"`
	CategoryID  uint   `json:"category_id" binding:"required"`
	Price       uint   `json:"price" gorm:"not null" binding:"required,numeric"`
	Image       string `json:"image" gorm:"not null" binding:"required"`
}

// for a new prodctItem

type ReqProductItem struct {
	ProductID         uint     `json:"product_id" binding:"required"`
	Price             uint     `json:"price" binding:"required,min=1"`
	VariationOptionID uint     `json:"variation_option_id" binding:"required"`
	QtyInStock        uint     `json:"qty_in_stock" binding:"required,min=1"`
	Images            []string `json:"images" binding:"required"`
}
