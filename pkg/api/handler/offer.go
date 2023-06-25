package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/response"
)

// SaveOffer godoc
// @summary api for admin to add new offer
// @id SaveOffer
// @tags Admin Offers
// @Param input body request.Offer{} true "input field"
// @Router /admin/offers [post]
// @Success 200 {object} response.Response{} "Successfully offer added"
// @Failure 409 {object} response.Response{} "Offer already exist"
// @Failure 400 {object} response.Response{} "Invalid inputs"
func (p *ProductHandler) SaveOffer(ctx *gin.Context) {

	var body request.Offer

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	err := p.productUseCase.SaveOffer(ctx, body)
	if err != nil {
		var statusCode int

		switch true {
		case errors.Is(err, usecase.ErrOfferNameAlreadyExist):
			statusCode = http.StatusConflict
		case errors.Is(err, usecase.ErrInvalidOfferEndDate):
			statusCode = http.StatusBadRequest
		default:
			statusCode = http.StatusInternalServerError
		}
		response.ErrorResponse(ctx, statusCode, "Failed to add offer", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully offer added", nil)
}

// FindAllOffers godoc
// @summary api for show all offers
// @id FindAllOffers
// @tags Admin Offers
// @Param page_number query int false "Page Number"
// @Param count query int false "Count"
// @Router /admin/offers [get]
// @Success 200 {object} response.Response{} ""Successfully found all offers"
// @Failure 500 {object} response.Response{} "Failed to find all offers"
func (c *ProductHandler) FindAllOffers(ctx *gin.Context) {

	pagination := request.GetPagination(ctx)

	offers, err := c.productUseCase.FindAllOffers(ctx, pagination)
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

// RemoveOffer godoc
// @summary api for admin to delete offer
// @id RemoveOffer
// @tags Admin Offers
// @Param offer_id path  int true "Offer ID"
// @Router /admin/offers/{offer_id} [delete]
// @Success 200 {object} response.Response{} "successfully offer added"
// @Failure 400 {object} response.Response{} "invalid input"
func (c *ProductHandler) RemoveOffer(ctx *gin.Context) {

	offerID, err := request.GetParamAsUint(ctx, "offer_id")
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindParamFailMessage, err, nil)
		return
	}

	err = c.productUseCase.RemoveOffer(ctx, offerID)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, "Failed to remove offer", err, nil)
		return
	}

	response.SuccessResponse(ctx, 200, "successfully offer removed", nil)

}

// @summary api for admin to add offer for category
// @id SaveCategoryOffer
// @tags Admin Offers
// @Param input body request.OfferCategory{} true "input field"
// @Router /admin/offers/category [post]
// @Success 200 {object} response.Response{} "successfully offer added for category"
// @Failure 400 {object} response.Response{} "invalid input"
func (c *ProductHandler) SaveCategoryOffer(ctx *gin.Context) {

	var body request.OfferCategory

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	err := c.productUseCase.SaveCategoryOffer(ctx, body)
	if err != nil {
		var statusCode int
		switch true {
		case errors.Is(err, usecase.ErrOfferAlreadyEnded):
			statusCode = http.StatusBadRequest
		case errors.Is(err, usecase.ErrCategoryOfferAlreadyExist):
			statusCode = http.StatusConflict
		default:
			statusCode = http.StatusInternalServerError
		}
		response.ErrorResponse(ctx, statusCode, "Failed to add offer", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully offer added for given category")
}

// FindAllCategoryOffers godoc
// @summary api for admin to get all offers of categories
// @id FindAllCategoryOffers
// @tags Admin Offers
// @Param page_number query int false "Page Number"
// @Param count query int false "Count"
// @Router /admin/offers/category [get]
// @Success 200 {object} response.Response{} "successfully got all offer_category"
// @Failure 500 {object} response.Response{} "failed to get offers_category"
func (c *ProductHandler) FindAllCategoryOffers(ctx *gin.Context) {

	pagination := request.GetPagination(ctx)

	offerCategories, err := c.productUseCase.FindAllCategoryOffers(ctx, pagination)

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

// RemoveCategoryOffer godoc
// @summary api for admin to remove offer from a category
// @id RemoveCategoryOffer
// @tags Admin Offers
// @Param offer_category_id path  int true "Offer Category ID"
// @Router /admin/offers/category/{offer_category_id} [delete]
// @Success 200 {object} response.Response{} "successfully offer added for category"
// @Failure 400 {object} response.Response{} "invalid input"
func (c *ProductHandler) RemoveCategoryOffer(ctx *gin.Context) {

	offerCategoryID, err := request.GetParamAsUint(ctx, "offer_category_id")
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindParamFailMessage, err, nil)
		return
	}

	err = c.productUseCase.RemoveCategoryOffer(ctx, offerCategoryID)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, "Failed to remove offer form category", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully offer removed from category")
}

// ReplaceCategoryOffer godoc
// @summary api for admin to replace offer of a category to another
// @id ReplaceCategoryOffer
// @tags Admin Offers
// @Param input body request.OfferCategory{} true "input field"
// @Router /admin/offers/category/replace [post]
// @Success 200 {object} response.Response{} "successfully offer replaced for category"
// @Failure 400 {object} response.Response{} "invalid input"
func (c *ProductHandler) ReplaceCategoryOffer(ctx *gin.Context) {

	var body request.OfferCategory

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	var offerCategory domain.OfferCategory
	copier.Copy(&offerCategory, &body)

	err := c.productUseCase.ReplaceCategoryOffer(ctx, offerCategory)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, "Failed to replace offer for given category", err, nil)
		return
	}

	response.SuccessResponse(ctx, 200, "Successfully offer replaced for given category")
}

// FindAllProductsOffers godoc
// @summary api for admin to get all offers of products
// @id FindAllProductsOffers
// @tags Offers
// @Param page_number query int false "Page Number"
// @Param count query int false "Count"
// @Router /admin/offers/products [get]
// @Success 200 {object} response.Response{} "successfully got all offers_categories"
// @Failure 500 {object} response.Response{} "failed to get offer_products"
func (c *ProductHandler) FindAllProductsOffers(ctx *gin.Context) {

	pagination := request.GetPagination(ctx)

	offersOfCategories, err := c.productUseCase.FindAllProductOffers(ctx, pagination)
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

// SaveProductOffer godoc
// @summary api for admin to add offer for product
// @id SaveProductOffer
// @tags Offers
// @Param input body request.OfferProduct{} true "input field"
// @Router /admin/offers/products [post]
// @Success 200 {object} response.Response{} "successfully offer added for product"
// @Failure 400 {object} response.Response{} "invalid input"
func (c *ProductHandler) SaveProductOffer(ctx *gin.Context) {

	var body request.OfferProduct

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	var offerProduct domain.OfferProduct
	copier.Copy(&offerProduct, &body)

	err := c.productUseCase.SaveProductOffer(ctx, offerProduct)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, "Failed to add offer for given product", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully offer added to given product")
}

// RemoveProductOffer godoc
// @summary api for admin to remove offer from product
// @id RemoveProductOffer
// @tags Offers
// @param offer_product_id path int true "offer_product_id"
// @Router /admin/offers/products/{offer_product_id} [delete]
// @Success 200 {object} response.Response{} "successfully offer removed from product"
// @Failure 400 {object} response.Response{} "invalid input on params"
func (c *ProductHandler) RemoveProductOffer(ctx *gin.Context) {

	offerProductID, err := request.GetParamAsUint(ctx, "offer_product_id")
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindParamFailMessage, err, nil)
		return
	}

	err = c.productUseCase.RemoveProductOffer(ctx, domain.OfferProduct{ID: offerProductID})

	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, "Failed to remove offer form product", err, nil)
		return
	}

	response.SuccessResponse(ctx, 200, "Successfully offer removed from product")
}

// ReplaceProductOffer godoc
// @summary api for admin to replace a new offer on an existing offer for a product
// @id ReplaceProductOffer
// @tags Offers
// @Param input body request.OfferProduct{} true "input field"
// @Router /admin/offers/products [put]
// @Success 200 {object} response.Response{} "successfully offer replaced for product"
// @Failure 400 {object} response.Response{} "invalid input"
func (c *ProductHandler) ReplaceProductOffer(ctx *gin.Context) {

	var body request.OfferProduct

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	var offerProduct domain.OfferProduct
	copier.Copy(&offerProduct, &body)

	err := c.productUseCase.ReplaceProductOffer(ctx, offerProduct)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, "Failed to replace offer for given product", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully offer replaced for  given product")
}
