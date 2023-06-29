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

	// offer
	SaveOffer(ctx *gin.Context)
	RemoveOffer(ctx *gin.Context)
	GetAllOffers(ctx *gin.Context)

	// category offer
	GetAllCategoryOffers(ctx *gin.Context)
	SaveCategoryOffer(ctx *gin.Context)
	RemoveCategoryOffer(ctx *gin.Context)
	ChangeCategoryOffer(ctx *gin.Context)

	// product offer
	GetAllProductsOffers(ctx *gin.Context)
	SaveProductOffer(ctx *gin.Context)
	RemoveProductOffer(ctx *gin.Context)
	ChangeProductOffer(ctx *gin.Context)
}
