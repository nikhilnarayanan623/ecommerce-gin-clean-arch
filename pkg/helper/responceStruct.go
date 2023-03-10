package helper

var ResoposeMap map[string]string

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

// admin side
type RespCategory struct {
	ID               uint   `json:"id"`
	CategoryName     string `json:"category_name"`
	CategoryID       uint   `json:"category_id"`
	MainCategoryName string `json:"main_category_name"`
}

type ResponseProduct struct {
	ProductName  string `json:"product_name" gorm:"not null" validate:"required,min=5,max=50"`
	Description  string `json:"description" gorm:"not null" validate:"required,min=10,max=100"`
	CategoryName string `json:"category_name"`
	Price        uint   `json:"price" gorm:"not null" validate:"required,numeric"`
	Image        string `json:"image" gorm:"not null"`
}

// admin
type ResAdminLogin struct {
	ID       uint   `json:"id" `
	UserName string `json:"user_name"`
	Email    string `json:"email"`
}
