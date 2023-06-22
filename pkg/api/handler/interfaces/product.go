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
	SaveProduct(ctx *gin.Context)
	UpdateProduct(ctx *gin.Context)

	SaveProductItem(ctx *gin.Context)
	FindAllProductItems(ctx *gin.Context)
}
