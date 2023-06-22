package interfaces

import "github.com/gin-gonic/gin"

type AdminHandler interface {
	AdminHome(ctx *gin.Context)
	FindAllUsers(ctx *gin.Context)
	BlockUser(ctx *gin.Context)

	AdminSignUp(ctx *gin.Context)

	UpdateStock(ctx *gin.Context)
	FindAllStocks(ctx *gin.Context)
	FullSalesReport(ctx *gin.Context)
}
