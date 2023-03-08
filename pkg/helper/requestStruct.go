package helper

// user this for admin and user
type LoginStruct struct {
	Email    string `json:"email" validate:"login"`
	UserName string `json:"user_name" validate:"login"`
	Phone    string `json:"phone" validate:"login"`
	Password string `json:"password" validate:"required,min=3"`
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
