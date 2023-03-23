package req

type ReqUpdateOrder struct {
	ShopOrderID   uint `json:"shop_order_id" binding:"required"`
	OrderStatusID uint `json:"order_status_id"`
}
