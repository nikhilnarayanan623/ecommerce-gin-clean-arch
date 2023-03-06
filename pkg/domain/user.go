package domain

type Users struct {
	ID          uint   `json:"id" gorm:"primaryKey;unique"`
	FirstName   string `json:"first_name" gorm:"not null" validate:"required,min=2,max=50"`
	LastName    string `json:"last_name" gorm:"not null" validate:"required,min=1,max=50"`
	Age         uint   `json:"age" gorm:"not null" validate:"required,numeric"`
	Email       string `json:"email" gorm:"unique;not null" validate:"required,email"`
	Phone       string `json:"phone" gorm:"unique;not null" validate:"required,min=10,max=10"`
	Password    string `json:"password" gorm:"not null" validate:"required"`
	BlockStatus bool   `json:"block_status" gorm:"not null"`
}

// many to many join
type UserAdress struct {
	UsersID   uint `json:"user_id" gorm:"not null"`
	Users     Users
	AddressID uint `json:"address_id" gorm:"not null"`
	Address   Address
}

type Address struct {
	ID          uint   `json:"id" gorm:"primaryKey;unique"`
	Name        string `json:"name" gorm:"not null" validate:"required,min=2,max=50"`
	PhoneNumber string `json:"phone_number" gorm:"not null" validate:"required,min=10,max=10"`
	House       string `json:"house" gorm:"not null" validate:"required"`
	Area        string `json:"area" gorm:"not null"`
	LandMark    string `json:"land_mark" gorm:"not null" validate:"required"`
	City        string `json:"city" gorm:"not null"`
	Pincode     uint   `json:"pincode" gorm:"not null" validate:"required,numeric"`
	CountryID   uint   `jsong:"country_id" gorm:"not null"`
	Country     Country
}

type Country struct {
	ID          uint   `json:"id" gorm:"primaryKey;unique;"`
	CountryName string `json:"country_name" gorm:"unique;not null"`
}

// Wish List
type WishList struct {
	ID            uint `json:"id" gorm:"primaryKey;not null"`
	UsersID       uint `json:"user_id" gorm:"not null"`
	Users         Users
	ProductItemID uint `json:"product_item_id" gorm:"not null"`
	ProductItem   ProductItem
}

// Cart
type Cart struct {
	ID      uint `json:"id" gorm:"primaryKey;not null"`
	UsersID uint `json:"user_id" gorm:"not null"`
	Users   Users
}

type CartItem struct {
	CartID        uint `josn:"cart_id" gorm:"not null"`
	Cart          Cart
	ProductItemID uint `json:"product_item_id" gorm:"not null"`
	ProductItem   ProductItem
}
