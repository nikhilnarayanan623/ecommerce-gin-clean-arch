package helper

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

// response for product
type ResponseProduct struct {
	ProductName  string `json:"product_name"`
	Description  string `json:"description" `
	CategoryName string `json:"category_name"`
	Price        uint   `json:"price"`
	Image        string `json:"image"`
}

// home page response
type ResUserHome struct {
	Products []ResponseProduct `json:"products"`
	User     UserRespStrcut    `json:"user"`
}

/*
	{ResposeCartItem

		ProductItemId
		ProductName
		Qty
		SubTotal
	}
	{ResponseCart

		[]ResponseCartItem
		TotalPrice
	}

*/

type ResponseCartItem struct {
	ProductItemId uint   `json:"product_item_id"`
	ProductName   string `json:"product_name"`
	Qty           uint   `json:"qty"`
	SubTotal      uint   `json:"sub_total"`
}
type ResponseCart struct {
	CartItem   []ResponseCartItem
	TotalPrice uint `json:"total_price"`
}
