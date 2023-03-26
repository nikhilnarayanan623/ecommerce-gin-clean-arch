package req

type ReqUpdateOrder struct {
	ShopOrderID   uint `json:"shop_order_id" binding:"required"`
	OrderStatusID uint `json:"order_status_id"`
}

// return request
type ReqReturn struct {
	ShopOrderID  uint   `json:"shop_order_id" binding:"required"`
	ReturnReason string `json:"return_reason" binding:"required,min=6,max=50"`
}

type ReqUpdatReturnReq struct {
	OrderReturnID uint   `json:"order_return_id" binding:"required"`
	OrderStatusID uint   `json:"order_status_id" binding:"required"`
	AdminComment  string `json:"admin_comment" bindin:"requied,min=6,max=50"`
}
