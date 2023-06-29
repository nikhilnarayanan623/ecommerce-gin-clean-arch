package interfaces

import "github.com/gin-gonic/gin"

type OrderHandler interface {
	GetAllOrderStatuses(ctx *gin.Context)

	SaveOrder(ctx *gin.Context)

	// ApproveOrderCOD(ctx *gin.Context)
	CancelOrder(ctx *gin.Context)
	SubmitReturnRequest(ctx *gin.Context)
	GetAllOrderItemsUser() func(ctx *gin.Context)
	GetUserOrder(ctx *gin.Context)

	//admin side
	GetAllShopOrders(ctx *gin.Context)
	GetAllOrderItemsAdmin() func(ctx *gin.Context)
	UpdateOrderStatus(ctx *gin.Context)
	GetAllOrderReturns(ctx *gin.Context)
	GetAllPendingReturns(ctx *gin.Context)
	UpdateReturnRequest(ctx *gin.Context)

	// wallet
	GetUserWallet(ctx *gin.Context)
	GetUserWalletTransactions(ctx *gin.Context)
}
