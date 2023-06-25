package request

// for a new product
type Product struct {
	Name        string `json:"product_name" binding:"required,min=3,max=50"`
	Description string `json:"description" binding:"required,min=10,max=100"`
	CategoryID  uint   `json:"category_id" binding:"required"`
	Price       uint   `json:"price" binding:"required,numeric"`
	Image       string `json:"image" binding:"required"`
}
type UpdateProduct struct {
	ID          uint   `json:"product_id" binding:"required"`
	Name        string `json:"product_name" binding:"required,min=3,max=50"`
	Description string `json:"description" binding:"required,min=10,max=100"`
	CategoryID  uint   `json:"category_id" binding:"required"`
	Price       uint   `json:"price" binding:"required,numeric"`
	Image       string `json:"image" binding:"required"`
}

// for a new productItem
type ProductItem struct {
	ProductID         uint     `json:"product_id" binding:"required"`
	Price             uint     `json:"price" binding:"required,min=1"`
	VariationOptionID []uint   `json:"variation_option_id" binding:"required"`
	QtyInStock        uint     `json:"qty_in_stock" binding:"required,min=1"`
	SKU               string   `json:"-"`
	Images            []string `json:"images" binding:"required,dive,min=1"`
}

type Variation struct {
	Names []string `json:"variation_names" binding:"required,dive,min=1"`
}

type VariationOption struct {
	Values []string `json:"variation_value" binding:"required,dive,min=1"`
}

type Category struct {
	Name string `json:"category_name" binding:"required"`
}

type SubCategory struct {
	CategoryID uint   `json:"category_id" binding:"required"`
	Name       string `json:"category_name" binding:"required"`
}
