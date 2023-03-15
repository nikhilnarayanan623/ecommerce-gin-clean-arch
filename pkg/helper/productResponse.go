package helper

// response for product
type ResponseProduct struct {
	ID           uint   `json:"id"`
	ProductName  string `json:"product_name"`
	Description  string `json:"description" `
	CategoryName string `json:"category_name"`
	Price        uint   `json:"price"`
	Image        string `json:"image"`
}

// fo a spedific category representation
type Category struct {
	ID               uint   `json:"id"`
	CategoryName     string `json:"category_name"`
	CategoryID       uint   `json:"category_id"`
	MainCategoryName string `json:"main_category_name"`
}

// fo a spedific variation representation
type VariationName struct {
	ID            uint   `json:"id"`
	VariationName string `json:"variation_name"`
	CategoryID    uint   `json:"category_id"`
	CategoryName  string `json:"category_name"`
}

// fo a spedific variation Value representation
type VariationValue struct {
	ID             uint   `json:"id"`
	VariationValue string `json:"variation_value"`
	VariationID    uint   `json:"variation_id"`
	VariationName  string `json:"variation_name"`
}

// fo all category, variation, variation_value
type RespFullCategory struct {
	Category       []Category
	VariationName  []VariationName
	VariationValue []VariationValue
}

// for reponse a specific products all product items
type RespProductItems struct {
	ID          uint   `json:"id"`
	ProductName string `json:"product_name"`
	ProductID   uint   `json:"product_id"`
	Price       uint   `json:"price"`
	QtyInStock  uint   `json:"qty_in_stock"`

	VariationOptionID uint   `json:"variation_option_id"`
	VariationValue     string `json:"variation_value"`
}
