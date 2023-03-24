package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/res"
)

// AddOffer godoc
// @summary api for admin to add new offer
// @id AddOffer
// @tags Offers
// @Param input body domain.Offer{} true "input field"
// @Router /admin/offers [post]
// @Success 200 {object} res.Response{} "successfully offer added"
// @Failure 400 {object} res.Response{} "invalid input"
func (p *ProductHandler) AddOffer(ctx *gin.Context) {

	var body domain.Offer

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response := res.ErrorResponse(400, "invalid input", err.Error(), nil)
		ctx.JSON(400, response)
		return
	}
	fmt.Println(body.StartDate)
	err := p.productUseCase.AddOffer(ctx, body)
	if err != nil {
		response := res.ErrorResponse(400, "faild to add offer", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := res.SuccessResponse(200, "successfully offer added", nil)
	ctx.JSON(http.StatusOK, response)
}

func (c *ProductHandler) RemoveOffer(ctx *gin.Context) {

	offerID, err := helper.StringToUint(ctx.Param("offer_id"))
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
// @Param input body domain.Offer{} true "input field"
// @Router /admin/offers/ [get]
// @Success 200 {object} res.Response{} ""successfully got all offer"
// @Failure 500 {object} res.Response{} "faild to get offers"
func (c *ProductHandler) ShowAllOffers(ctx *gin.Context) {

	resOfer, err := c.productUseCase.GetAllOffers(ctx)
	if err != nil {
		response := res.ErrorResponse(500, "faild to get offers", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if resOfer.Offers == nil {
		response := res.SuccessResponse(200, "there is no offers to show", resOfer)
		ctx.JSON(http.StatusOK, response)
		return
	}

	response := res.SuccessResponse(200, "successfully got all offer", resOfer)
	ctx.JSON(http.StatusOK, response)
}

// OfferCategoryPage godoc
// @summary api for show all offers
// @id OfferCategoryPage
// @tags Offers
// @Param input body domain.Offer{} true "input field"
// @Router /admin/offers/category [get]
// @Success 200 {object} res.Response{} ""successfully all offers and categories got for offer category page"
// @Failure 500 {object} res.Response{} "faild to get offers"
func (c *ProductHandler) OfferCategoryPage(ctx *gin.Context) {

	resOfferCategoryPage, err := c.productUseCase.OfferCategoryPage(ctx)
	if err != nil {
		response := res.ErrorResponse(500, "faild to get offer category page data", err.Error(), nil)
		ctx.JSON(500, response)
		return
	}

	if resOfferCategoryPage.Offers == nil {
		response := res.SuccessResponse(200, "there is no offer so can't add offer for category", nil)
		ctx.JSON(http.StatusOK, response)
		return
	}

	response := res.SuccessResponse(200, "successfully all offers and categories got for offer category page", resOfferCategoryPage)
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

// AddOfferProduct godoc
// @summary api for admin to add offer for product
// @id AddOfferProduct
// @tags Offers
// @Param input body req.ReqOfferProduct{} true "input field"
// @Router /admin/offers/product [post]
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
