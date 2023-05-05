package interfaces

import "github.com/gin-gonic/gin"

type UserHandler interface {
	Home(ctx *gin.Context)
}
