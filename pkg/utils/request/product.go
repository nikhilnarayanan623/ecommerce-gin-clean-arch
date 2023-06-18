package request

// for a new product
type Product struct {
	ProductName string `json:"product_name" gorm:"not null" binding:"required,min=3,max=50"`
	Description string `json:"description" gorm:"not null" binding:"required,min=10,max=100"`
	CategoryID  uint   `json:"category_id" binding:"required"`
	Price       uint   `json:"price" gorm:"not null" binding:"required,numeric"`
	Image       string `json:"image" gorm:"not null" binding:"required"`
}
type UpdateProduct struct {
	ID          uint   `json:"id" binding:"required"`
	ProductName string `json:"product_name" gorm:"not null" binding:"required,min=3,max=50"`
	Description string `json:"description" gorm:"not null" binding:"required,min=10,max=100"`
	CategoryID  uint   `json:"category_id" binding:"required"`
	Price       uint   `json:"price" gorm:"not null" binding:"required,numeric"`
	Image       string `json:"image" gorm:"not null" binding:"required"`
}

// for a new prodctItem
type ProductItem struct {
	ProductID         uint     `json:"product_id" binding:"required"`
	Price             uint     `json:"price" binding:"required,min=1"`
	VariationOptionID uint     `json:"variation_option_id" binding:"required"`
	QtyInStock        uint     `json:"qty_in_stock" binding:"required,min=1"`
	Images            []string `json:"images" binding:"required"`
}

type Variation struct {
	CategoryID    uint   `json:"category_id"  binding:"required,numeric"`
	VariationName string `json:"variation_name" binding:"required"`
}

type VariationOption struct {
	VariationID    uint   `json:"variation_id" gorm:"not null" binding:"required,numeric"` // a specific field of variation like color/size
	VariationValue string `json:"variation_value" gorm:"not null" binding:"required"`      // the variations value like blue/XL
}

// product side
type Category struct {
	CategoryName string `json:"category_name"` // new category name
	ID           uint   `json:"id"`            // id any of main category
}
