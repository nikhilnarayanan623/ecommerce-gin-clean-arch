package interfaces

import "github.com/gin-gonic/gin"

type BrandHandler interface {
	Save(ctx *gin.Context)
	FindOne(ctx *gin.Context)
	FindAll(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}
