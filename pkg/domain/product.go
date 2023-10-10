package domain

import "time"

// represent a model of product
type Product struct {
	ID            uint      `json:"id" gorm:"primaryKey;not null"`
	Name          string    `json:"product_name" gorm:"not null" binding:"required,min=3,max=50"`
	Description   string    `json:"description" gorm:"not null" binding:"required,min=10,max=100"`
	CategoryID    uint      `json:"category_id" binding:"omitempty,numeric"`
	Category      Category  `json:"-"`
	BrandID       uint      `gorm:"not null"`
	Brand         Brand     `json:"-"`
	Price         uint      `json:"price" gorm:"not null" binding:"required,numeric"`
	DiscountPrice uint      `json:"discount_price"`
	Image         string    `json:"image" gorm:"not null"`
	CreatedAt     time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// this for a specific variant of product
type ProductItem struct {
	ID            uint `json:"id" gorm:"primaryKey;not null"`
	ProductID     uint `json:"product_id" gorm:"not null" binding:"required,numeric"`
	Product       Product
	QtyInStock    uint      `json:"qty_in_stock" gorm:"not null" binding:"required,numeric"`
	Price         uint      `json:"price" gorm:"not null" binding:"required,numeric"`
	SKU           string    `json:"sku" gorm:"unique;not null"`
	DiscountPrice uint      `json:"discount_price"`
	CreatedAt     time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// for a products category main and sub category as self joining
type Category struct {
	ID         uint      `json:"-" gorm:"primaryKey;not null"`
	CategoryID uint      `json:"category_id"`
	Category   *Category `json:"-"`
	Name       string    `json:"category_name" gorm:"not null" binding:"required,min=1,max=30"`
}

type Brand struct {
	ID   uint   `json:"id" gorm:"primaryKey;not null"`
	Name string `json:"brand_name" gorm:"unique;not null"`
}

// variation means size color etc..
type Variation struct {
	ID         uint     `json:"-" gorm:"primaryKey;not null"`
	CategoryID uint     `json:"category_id" gorm:"not null" binding:"required,numeric"`
	Category   Category `json:"-"`
	Name       string   `json:"variation_name" gorm:"not null" binding:"required"`
}

// variation option means values are like s,m,xl for size and blue,white,black for Color
type VariationOption struct {
	ID          uint      `json:"-" gorm:"primaryKey;not null"`
	VariationID uint      `json:"variation_id" gorm:"not null" binding:"required,numeric"` // a specific field of variation like color/size
	Variation   Variation `json:"-"`
	Value       string    `json:"variation_value" gorm:"not null" binding:"required"` // the variations value like blue/XL
}

type ProductConfiguration struct {
	ProductItemID     uint            `json:"product_item_id" gorm:"not null"`
	ProductItem       ProductItem     `json:"-"`
	VariationOptionID uint            `json:"variation_option_id" gorm:"not null"`
	VariationOption   VariationOption `json:"-"`
}

// to store a url of productItem Id along a unique url
// so we can ote multiple images url for a ProductItem
// one to many connection
type ProductImage struct {
	ID            uint        `json:"id" gorm:"primaryKey;not null"`
	ProductItemID uint        `json:"product_item_id" gorm:"not null"`
	ProductItem   ProductItem `json:"-"`
	Image         string      `json:"image" gorm:"not null"`
}

// offer
type Offer struct {
	ID           uint      `json:"id" gorm:"primaryKey;not null" swaggerignore:"true"`
	Name         string    `json:"offer_name" gorm:"not null;unique" binding:"required"`
	Description  string    `json:"description" gorm:"not null" binding:"required,min=6,max=50"`
	DiscountRate uint      `json:"discount_rate" gorm:"not null" binding:"required,numeric,min=1,max=100"`
	StartDate    time.Time `json:"start_date" gorm:"not null" binding:"required"`
	EndDate      time.Time `json:"end_date" gorm:"not null" binding:"required,gtfield=StartDate"`
}

type OfferCategory struct {
	ID         uint     `json:"id" gorm:"primaryKey;not null"`
	OfferID    uint     `json:"offer_id" gorm:"not null"`
	Offer      Offer    `json:"-"`
	CategoryID uint     `json:"category_id" gorm:"not null"`
	Category   Category `json:"-"`
}

type OfferProduct struct {
	ID        uint `json:"id" gorm:"primaryKey;not null"`
	OfferID   uint `json:"offer_id" gorm:"not null"`
	Offer     Offer
	ProductID uint `json:"product_id" gorm:"not null"`
	Product   Product
}
