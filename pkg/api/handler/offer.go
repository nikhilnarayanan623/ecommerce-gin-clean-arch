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
// @Router /admin/offers/category [get]
// @Success 200 {object} res.Response{} "successfully all offers and categories got for offer category page"
// @Failure 500 {object} res.Response{} "faild to get offers"
func (c *ProductHandler) OfferCategoryPage(ctx *gin.Context) {

	resOfferCategoryData, err := c.productUseCase.OfferCategoryPage(ctx)
	if err != nil {
		response := res.ErrorResponse(500, "faild to get offer category page data", err.Error(), nil)
		ctx.JSON(500, response)
		return
	}

	if resOfferCategoryData.Offers == nil {
		response := res.SuccessResponse(200, "there is no offer so can't add offer for category", nil)
		ctx.JSON(http.StatusOK, response)
		return
	}

	response := res.SuccessResponse(200, "successfully all offers and categories got for offer category page", resOfferCategoryData)
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

	offerCategoryID, err := helper.StringToUint(ctx.Param("offer_category_id"))
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

// OfferProductsPage godoc
// @summary api for show all offers and products
// @id OfferProductsPage
// @tags Offers
// @Router /admin/offers/products [get]
// @Success 200 {object} res.Response{} "successfully all offers and categories got for offer products page"
// @Failure 500 {object} res.Response{} "faild to get offers"
func (c *ProductHandler) OfferProductsPage(ctx *gin.Context) {

	resOfferProductsData, err := c.productUseCase.OfferProductsPage(ctx)
	if err != nil {
		response := res.ErrorResponse(500, "faild to get offer products page data", err.Error(), nil)
		ctx.JSON(500, response)
		return
	}

	if resOfferProductsData.Offers == nil {
		response := res.SuccessResponse(200, "there is no offer so can't add offer for product", nil)
		ctx.JSON(http.StatusOK, response)
		return
	}

	response := res.SuccessResponse(200, "successfully all offers and categories got for offer products page", resOfferProductsData)
	ctx.JSON(http.StatusOK, response)

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

func (c *ProductHandler) RemoveOfferProduct(ctx *gin.Context) {

	offerProdctID, err := helper.StringToUint(ctx.Param("offer_product_id"))
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
// @Router /admin/offers/product [post]
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
