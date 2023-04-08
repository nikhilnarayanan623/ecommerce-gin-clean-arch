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
	UserID uint `json:"user_id" binding:"required,numeric"`
}

type ReqPagination struct {
	Count      uint `json:"count" binding:"required,min=1"`
	PageNumber uint `json:"page_number" binding:"required,min=1"`
}
