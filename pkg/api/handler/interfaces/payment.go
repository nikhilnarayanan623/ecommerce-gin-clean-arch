package interfaces

import "github.com/gin-gonic/gin"

type PaymentHandler interface {

	// payment
	CartOrderPaymentSelectPage(ctx *gin.Context)
	// AddPaymentMethod(ctx *gin.Context)
	UpdatePaymentMethod(ctx *gin.Context)
	GetAllPaymentMethodsAdmin() func(ctx *gin.Context)
	GetAllPaymentMethodsUser() func(ctx *gin.Context)

	PaymentCOD(ctx *gin.Context)

	RazorpayCheckout(ctx *gin.Context)
	RazorpayVerify(ctx *gin.Context)

	StripePaymentVeify(ctx *gin.Context)
	StripPaymentCheckout(ctx *gin.Context)
}
