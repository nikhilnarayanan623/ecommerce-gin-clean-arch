package interfaces

import "github.com/gin-gonic/gin"

type ProductHandler interface {
	FindAllCategories(ctx *gin.Context)
	SaveCategory(ctx *gin.Context)
	SaveSubCategory(ctx *gin.Context)
	SaveVariation(ctx *gin.Context)
	SaveVariationOption(ctx *gin.Context)
	FindAllVariations(ctx *gin.Context)

	FindAllProducts(ctx *gin.Context)
	SaveProducts(ctx *gin.Context)
	UpdateProduct(ctx *gin.Context)

	SaveProductItem(ctx *gin.Context)
	FindProductItems(ctx *gin.Context)
}
