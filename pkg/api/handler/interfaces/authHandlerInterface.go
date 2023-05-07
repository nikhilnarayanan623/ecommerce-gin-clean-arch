package interfaces

import "github.com/gin-gonic/gin"

type AuthHandler interface {
	UserLogin(ctx *gin.Context)
	UserRenewRefreshToken() gin.HandlerFunc
}
