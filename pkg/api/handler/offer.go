package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/response"
)

// AddOffer godoc
// @summary api for admin to add new offer
// @id AddOffer
// @tags Offers
// @Param input body request.Offer{} true "input field"
// @Router /admin/offers [post]
// @Success 200 {object} response.Response{} "successfully offer added"
// @Failure 400 {object} response.Response{} "invalid input"
func (p *ProductHandler) AddOffer(ctx *gin.Context) {

	var body request.Offer

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	var offer domain.Offer

	copier.Copy(&offer, &body)

	err := p.productUseCase.AddOffer(ctx, offer)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, "Failed to add offer", err, nil)

		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully offer added", nil)
}

func (c *ProductHandler) RemoveOffer(ctx *gin.Context) {

	offerID, err := request.GetParamAsUint(ctx, "offer_id")
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindParamFailMessage, err, nil)
		return
	}

	err = c.productUseCase.RemoveOffer(ctx, offerID)

	if err != nil {
		response.ErrorResponse(ctx, 400, "Failed to remove offer", err, nil)
		return
	}

	response.SuccessResponse(ctx, 200, "successfully offer removed", nil)

}

// FindAllOffers godoc
// @summary api for show all offers
// @id FindAllOffers
// @tags Offers
// @Router /admin/offers/ [get]
// @Success 200 {object} response.Response{} ""successfully got all offers"
// @Failure 500 {object} response.Response{} "faild to get offers"
func (c *ProductHandler) FindAllOffers(ctx *gin.Context) {

	offers, err := c.productUseCase.FindAllOffers(ctx)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find all offers", err, nil)

		return
	}

	if offers == nil {
		response.SuccessResponse(ctx, http.StatusOK, "No offer found", offers)

		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully found all offers", offers)
}

// FindAllOfferCategories godoc
// @summary api for admin to get all offers of categories
// @id FindAllOfferCategories
// @tags Offers
// @Router /admin/offers/category [get]
// @Success 200 {object} response.Response{} "successfully got all offer_category"
// @Failure 500 {object} response.Response{} "failed to get offers_category"
func (c *ProductHandler) FindAllOfferCategories(ctx *gin.Context) {

	offerCategories, err := c.productUseCase.FindAllOffersOfCategories(ctx)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find offer categories", err, nil)
		return
	}

	if len(offerCategories) == 0 {
		response.SuccessResponse(ctx, http.StatusOK, "No offer categories found", nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully found offers categories", offerCategories)
}

// AddOfferCategory godoc
// @summary api for admin to add offer for category
// @id AddOfferCategory
// @tags Offers
// @Param input body req.OfferCategory{} true "input field"
// @Router /admin/offers/category [post]
// @Success 200 {object} response.Response{} "successfully offer added for category"
// @Failure 400 {object} response.Response{} "invalid input"
func (c *ProductHandler) AddOfferCategory(ctx *gin.Context) {

	var body request.OfferCategory
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	var offerCategory domain.OfferCategory
	copier.Copy(&offerCategory, &body)

	err := c.productUseCase.AddOfferCategory(ctx, offerCategory)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, "Failed to add offer", err, nil)

		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully offer added for given category")
}

func (c *ProductHandler) RemoveOfferCategory(ctx *gin.Context) {

	offerCategoryID, err := request.GetParamAsUint(ctx, "offer_category_id")
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindParamFailMessage, err, nil)
		return
	}

	err = c.productUseCase.RemoveOfferCategory(ctx, domain.OfferCategory{ID: offerCategoryID})

	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, "Failed to remove offer form category", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully offer removed from category")
}

// ReplaceOfferCategory godoc
// @summary api for admin to add offer for category
// @id ReplaceOfferCategory
// @tags Offers
// @Param input body request.OfferCategory{} true "input field"
// @Router /admin/offers/category/replace [post]
// @Success 200 {object} response.Response{} "successfully offer replaced for category"
// @Failure 400 {object} response.Response{} "invalid input"
func (c *ProductHandler) ReplaceOfferCategory(ctx *gin.Context) {

	var body request.OfferCategory

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	var offerCategory domain.OfferCategory
	copier.Copy(&offerCategory, &body)

	err := c.productUseCase.ReplaceOfferCategory(ctx, offerCategory)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, "Failed to replace offer for given category", err, nil)
		return
	}

	response.SuccessResponse(ctx, 200, "Successfully offer replaced for given category")
}

// GetOffersOfProducts godoc
// @summary api for admin to get all offers of products
// @id GetOffersOfProducts
// @tags Offers
// @Router /admin/offers/products [get]
// @Success 200 {object} response.Response{} "successfully got all offers_categories"
// @Failure 500 {object} response.Response{} "failed to get offer_products"
func (c *ProductHandler) FindAllOffersOfProducts(ctx *gin.Context) {

	offersOfCategories, err := c.productUseCase.FindAllOffersOfProducts(ctx)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find all offer products", err, nil)
		return
	}

	if offersOfCategories == nil {
		response.SuccessResponse(ctx, http.StatusOK, "No offer products found", nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully found all offer products", offersOfCategories)
}

// AddOfferProduct godoc
// @summary api for admin to add offer for product
// @id AddOfferProduct
// @tags Offers
// @Param input body request.OfferProduct{} true "input field"
// @Router /admin/offers/products [post]
// @Success 200 {object} response.Response{} "successfully offer added for product"
// @Failure 400 {object} response.Response{} "invalid input"
func (c *ProductHandler) AddOfferProduct(ctx *gin.Context) {

	var body request.OfferProduct

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	var offerProduct domain.OfferProduct
	copier.Copy(&offerProduct, &body)

	err := c.productUseCase.AddOfferProduct(ctx, offerProduct)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, "Failed to add offer for given product", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully offer added to given product")
}

// RemoveOfferProduct godoc
// @summary api for admin to remove offer from product
// @id RemoveOfferProduct
// @tags Offers
// @param offer_product_id path int true "offer_product_id"
// @Router /admin/offers/products/ [delete]
// @Success 200 {object} response.Response{} "successfully offer removed from product"
// @Failure 400 {object} response.Response{} "invalid input on params"
func (c *ProductHandler) RemoveOfferProduct(ctx *gin.Context) {

	offerProductID, err := request.GetParamAsUint(ctx, "offer_product_id")
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindParamFailMessage, err, nil)
		return
	}

	err = c.productUseCase.RemoveOfferProducts(ctx, domain.OfferProduct{ID: offerProductID})

	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, "Failed to remove offer form product", err, nil)
		return
	}

	response.SuccessResponse(ctx, 200, "Successfully offer removed from product")
}

// ReplaceOfferProduct godoc
// @summary api for admin to replace a new offer on an existing offer for a product
// @id ReplaceOfferProduct
// @tags Offers
// @Param input body request.OfferProduct{} true "input field"
// @Router /admin/offers/products [put]
// @Success 200 {object} response.Response{} "successfully offer replaced for product"
// @Failure 400 {object} response.Response{} "invalid input"
func (c *ProductHandler) ReplaceOfferProduct(ctx *gin.Context) {

	var body request.OfferProduct

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	var offerProduct domain.OfferProduct
	copier.Copy(&offerProduct, &body)

	err := c.productUseCase.ReplaceOfferProducts(ctx, offerProduct)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, "Failed to replace offer for given product", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully offer replaced for  given product")
}
