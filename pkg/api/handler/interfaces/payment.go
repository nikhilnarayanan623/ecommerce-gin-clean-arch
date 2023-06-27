package interfaces

import "github.com/gin-gonic/gin"

type PaymentHandler interface {

	// payment
	CartOrderPaymentSelectPage(ctx *gin.Context)
	// AddPaymentMethod(ctx *gin.Context)
	UpdatePaymentMethod(ctx *gin.Context)
	FindAllPaymentMethodsAdmin() func(ctx *gin.Context)
	FindAllPaymentMethodsUser() func(ctx *gin.Context)
}
