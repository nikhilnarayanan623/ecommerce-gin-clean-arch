package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/response"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase"
	usecaseInterface "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
)

type offerHandler struct {
	offerUseCase usecaseInterface.OfferUseCase
}

func NewOfferHandler(offerUseCase usecaseInterface.OfferUseCase) interfaces.OfferHandler {
	return &offerHandler{
		offerUseCase: offerUseCase,
	}
}

// SaveOffer godoc
//	@Summary		Add offer (Admin)
//	@Security		BearerAuth
//	@Description	API for admin to add an offer (Admin)
//	@Id				SaveOffer
//	@Tags			Admin Offers
//	@Param			input	body	request.Offer{}	true	"input field"
//	@Router			/admin/offers [post]
//	@Success		200	{object}	response.Response{}	"Successfully offer added"
//	@Failure		409	{object}	response.Response{}	"Offer already exist"
//	@Failure		400	{object}	response.Response{}	"Invalid inputs"
func (p *offerHandler) SaveOffer(ctx *gin.Context) {

	var body request.Offer

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	err := p.offerUseCase.SaveOffer(ctx, body)
	if err != nil {
		var statusCode int

		switch {
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

// GetAllOffers godoc
//	@Summary		Get all offers (Admin)
//	@Security		BearerAuth
//	@Description	API for admin to get all offers
//	@Id				GetAllOffers
//	@Tags			Admin Offers
//	@Param			page_number	query	int	false	"Page Number"
//	@Param			count		query	int	false	"Count"
//	@Router			/admin/offers [get]
//	@Success		200	{object}	response.Response{}	""Successfully	found	all	offers"
//	@Failure		500	{object}	response.Response{}	"Failed to get all offers"
func (c *offerHandler) GetAllOffers(ctx *gin.Context) {

	pagination := request.GetPagination(ctx)

	offers, err := c.offerUseCase.FindAllOffers(ctx, pagination)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get all offers", err, nil)

		return
	}

	if offers == nil {
		response.SuccessResponse(ctx, http.StatusOK, "No offer found", offers)

		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully found all offers", offers)
}

// RemoveOffer godoc
//	@summary		Remove offer (Admin)
//	@Security		BearerAuth
//	@Description	API admin to remove an offer
//	@Id				RemoveOffer
//	@Tags			Admin Offers
//	@Param			offer_id	path	int	true	"Offer ID"
//	@Router			/admin/offers/{offer_id} [delete]
//	@Success		200	{object}	response.Response{}	"successfully offer added"
//	@Failure		400	{object}	response.Response{}	"invalid input"
func (c *offerHandler) RemoveOffer(ctx *gin.Context) {

	offerID, err := request.GetParamAsUint(ctx, "offer_id")
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindParamFailMessage, err, nil)
		return
	}

	err = c.offerUseCase.RemoveOffer(ctx, offerID)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, "Failed to remove offer", err, nil)
		return
	}

	response.SuccessResponse(ctx, 200, "successfully offer removed", nil)

}

//	@Summary		Add category offer (Admin)
//	@Security		BearerAuth
//	@Description	API for admin to add an offer category
//	@Id				SaveCategoryOffer
//	@Tags			Admin Offers
//	@Param			input	body	request.OfferCategory{}	true	"input field"
//	@Router			/admin/offers/category [post]
//	@Success		200	{object}	response.Response{}	"successfully offer added for category"
//	@Failure		400	{object}	response.Response{}	"invalid input"
func (c *offerHandler) SaveCategoryOffer(ctx *gin.Context) {

	var body request.OfferCategory

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	err := c.offerUseCase.SaveCategoryOffer(ctx, body)
	if err != nil {
		var statusCode int
		switch {
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

// GetAllCategoryOffers godoc
//	@Summary		Get all category offers (Admin)
//	@Security		BearerAuth
//	@Description	API for admin to get all category offers
//	@Id				GetAllCategoryOffers
//	@Tags			Admin Offers
//	@Param			page_number	query	int	false	"Page Number"
//	@Param			count		query	int	false	"Count"
//	@Router			/admin/offers/category [get]
//	@Success		200	{object}	response.Response{}	"successfully got all offer_category"
//	@Failure		500	{object}	response.Response{}	"failed to get offers_category"
func (c *offerHandler) GetAllCategoryOffers(ctx *gin.Context) {

	pagination := request.GetPagination(ctx)

	offerCategories, err := c.offerUseCase.FindAllCategoryOffers(ctx, pagination)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get offer categories", err, nil)
		return
	}

	if len(offerCategories) == 0 {
		response.SuccessResponse(ctx, http.StatusOK, "No offer categories found", nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully found offers categories", offerCategories)
}

// RemoveCategoryOffer godoc
//	@Summary		Remove category offer (Admin)
//	@Security		BearerAuth
//	@Description	API admin to remove a offer from category
//	@Id				RemoveCategoryOffer
//	@Tags			Admin Offers
//	@Param			offer_category_id	path	int	true	"Offer Category ID"
//	@Router			/admin/offers/category/{offer_category_id} [delete]
//	@Success		200	{object}	response.Response{}	"successfully offer added for category"
//	@Failure		400	{object}	response.Response{}	"invalid input"
func (c *offerHandler) RemoveCategoryOffer(ctx *gin.Context) {

	offerCategoryID, err := request.GetParamAsUint(ctx, "offer_category_id")
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindParamFailMessage, err, nil)
		return
	}

	err = c.offerUseCase.RemoveCategoryOffer(ctx, offerCategoryID)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, "Failed to remove offer form category", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully offer removed from category")
}

// ChangeCategoryOffer godoc
//	@Summary		Change product offer (Admin)
//	@Security		BearerAuth
//	@Description	API admin to change product offer to another offer
//	@Id				ChangeCategoryOffer
//	@Tags			Admin Offers
//	@Param			input	body	request.UpdateCategoryOffer{}	true	"input field"
//	@Router			/admin/offers/category [patch]
//	@Success		200	{object}	response.Response{}	"successfully offer replaced for category"
//	@Failure		400	{object}	response.Response{}	"invalid input"
func (c *offerHandler) ChangeCategoryOffer(ctx *gin.Context) {

	var body request.UpdateCategoryOffer

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	err := c.offerUseCase.ChangeCategoryOffer(ctx, body.CategoryOfferID, body.OfferID)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, "Failed to change offer for given category offer", err, nil)
		return
	}

	response.SuccessResponse(ctx, 200, "Successfully offer changed for given category offer")
}

// SaveProductOffer godoc
//	@Summary		Add product offer (Admin)
//	@Security		BearerAuth
//	@Description	API for admin to add an offer for product
//	@Id				SaveProductOffer
//	@Tags			Admin Offers
//	@Param			input	body	request.OfferProduct{}	true	"input field"
//	@Router			/admin/offers/products [post]
//	@Success		200	{object}	response.Response{}	"successfully offer added for product"
//	@Failure		400	{object}	response.Response{}	"invalid input"
func (c *offerHandler) SaveProductOffer(ctx *gin.Context) {

	var body request.OfferProduct

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	var offerProduct domain.OfferProduct
	copier.Copy(&offerProduct, &body)

	err := c.offerUseCase.SaveProductOffer(ctx, offerProduct)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, "Failed to add offer for given product", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully offer added to given product")
}

// GetAllProductsOffers godoc
//	@Summary		Get all product offers (Admin)
//	@Security		BearerAuth
//	@Description	API for admin to get all product offers
//	@Id				GetAllProductsOffers
//	@Tags			Admin Offers
//	@Param			page_number	query	int	false	"Page Number"
//	@Param			count		query	int	false	"Count"
//	@Router			/admin/offers/products [get]
//	@Success		200	{object}	response.Response{}	"successfully got all offers_categories"
//	@Failure		500	{object}	response.Response{}	"failed to get offer_products"
func (c *offerHandler) GetAllProductsOffers(ctx *gin.Context) {

	pagination := request.GetPagination(ctx)

	offersOfCategories, err := c.offerUseCase.FindAllProductOffers(ctx, pagination)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get all offer products", err, nil)
		return
	}

	if offersOfCategories == nil {
		response.SuccessResponse(ctx, http.StatusOK, "No offer products found", nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully found all offer products", offersOfCategories)
}

// RemoveProductOffer godoc
//	@Summary		Remove product offer (Admin)
//	@Security		BearerAuth
//	@Description	API admin to remove a offer from product
//	@Id				RemoveProductOffer
//	@Tags			Admin Offers
//	@param			offer_product_id	path	int	true	"offer_product_id"
//	@Router			/admin/offers/products/{offer_product_id} [delete]
//	@Success		200	{object}	response.Response{}	"Successfully offer removed from product"
//	@Failure		400	{object}	response.Response{}	"invalid input on params"
func (c *offerHandler) RemoveProductOffer(ctx *gin.Context) {

	offerProductID, err := request.GetParamAsUint(ctx, "offer_product_id")
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindParamFailMessage, err, nil)
		return
	}

	err = c.offerUseCase.RemoveProductOffer(ctx, offerProductID)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, "Failed to remove offer form product", err, nil)
		return
	}

	response.SuccessResponse(ctx, 200, "Successfully offer removed from product")
}

// ChangeProductOffer godoc
//	@Summary		Change product offer (Admin)
//	@Security		BearerAuth
//	@Description	API admin to change product offer to another offer
//	@Id				ChangeProductOffer
//	@Tags			Admin Offers
//	@Param			input	body	request.UpdateProductOffer{}	true	"input field"
//	@Router			/admin/offers/products [patch]
//	@Success		200	{object}	response.Response{}	"Successfully offer changed for  given product offer"
//	@Failure		400	{object}	response.Response{}	"invalid input"
func (c *offerHandler) ChangeProductOffer(ctx *gin.Context) {

	var body request.UpdateProductOffer

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	err := c.offerUseCase.ChangeProductOffer(ctx, body.ProductOfferID, body.OfferID)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, "Failed to change offer for given product offer", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully offer changed for  given product offer")
}
