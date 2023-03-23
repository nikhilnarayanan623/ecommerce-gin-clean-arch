package req

type ReqUser struct {
	UserName  string `json:"user_name" binding:"omitempty,min=3,max=15"`
	FirstName string `json:"first_name"  binding:"omitempty,min=2,max=50"`
	LastName  string `json:"last_name"  binding:"omitempty,min=1,max=50"`
	Age       uint   `json:"age"  binding:"omitempty,numeric"`
	Email     string `json:"email" gorm:"unique;not null" binding:"omitempty,email"`
	Phone     string `json:"phone" gorm:"unique;not null" binding:"omitempty,min=10,max=10"`
}
