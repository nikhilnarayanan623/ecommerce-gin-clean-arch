package req

import "time"

// offer
type ReqOffer struct {
	OfferName    string    `json:"offer_name" binding:"required"`
	Description  string    `json:"description" binding:"required,min=6,max=50"`
	DiscountRate uint      `json:"discount_rate" binding:"required,numeric,min=1,max=100"`
	StartDate    time.Time `json:"start_date" binding:"required"`
	EndDate      time.Time `json:"end_date" binding:"required,gtfield=StartDate"`
}
type ReqOfferCategory struct {
	OfferID    uint `json:"offer_id" binding:"required"`
	CategoryID uint `json:"category_id" binding:"required"`
}

type ReqOfferProduct struct {
	OfferID   uint `json:"offer_id" binding:"required"`
	ProductID uint `json:"product_id" binding:"required"`
}
