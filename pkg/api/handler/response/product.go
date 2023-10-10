package response

import (
	"time"
)

// response for product
type Product struct {
	ID               uint      `json:"product_id"`
	CategoryID       uint      `json:"category_id"`
	Price            uint      `json:"price"`
	DiscountPrice    uint      `json:"discount_price"`
	Name             string    `json:"product_name"`
	Description      string    `json:"description" `
	CategoryName     string    `json:"category_name"`
	MainCategoryName string    `json:"main_category_name"`
	BrandID          uint      `json:"brand_id"`
	BrandName        string    `json:"brand_name"`
	Image            string    `json:"image"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// for a specific category representation
type Category struct {
	ID          uint          `json:"category_id"`
	Name        string        `json:"category_name"`
	SubCategory []SubCategory `json:"sub_category" gorm:"-"`
}

type SubCategory struct {
	ID   uint   `json:"category_id"`
	Name string `json:"category_name"`
}

// for a specific variation representation
type Variation struct {
	ID               uint              `json:"variation_id"`
	Name             string            `json:"variation_name"`
	VariationOptions []VariationOption `gorm:"-"`
}

// for a specific variation Value representation
type VariationOption struct {
	ID    uint   `json:"variation_option_id"`
	Value string `json:"variation_value"`
}

// for response a specific products all product items
type ProductItems struct {
	ID               uint                    `json:"product_item_id"`
	Name             string                  `json:"product_name"`
	ProductID        uint                    `json:"product_id"`
	Price            uint                    `json:"price"`
	DiscountPrice    uint                    `json:"discount_price"`
	SKU              string                  `json:"sku"`
	QtyInStock       uint                    `json:"qty_in_stock"`
	CategoryName     string                  `json:"category_name"`
	MainCategoryName string                  `json:"main_category_name"`
	BrandID          uint                    `json:"brand_id"`
	BrandName        string                  `json:"brand_name"`
	VariationValues  []ProductVariationValue `json:"variation_values" gorm:"-"`
	Images           []string                `json:"images" gorm:"-"`
}

type ProductVariationValue struct {
	VariationID       uint   `json:"variation_id"`
	Name              string `json:"variation_name"`
	VariationOptionID uint   `json:"variation_option_id"`
	Value             string `json:"variation_value"`
}

// offer response
type OfferCategory struct {
	OfferCategoryID uint   `json:"offer_category_id"`
	CategoryID      uint   `json:"category_id"`
	CategoryName    string `json:"category_name"`
	DiscountRate    uint   `json:"discount_rate"`
	OfferID         uint   `json:"offer_id"`
	OfferName       string `json:"offer_name"`
}

type OfferProduct struct {
	OfferProductID uint   `json:"offer_product_id"`
	ProductID      uint   `json:"product_id"`
	ProductName    string `json:"product_name"`
	DiscountRate   uint   `json:"discount_rate"`
	OfferID        uint   `json:"offer_id"`
	OfferName      string `json:"offer_name"`
}
