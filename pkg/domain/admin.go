package domain

type Admin struct {
	ID       uint   `json:"id" gorm:"primaryKey;not null"`
	UserName string `json:"user_name" gorm:"not null" binding:"omitempty,min=3,max=15"`
	Email    string `json:"email" gorm:"not null" binding:"omitempty,email"`
	Password string `json:"password" gorm:"not null" binding:"required,min=5,max=30"`
}
