package interfaces

import "github.com/gin-gonic/gin"

type AuthHandler interface {
	UserLogin(ctx *gin.Context)
	UserLoginOtpVerify(ctx *gin.Context)
	UserLoginOtpSend(ctx *gin.Context)
	UserRenewRefreshToken() gin.HandlerFunc
}
