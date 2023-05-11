package interfaces

import "github.com/gin-gonic/gin"

type AdminHandler interface {
	AdminLogin(ctx *gin.Context)
	ListUsers(ctx *gin.Context)
}
