package interfaces

import "github.com/gin-gonic/gin"

type ProductHandler interface {
	GetAllCategories(ctx *gin.Context)
	SaveCategory(ctx *gin.Context)
	SaveSubCategory(ctx *gin.Context)
	SaveVariation(ctx *gin.Context)
	SaveVariationOption(ctx *gin.Context)
	GetAllVariations(ctx *gin.Context)

	GetAllProductsAdmin() func(ctx *gin.Context)
	GetAllProductsUser() func(ctx *gin.Context)

	SaveProduct(ctx *gin.Context)
	UpdateProduct(ctx *gin.Context)

	SaveProductItem(ctx *gin.Context)
	GetAllProductItemsAdmin() func(ctx *gin.Context)
	GetAllProductItemsUser() func(ctx *gin.Context)
}
