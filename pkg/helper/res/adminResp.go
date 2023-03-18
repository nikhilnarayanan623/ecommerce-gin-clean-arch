package res

var ResoposeMap map[string]string

// admin
type ResAdminLogin struct {
	ID       uint   `json:"id" `
	UserName string `json:"user_name"`
	Email    string `json:"email"`
}

// response of category for showing the category
type RespCategory struct {
	ID               uint   `json:"id"`
	CategoryName     string `json:"category_name"`
	CategoryID       uint   `json:"category_id"`
	MainCategoryName string `json:"main_category_name"`
}

// reponse for get all variations with its respective category
