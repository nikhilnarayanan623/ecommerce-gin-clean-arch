package interfaces

import "github.com/gin-gonic/gin"

type ProductHandler interface {
	FindAllCategories(ctx *gin.Context)
	SaveCategory(ctx *gin.Context)
	SaveSubCategory(ctx *gin.Context)
	SaveVariation(ctx *gin.Context)
	SaveVariationOption(ctx *gin.Context)
	FindAllVariations(ctx *gin.Context)

	FindAllProductsAdmin() func(ctx *gin.Context)
	FindAllProductsUser() func(ctx *gin.Context)

	SaveProduct(ctx *gin.Context)
	UpdateProduct(ctx *gin.Context)

	SaveProductItem(ctx *gin.Context)
	FindAllProductItemsAdmin() func(ctx *gin.Context)
	FindAllProductItemsUser() func(ctx *gin.Context)

	// offer
	SaveOffer(ctx *gin.Context)
	RemoveOffer(ctx *gin.Context)
	FindAllOffers(ctx *gin.Context)

	// category offer
	FindAllCategoryOffers(ctx *gin.Context)
	SaveCategoryOffer(ctx *gin.Context)
	RemoveCategoryOffer(ctx *gin.Context)
	ChangeCategoryOffer(ctx *gin.Context)

	// product offer
	FindAllProductsOffers(ctx *gin.Context)
	SaveProductOffer(ctx *gin.Context)
	RemoveProductOffer(ctx *gin.Context)
	ChangeProductOffer(ctx *gin.Context)
}
