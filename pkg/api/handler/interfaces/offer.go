package interfaces

import "github.com/gin-gonic/gin"

type OfferHandler interface {

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
