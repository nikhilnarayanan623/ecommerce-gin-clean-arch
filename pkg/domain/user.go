package domain

type Users struct {
	ID          uint   `json:"id" gorm:"primaryKey;unique"`
	FirstName   string `json:"first_name" gorm:"not null" validate:"required,min=2,max=50"`
	LastName    string `json:"last_name" gorm:"not null" validate:"required,min=2,max=50"`
	Age         uint   `json:"age" gorm:"not null" validate:"require,numeric"`
	Email       string `json:"email" gorm:"not nul;unique" validate:"required,email"`
	Phone       string `json:"phone" gorm:"not null;unique" validate:"required,min=10,max=10"`
	Password    string `json:"password" gorm:"not null" validate:"required"`
	BlockStatus bool   `json:"block_status" gorm:"not null"`
}

// many to many join
type UserAdress struct {
	UsersID   uint
	Users     Users
	AddressID uint
	Address   Address
}

type Address struct {
	ID          uint   `json:"id" gorm:"primaryKey;unique"`
	Name        string `json:"name" gorm:"not null" validate:"required,min=2,max=50"`
	PhoneNumber string `json:"phone_number" gorm:"not null" validate:"min=10,max=10"`
	House       string `json:"house" gorm:"not null"`
	Area        string `json:"area" gorm:"not null"`
	LandMark    string `json:"land_mark" gorm:"not null"`
	City        string `json:"city" gorm:"not null"`
	Pincode     uint   `json:"pincode" gorm:"not null" validate:"required,numeric"`
	CountryID   uint   `gorm:"not null"`
	Country     Country
}

type Country struct {
	ID          uint   `json:"id" gorm:"primaryKey;unique;"`
	CountryName string `json:"country_name" gorm:"unique;not null"`
}
