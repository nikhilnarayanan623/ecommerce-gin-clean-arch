package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
)

// AddOffer godoc
// @summary api for admin to add new offer
// @id AddOffer
// @tags Offers
// @Param input body req.ReqOffer{} true "input field"
// @Router /admin/offers [post]
// @Success 200 {object} res.Response{} "successfully offer added"
// @Failure 400 {object} res.Response{} "invalid input"
func (p *ProductHandler) AddOffer(ctx *gin.Context) {

	var body req.ReqOffer

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response := res.ErrorResponse(400, "invalid input", err.Error(), nil)
		ctx.JSON(400, response)
		return
	}

	var offer domain.Offer

	copier.Copy(&offer, &body)

	err := p.productUseCase.AddOffer(ctx, offer)
	if err != nil {
		response := res.ErrorResponse(400, "faild to add offer", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := res.SuccessResponse(200, "successfully offer added", nil)
	ctx.JSON(http.StatusOK, response)
}

func (c *ProductHandler) RemoveOffer(ctx *gin.Context) {

	offerID, err := utils.StringToUint(ctx.Param("offer_id"))
	if err != nil {
		response := res.ErrorResponse(400, "invalid input on param", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	err = c.productUseCase.RemoveOffer(ctx, offerID)
	if err != nil {
		response := res.ErrorResponse(400, "faild to rmove offer", err.Error(), nil)
		ctx.JSON(400, response)
		return
	}

	response := res.SuccessResponse(200, "successfully offer removed", nil)
	ctx.JSON(http.StatusOK, response)

}

// ShowAllOffers godoc
// @summary api for show all offers
// @id ShowAllOffers
// @tags Offers
// @Router /admin/offers/ [get]
// @Success 200 {object} res.Response{} ""successfully got all offers"
// @Failure 500 {object} res.Response{} "faild to get offers"
func (c *ProductHandler) ShowAllOffers(ctx *gin.Context) {

	offers, err := c.productUseCase.GetAllOffers(ctx)
	if err != nil {
		response := res.ErrorResponse(500, "faild to get offers", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if offers == nil {
		response := res.SuccessResponse(200, "there is no offers to show", offers)
		ctx.JSON(http.StatusOK, response)
		return
	}

	response := res.SuccessResponse(200, "successfully got all offers", offers)
	ctx.JSON(http.StatusOK, response)
}

// GetOfferCategory godoc
// @summary api for admin to get all offers of categories
// @id GetOfferCategory
// @tags Offers
// @Router /admin/offers/category [get]
// @Success 200 {object} res.Response{} "successfully got all offer_category"
// @Failure 500 {object} res.Response{} "faild to get offers_category"
func (c *ProductHandler) GetOfferCategories(ctx *gin.Context) {

	offerCategories, err := c.productUseCase.GetAllOffersOfCategories(ctx)
	if err != nil {
		response := res.ErrorResponse(500, "faild to get offer_categories", err.Error(), nil)
		ctx.JSON(500, response)
		return
	}

	if offerCategories == nil {
		response := res.SuccessResponse(200, "there is no offer_cateogies avialable", nil)
		ctx.JSON(http.StatusOK, response)
		return
	}

	response := res.SuccessResponse(200, "faild to get offers_category", offerCategories)
	ctx.JSON(http.StatusOK, response)

}

// AddOfferCategory godoc
// @summary api for admin to add offer for category
// @id AddOfferCategory
// @tags Offers
// @Param input body req.ReqOfferCategory{} true "input field"
// @Router /admin/offers/category [post]
// @Success 200 {object} res.Response{} "successfully offer added for category"
// @Failure 400 {object} res.Response{} "invalid input"
func (c *ProductHandler) AddOfferCategory(ctx *gin.Context) {

	var body req.ReqOfferCategory
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response := res.ErrorResponse(400, "invalid inputs", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	var offerCategory domain.OfferCategory
	copier.Copy(&offerCategory, &body)

	err := c.productUseCase.AddOfferCategory(ctx, offerCategory)
	if err != nil {
		response := res.ErrorResponse(400, "faild to add offer for given category", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := res.SuccessResponse(200, "successfully offer added for given category")
	ctx.JSON(200, response)
}

func (c *ProductHandler) RemoveOfferCategory(ctx *gin.Context) {

	offerCategoryID, err := utils.StringToUint(ctx.Param("offer_category_id"))
	if err != nil {
		response := res.ErrorResponse(400, "invalid input on params", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	err = c.productUseCase.RemoveOfferCategory(ctx, domain.OfferCategory{ID: offerCategoryID})

	if err != nil {
		response := res.ErrorResponse(400, "faild to remove offer form category", err.Error(), offerCategoryID)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := res.SuccessResponse(200, "successfully offer removed from cateogry")
	ctx.JSON(http.StatusOK, response)
}

// ReplaceOfferCategory godoc
// @summary api for admin to add offer for category
// @id ReplaceOfferCategory
// @tags Offers
// @Param input body req.ReqOfferCategory{} true "input field"
// @Router /admin/offers/category/replace [post]
// @Success 200 {object} res.Response{} "successfully offer replaced for category"
// @Failure 400 {object} res.Response{} "invalid input"
func (c *ProductHandler) ReplaceOfferCategory(ctx *gin.Context) {
	var body req.ReqOfferCategory
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response := res.ErrorResponse(400, "invalid inputs", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	var offerCategory domain.OfferCategory
	copier.Copy(&offerCategory, &body)

	err := c.productUseCase.ReplaceOfferCategory(ctx, offerCategory)
	if err != nil {
		response := res.ErrorResponse(400, "faild to replace offer for given category", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := res.SuccessResponse(200, "successfully offer replaced for given category")
	ctx.JSON(200, response)

}

// GetOffersOfProducts godoc
// @summary api for admin to get all offers of products
// @id GetOffersOfProducts
// @tags Offers
// @Router /admin/offers/products [get]
// @Success 200 {object} res.Response{} "successfully got all offers_categories"
// @Failure 500 {object} res.Response{} "faild to get offer_products"
func (c *ProductHandler) GetOffersOfProducts(ctx *gin.Context) {

	offersOfCategories, err := c.productUseCase.GetAllOffersOfProducts(ctx)
	if err != nil {
		response := res.ErrorResponse(500, "faild to get offer_products", err.Error(), nil)
		ctx.JSON(500, response)
		return
	}

	if offersOfCategories == nil {
		response := res.SuccessResponse(200, "there is no offer_products available", nil)
		ctx.JSON(http.StatusOK, response)
		return
	}

	response := res.SuccessResponse(200, "successfully got all offers_categories", offersOfCategories)
	ctx.JSON(http.StatusOK, response)

}

// AddOfferProduct godoc
// @summary api for admin to add offer for product
// @id AddOfferProduct
// @tags Offers
// @Param input body req.ReqOfferProduct{} true "input field"
// @Router /admin/offers/products [post]
// @Success 200 {object} res.Response{} "successfully offer added for product"
// @Failure 400 {object} res.Response{} "invalid input"
func (c *ProductHandler) AddOfferProduct(ctx *gin.Context) {

	var body req.ReqOfferProduct
	if err := ctx.Bind(&body); err != nil {
		response := res.ErrorResponse(400, "invalid inputs", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	var offerProduct domain.OfferProduct
	copier.Copy(&offerProduct, &body)

	err := c.productUseCase.AddOfferProduct(ctx, offerProduct)
	if err != nil {
		response := res.ErrorResponse(400, "faild to add offer for given product", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	response := res.SuccessResponse(200, "successfully offer added to given product")
	ctx.JSON(http.StatusOK, response)
}

// RemoveOfferProduct godoc
// @summary api for admin to remove offer from product
// @id RemoveOfferProduct
// @tags Offers
// @param offer_product_id path int true "offer_product_id"
// @Router /admin/offers/products/ [delete]
// @Success 200 {object} res.Response{} "successfully offer removed from product"
// @Failure 400 {object} res.Response{} "invalid input on params"
func (c *ProductHandler) RemoveOfferProduct(ctx *gin.Context) {

	offerProdctID, err := utils.StringToUint(ctx.Param("offer_product_id"))
	if err != nil {
		response := res.ErrorResponse(400, "invalid input on params", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	err = c.productUseCase.RemoveOfferProducts(ctx, domain.OfferProduct{ID: offerProdctID})

	if err != nil {
		response := res.ErrorResponse(400, "faild to remove offer form product", err.Error(), offerProdctID)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := res.SuccessResponse(200, "successfully offer removed from product")
	ctx.JSON(http.StatusOK, response)
}

// ReplaceOfferProduct godoc
// @summary api for admin to replace a new offer on an existing offer for a product
// @id ReplaceOfferProduct
// @tags Offers
// @Param input body req.ReqOfferProduct{} true "input field"
// @Router /admin/offers/products [put]
// @Success 200 {object} res.Response{} "successfully offer replaced for product"
// @Failure 400 {object} res.Response{} "invalid input"
func (c *ProductHandler) ReplaceOfferProduct(ctx *gin.Context) {

	var body req.ReqOfferProduct
	if err := ctx.Bind(&body); err != nil {
		response := res.ErrorResponse(400, "invalid inputs", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	var offerProduct domain.OfferProduct
	copier.Copy(&offerProduct, &body)

	err := c.productUseCase.ReplaceOfferProducts(ctx, offerProduct)
	if err != nil {
		response := res.ErrorResponse(400, "faild to replace offer for given product", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	response := res.SuccessResponse(200, "successfully offer replaced for  given product")
	ctx.JSON(http.StatusOK, response)
}
