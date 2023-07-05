package interfaces

import "github.com/gin-gonic/gin"

type StockHandler interface {
	UpdateStock(ctx *gin.Context)
	GetAllStocks(ctx *gin.Context)
}
