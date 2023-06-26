package interfaces

import "github.com/gin-gonic/gin"

type OrderHandler interface {
	FindAllOrderStatuses(ctx *gin.Context)

	//user side
	RazorpayCheckout(ctx *gin.Context)
	RazorpayVerify(ctx *gin.Context)

	PlaceOrderOnCOD(ctx *gin.Context)
	// ApproveOrderCOD(ctx *gin.Context)
	CancelOrder(ctx *gin.Context)
	SubmitReturnRequest(ctx *gin.Context)
	FindAllOrderItemsUser() func(ctx *gin.Context)
	FindUserOrder(ctx *gin.Context)

	//admin side
	FindAllShopOrders(ctx *gin.Context)
	FindAllOrderItemsAdmin() func(ctx *gin.Context)
	UpdateOrderStatus(ctx *gin.Context)
	FindAllOrderReturns(ctx *gin.Context)
	FindAllPendingReturns(ctx *gin.Context)
	UpdateReturnRequest(ctx *gin.Context)

	// wallet
	FindUserWallet(ctx *gin.Context)
	FindUserWalletTransactions(ctx *gin.Context)
}
