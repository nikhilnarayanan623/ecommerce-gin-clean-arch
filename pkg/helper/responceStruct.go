package helper

type SingleRespStruct struct {
	Error string `json:"error"`
}

// to responce with avoid unwanted details
type UserRespStrcut struct {
	ID          uint   `json:"id" copier:"must"`
	FirstName   string `json:"first_name" copier:"must"`
	LastName    string `json:"last_name" copier:"must"`
	Age         uint   `json:"age" copier:"must"`
	Email       string `json:"email" copier:"must"`
	Phone       string `json:"phone" copier:"must"`
	BlockStatus bool   `json:"block_status" copier:"must"`
}

type ResCart struct {
	CartItems  []ResCartItem
	TotalPrice uint `json:"total_price"`
}

type ResCartItem struct {
	ProductItemID uint   `json:"product_item_id"`
	ProductName   string `jsong:"product_name"`
	Qty           uint   `json:"qty"`
	OutOfStock    bool   `json:"out_of_stock"`
}
