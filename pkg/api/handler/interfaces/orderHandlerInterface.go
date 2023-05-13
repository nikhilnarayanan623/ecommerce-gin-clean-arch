package interfaces

import "github.com/gin-gonic/gin"

type OrderHandler interface {
	GetAllOrderStatuses(ctx *gin.Context)

	//user side
	CartOrderPayementSelectPage(ctx *gin.Context)
	PlaceOrderCartCOD(ctx *gin.Context)
	CancellOrder(ctx *gin.Context)
	SubmitReturnRequest(ctx *gin.Context)

	GetUserOrder(ctx *gin.Context)
	GetOrderItemsByShopOrderItems(ctx *gin.Context)

	//admin side
	GetAllShopOrders(ctx *gin.Context)
	UdateOrderStatus(ctx *gin.Context)
	GetAllOrderReturns(ctx *gin.Context)
	GetAllPendingReturns(ctx *gin.Context)
	UpdateReturnRequest(ctx *gin.Context)
}
