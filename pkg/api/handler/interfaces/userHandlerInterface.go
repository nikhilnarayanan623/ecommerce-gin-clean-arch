package interfaces

import "github.com/gin-gonic/gin"

type UserHandler interface {
	Home(ctx *gin.Context)
	GetUserProfile(ctx *gin.Context)
	UpdateUserProfile(ctx *gin.Context)

	AddAddress(ctx *gin.Context)
	GetAllAddresses(ctx *gin.Context)
	UpdateAddress(ctx *gin.Context)
}
