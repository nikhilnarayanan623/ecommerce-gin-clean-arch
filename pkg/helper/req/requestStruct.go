package req

// login struct for user and admin
type LoginStruct struct {
	UserName string `json:"user_name" binding:"omitempty,min=3,max=15"`
	Phone    string `json:"phone" binding:"omitempty,min=10,max=10"`
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"required,min=5,max=30"`
}
type OTPLoginStruct struct {
	Email    string `json:"email" binding:"omitempty,email"`
	UserName string `json:"user_name" binding:"omitempty,min=3,max=16"`
	Phone    string `json:"phone" binding:"omitempty,min=10,max=10"`
}

type OTPVerifyStruct struct {
	OTP    string `json:"otp" binding:"required,min=4,max=8"`
	UserID uint   `json:"user_id" binding:"required,numeric"`
}

type BlockStruct struct {
	UserID uint `json:"-" binding:"required,numeric"`
}

// product side
type ReqCategory struct {
	CategoryName string `json:"category_name"` // new category name
	ID           uint   `json:"id"`            // id any of main category
}

// user side
type ReqCart struct {
	UserID        uint `json:"-"`
	ProductItemID uint `json:"product_item_id" binding:"required"`
}

type ReqCartCount struct {
	UserID        uint `json:"-"`
	ProductItemID uint `json:"product_item_id" binding:"required"`
	Count         uint `json:"count" binding:"omitempty,gte=1"`
}

// for address add address
type ReqAddress struct {
	ID          uint   `json:"-"`
	Name        string `json:"name" binding:"required,min=2,max=50"`
	PhoneNumber string `json:"phone_number" binding:"required,min=10,max=10"`
	House       string `json:"house" binding:"required"`
	Area        string `json:"area"`
	LandMark    string `json:"land_mark" binding:"required"`
	City        string `json:"city"`
	Pincode     uint   `json:"pincode" binding:"required"`
	CountryID   uint   `json:"country_id" binding:"required"`

	IsDefault *bool `json:"is_default"`
}

// for address
type ReqEditAddress struct {
	ID          uint   `json:"address_id" binding:"required"`
	Name        string `json:"name" binding:"required,min=2,max=50"`
	PhoneNumber string `json:"phone_number" binding:"required,min=10,max=10"`
	House       string `json:"house" binding:"required"`
	Area        string `json:"area"`
	LandMark    string `json:"land_mark" binding:"required"`
	City        string `json:"city"`
	Pincode     uint   `json:"pincode" binding:"required"`
	CountryID   uint   `json:"country_id" binding:"required"`

	IsDefault *bool `json:"is_default"`
}

// offer
type ReqOfferCategory struct {
	OfferID    uint `json:"offer_id" binding:"required"`
	CategoryID uint `json:"category_id" binding:"required"`
}

type ReqOfferProduct struct {
	OfferID   uint `json:"offer_id" binding:"required"`
	ProductID uint `json:"product_id" binding:"required"`
}

type ReqApplyCoupon struct {
	CouponCode string `json:"coupon_code" binding:"required"`
}
