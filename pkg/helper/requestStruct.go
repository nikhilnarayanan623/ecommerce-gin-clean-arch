package helper

// user this for admin and user
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
	OTP string `json:"otp" binding:"required,min=4,max=8"`
	ID  uint   `json:"id" binding:"required,numeric"`
}

type CategoryStruct struct {
	CategoryID     uint   `json:"category_id"`
	CategoryName   string `json:"category_name"`
	VariationName  string `json:"variation_name"`
	VariationValue string `json:"variation_value"`
}

type BlockStruct struct {
	ID uint `json:"id" copier:"must"`
}

// admin side
type ReqCategory struct {
	CategoryName string `json:"category_name"` // new category name
	ID           uint   `json:"id"`            // id any of main category
}

type ProductRequest struct {
	ProductName string `json:"product_name" gorm:"not null" validate:"required,min=5,max=50"`
	Description string `json:"description" gorm:"not null" validate:"required,min=10,max=100"`
	CategoryID  uint   `json:"category_id"`
	Price       uint   `json:"price" gorm:"not null" validate:"required,numeric"`
	Image       string `json:"image" gorm:"not null"`
}
