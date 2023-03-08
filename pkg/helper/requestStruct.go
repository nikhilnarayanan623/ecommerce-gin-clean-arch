package helper

type LoginStruct struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
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
