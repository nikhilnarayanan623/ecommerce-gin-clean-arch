package interfaces

import "github.com/gin-gonic/gin"

type ProductHandler interface {
	GetAlllCategories(ctx *gin.Context)
	AddCategory(ctx *gin.Context)
	AddVariation(ctx *gin.Context)
	AddVariationOption(ctx *gin.Context)

	GetAllProducts(ctx *gin.Context)
	AddProducts(ctx *gin.Context)
	UpdateProduct(ctx *gin.Context)

	AddProductItem(ctx *gin.Context)
	GetProductItems(ctx *gin.Context)
}
