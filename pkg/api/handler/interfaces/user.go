package interfaces

import "github.com/gin-gonic/gin"

type UserHandler interface {
	Home(ctx *gin.Context)
	FindProfile(ctx *gin.Context)
	UpdateProfile(ctx *gin.Context)

	SaveAddress(ctx *gin.Context)
	FindAllAddresses(ctx *gin.Context)
	UpdateAddress(ctx *gin.Context)
	SaveToWishList(ctx *gin.Context)
	RemoveFromWishList(ctx *gin.Context)
	FindWishList(ctx *gin.Context)
}
