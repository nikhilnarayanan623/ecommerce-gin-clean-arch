package interfaces

import "github.com/gin-gonic/gin"

type AdminHandler interface {
	AdminHome(ctx *gin.Context)
	ListUsers(ctx *gin.Context)
	BlockUser(ctx *gin.Context)

	UpdateStock(ctx *gin.Context)
	GetAllStockDetails(ctx *gin.Context)
	FullSalesReport(ctx *gin.Context)
}
