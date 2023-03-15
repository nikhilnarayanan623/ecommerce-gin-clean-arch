package domain

// represent a model of product
type Product struct {
	ID          uint     `json:"id" gorm:"primaryKey;not null"`
	ProductName string   `json:"product_name" gorm:"not null" binding:"required,min=3,max=50"`
	Description string   `json:"description" gorm:"not null" validate:"required,min=10,max=100"`
	CategoryID  uint     `json:"category_id" binding:"omitempty,numeric"`
	Category    Category `json:"-"` //when binding this fields inside taking so added new requbody on helper
	Price       uint     `json:"price" gorm:"not null" validte:"required,numeric"`
	Image       string   `json:"image" gorm:"not null"`
}

// this for a specift variant of product
type ProductItem struct {
	ID        uint `json:"id" gorm:"primaryKey;not null"`
	ProductID uint `json:"product_id" gorm:"not null" binding:"required,numeric"`
	Product   Product
	//images are stored in sperate table along with productItem Id
	QtyInStock uint `json:"qty_in_stck" gorm:"not null" binding:"required,numeric"` // no need of stockAvailble column , because from this qty we can get it
	Price      uint `json:"price" gorm:"not null" binding:"required,numeric"`
}

// for a products category main and sub category as self joining
type Category struct {
	ID           uint      `json:"id" gorm:"primaryKey;not null"`
	CategoryID   uint      `json:"category_id"`
	Category     *Category // self join for sub category
	CategoryName string    `json:"category_name" gorm:"unique;not null" binding:"required,min=1,max=30"`
}

// variation means size color etc..
type Variation struct {
	ID            uint `json:"id" gorm:"primaryKey;not null"`
	CategoryID    uint `json:"category_id" gorm:"not null" binding:"required,numeric"`
	Category      *Category
	VariationName string `json:"variation_name" gorm:"not null" binding:"required"`
}

// variation option means values are like s,m,xl for size and blue,white,black for Color
type VariationOption struct {
	ID             uint `json:"id" gorm:"primaryKey;not null"`
	VariationID    uint `json:"variation_id" gorm:"not null" binding:"required,numeric"` // a specific field of variation like color/size
	Variation      *Variation
	VariationValue string `json:"variation_value" gorm:"not null" binding:"required"` // the variations value like blue/XL
}

// used to many to many join like multile product have same size or color and many color or size have same product
// this configuraion means to connect a specifc product to Its variasionOption(jeans-size-m)
type ProductConfiguration struct {
	ProductItemID     uint `json:"product_item_id" gorm:"not null"`
	ProductItem       ProductItem
	VariationOptionID uint `json:"variation_option_id" gorm:"not null"`
	VariationOption   VariationOption
}

// to store a url of productItem Id along a unique url
// so we can sote multiple imagesurl for a ProductItem
// one to many connection
type ProductImage struct {
	ID            uint `json:"id" gorm:"primaryKey;not null"`
	ProductItemID uint `jsong:"product_item_id" gorm:"not null"`
	ProductItem   ProductItem
	Image         string `json:"image" gorm:"not null"`
}
