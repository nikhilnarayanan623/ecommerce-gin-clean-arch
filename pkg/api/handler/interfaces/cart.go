package interfaces

import "github.com/gin-gonic/gin"

type CartHandler interface {
	AddToCart(ctx *gin.Context)
	GetCart(ctx *gin.Context)
	UpdateCart(ctx *gin.Context)
	RemoveFromCart(ctx *gin.Context)
}
