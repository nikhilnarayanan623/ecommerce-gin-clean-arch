package interfaces

import "github.com/gin-gonic/gin"

type AuthHandler interface {
	//userSide
	UserLogin(ctx *gin.Context)
	UserLoginOtpVerify(ctx *gin.Context)
	UserLoginOtpSend(ctx *gin.Context)
	UserRenewAccessToken() gin.HandlerFunc

	//admin side
	AdminLogin(ctx *gin.Context)
	AdminRenewAccessToken() gin.HandlerFunc
}
