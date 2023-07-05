package interfaces

import "github.com/gin-gonic/gin"

type AdminHandler interface {
	GetAllUsers(ctx *gin.Context)
	BlockUser(ctx *gin.Context)

	AdminSignUp(ctx *gin.Context)
	GetFullSalesReport(ctx *gin.Context)
}
