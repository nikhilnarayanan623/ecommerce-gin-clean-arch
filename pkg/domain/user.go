package domain

type User struct {
	ID          uint   `json:"id" gorm:"primaryKey;unique"`
	UserName    string `json:"user_name" gorm:"not null" binding:"required,min=3,max=15"`
	FirstName   string `json:"first_name" gorm:"not null" binding:"required,min=2,max=50"`
	LastName    string `json:"last_name" gorm:"not null" binding:"required,min=1,max=50"`
	Age         uint   `json:"age" gorm:"not null" binding:"required,numeric"`
	Email       string `json:"email" gorm:"unique;not null" binding:"required,email"`
	Phone       string `json:"phone" gorm:"unique;not null" binding:"required,min=10,max=10"`
	Password    string `json:"password" gorm:"not null" binding:"required"`
	BlockStatus bool   `json:"block_status" gorm:"not null"`
}

// many to many join
type UserAdress struct {
	UserID    uint `json:"user_id" gorm:"not null"`
	User      User
	AddressID uint `json:"address_id" gorm:"not null"`
	Address   Address
}

type Address struct {
	ID          uint   `json:"id" gorm:"primaryKey;unique"`
	Name        string `json:"name" gorm:"not null" binding:"required,min=2,max=50"`
	PhoneNumber string `json:"phone_number" gorm:"not null" binding:"required,min=10,max=10"`
	House       string `json:"house" gorm:"not null" binding:"required"`
	Area        string `json:"area" gorm:"not null"`
	LandMark    string `json:"land_mark" gorm:"not null" binding:"required"`
	City        string `json:"city" gorm:"not null"`
	Pincode     uint   `json:"pincode" gorm:"not null" binding:"required,numeric,min=6,max=6"`
	CountryID   uint   `jsong:"country_id" gorm:"not null" binding:"required"`
	Country     Country
}

type Country struct {
	ID          uint   `json:"id" gorm:"primaryKey;unique;"`
	CountryName string `json:"country_name" gorm:"unique;not null"`
}

// Wish List
type WishList struct {
	ID            uint `json:"id" gorm:"primaryKey;not null"`
	UserID        uint `json:"user_id" gorm:"not null"`
	User          User
	ProductItemID uint `json:"product_item_id" gorm:"not null"`
	ProductItem   ProductItem
}

// Cart
type Cart struct {
	ID         uint `json:"id" gorm:"primaryKey;not null"`
	UserID     uint `json:"user_id" gorm:"not null"`
	User       User
	TotalPrice uint `json:"total_price" gorm:"not null"`
}

type CartItem struct {
	ID            uint `json:"id" gorm:"primaryKey;not null"`
	CartID        uint `josn:"cart_id" gorm:"not null"`
	Cart          Cart
	ProductItemID uint `json:"product_item_id" gorm:"not null"`
	ProductItem   ProductItem
	Qty           uint `json:"qty" gorm:"not null"`
}
