package domain

type Admin struct {
	ID       uint   `json:"id" gorm:"primaryKey;not null"`
	UserName string `json:"user_name" gorm:"not null"`
	Email    string `json:"email" gorm:"not null" validate:"email"`
	Password string `json:"password" grom:"not null" validate:"required"`
}
