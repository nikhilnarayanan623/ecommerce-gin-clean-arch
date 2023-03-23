package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
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
