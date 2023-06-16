package interfaces

import "github.com/gin-gonic/gin"

type PaymentHandler interface {

	// payment
	CartOrderPayementSelectPage(ctx *gin.Context)
	AddPaymentMethod(ctx *gin.Context)
	UpdatePaymentMethod(ctx *gin.Context)
	GetAllPaymentMethods(ctx *gin.Context)
}
