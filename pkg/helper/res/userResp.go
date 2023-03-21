package res

// user details response
type UserRespStrcut struct {
	ID          uint   `json:"id" copier:"must"`
	FirstName   string `json:"first_name" copier:"must"`
	LastName    string `json:"last_name" copier:"must"`
	Age         uint   `json:"age" copier:"must"`
	Email       string `json:"email" copier:"must"`
	Phone       string `json:"phone" copier:"must"`
	BlockStatus bool   `json:"block_status" copier:"must"`
}

// home page response
type ResUserHome struct {
	Products []ResponseProduct `json:"products"`
	User     UserRespStrcut    `json:"user"`
}

type ResponseCartItem struct {
	ProductItemId uint   `json:"product_item_id"`
	ProductName   string `json:"product_name"`
	Price         uint   `json:"price"`
	QtyInStock    uint   `json:"qty_in_stock"`
	Qty           uint   `json:"qty"`
	SubTotal      uint   `json:"sub_total"`
}

type ResponseCart struct {
	CartItems  []ResponseCartItem
	TotalPrice uint `json:"total_price"`
}

// address
type ResAddress struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	House       string `json:"house"`
	Area        string `json:"area"`
	LandMark    string `json:"land_mark"`
	City        string `json:"city"`
	Pincode     uint   `json:"pincode"`
	CountryID   uint   `json:"country_id"`
	CountryName string `json:"country_name"`

	IsDefault *bool `json:"is_default"`
}

// wish list response
type ResWishList struct {
	ProductItemID uint   `json:"product_item_id"`
	ProductName   string `json:"product_name"`
	Price         uint   `json:"price"`
	Image         string `json:"image"`
	QtyInStock    uint   `json:"qty_in_stock"`
	//OutOfStock    bool   `json:"out_of_stock"`
}
