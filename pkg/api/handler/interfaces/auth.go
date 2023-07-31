package interfaces

import "github.com/gin-gonic/gin"

type AuthHandler interface {
	//userSide
	UserLogin(ctx *gin.Context)
	UserSignUp(ctx *gin.Context)
	UserSignUpVerify(ctx *gin.Context)

	UserGoogleAuthInitialize(ctx *gin.Context)
	UserGoogleAuthLoginPage(ctx *gin.Context)
	UserGoogleAuthCallBack(ctx *gin.Context)

	UserLoginOtpVerify(ctx *gin.Context)
	UserLoginOtpSend(ctx *gin.Context)

	UserRenewAccessToken() gin.HandlerFunc

	//admin side
	AdminLogin(ctx *gin.Context)
	AdminRenewAccessToken() gin.HandlerFunc
}
