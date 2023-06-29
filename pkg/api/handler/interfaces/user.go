package interfaces

import "github.com/gin-gonic/gin"

type UserHandler interface {
	GetProfile(ctx *gin.Context)
	UpdateProfile(ctx *gin.Context)

	SaveAddress(ctx *gin.Context)
	GetAllAddresses(ctx *gin.Context)
	UpdateAddress(ctx *gin.Context)
	SaveToWishList(ctx *gin.Context)
	RemoveFromWishList(ctx *gin.Context)
	GetWishList(ctx *gin.Context)
}
