package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/response"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase"
	usecaseInterface "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
)

type brandHandler struct {
	brandUseCase usecaseInterface.BrandUseCase
}

func NewBrandHandler(brandUseCase usecaseInterface.BrandUseCase) interfaces.BrandHandler {
	return &brandHandler{
		brandUseCase: brandUseCase,
	}
}

// @Summary		Save Brand
// @Description	API for admin to save new brand
// @Security		BearerAuth
// @Tags			Admin Brand
// @Id				SaveBrand
// @Param			inputs	body	request.Brand{}	true	"Input Field"
// @Router			/admin/brands [post]
// @Success		200	{object}	response.Response{domain.Brand{}}	"successfully brand created"
// @Failure		400	{object}	response.Response{}	"invalid input"
// @Failure		409	{object}	response.Response{}	"brand name already exist"
// @Failure		500	{object}	response.Response{}	"failed to create brand"
func (b *brandHandler) Save(ctx *gin.Context) {

	var body request.Brand

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, body)
		return
	}

	brand := domain.Brand{
		Name: body.Name,
	}

	brand, err := b.brandUseCase.Save(brand)

	if err != nil {
		var (
			statusCode = http.StatusInternalServerError
			message    = "failed to save brand"
		)
		if err == usecase.ErrBrandAlreadyExist {
			statusCode = http.StatusConflict
			message = "brand name already exist different other name"
		}
		response.ErrorResponse(ctx, statusCode, message, err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusCreated, "successfully brand created", brand)
}

// @Summary		Find One Brand
// @Description	API for admin to find one brand
// @Security		BearerAuth
// @Tags			Admin Brand
// @Id				FindOneBrand
// @Param			brand_id	path	int	true	"Brand ID"
// @Router			/admin/brands/{brand_id} [get]
// @Success		200	{object}	response.Response{domain.Brand{}}	"successfully brand found"
// @Failure		400	{object}	response.Response{}	"invalid input"
// @Failure		500	{object}	response.Response{}	"failed to find brand"
func (b *brandHandler) FindOne(ctx *gin.Context) {

	brandID, err := request.GetParamAsUint(ctx, "brand_id")
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindParamFailMessage, err, nil)
		return
	}

	brand, err := b.brandUseCase.FindOne(brandID)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "failed to find brand", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "successfully found brand", brand)
}

// @Summary		Find All Brand
// @Description	API for admin to find all brands
// @Security		BearerAuth
// @Tags			Admin Brand
// @Id				FindAllBrands
// @Param			page_number	query	int	false	"Page number"
// @Param			count		query	int	false	"Count"
// @Router			/admin/brands [get]
// @Success		200	{object}	response.Response{[]domain.Brand{}}	"successfully found all brands"
// @Success		204	{object}	response.Response{[]domain.Brand{}}	"there is no brands to show"
// @Failure		500	{object}	response.Response{}	"failed to find brand"
func (b *brandHandler) FindAll(ctx *gin.Context) {

	pagination := request.GetPagination(ctx)

	brands, err := b.brandUseCase.FindAll(pagination)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "failed to find all brands", err, nil)
		return
	}

	if len(brands) == 0 {
		response.SuccessResponse(ctx, http.StatusNoContent, "there is no brands available to show")
	}

	response.SuccessResponse(ctx, http.StatusOK, "successfully found all brands", brands)
}

// @Summary		Save Brand
// @Description	API for admin to update brand
// @Security		BearerAuth
// @Tags			Admin Brand
// @Id				UpdateBrand
// @Param			brand_id	path	int	true	"Brand ID"
// @Param			inputs	body	request.Brand{}	true	"Input Field"
// @Router			/admin/brands/{brand_id} [put]
// @Success		200	{object}	response.Response{domain.Brand{}}	"successfully brand updated"
// @Failure		400	{object}	response.Response{}	"invalid input"
// @Failure		500	{object}	response.Response{}	"failed to update brand"
func (b *brandHandler) Update(ctx *gin.Context) {

	brandID, err := request.GetParamAsUint(ctx, "brand_id")
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindParamFailMessage, err, nil)
		return
	}

	var body request.Brand

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, body)
		return
	}

	brand := domain.Brand{
		ID:   brandID,
		Name: body.Name,
	}

	err = b.brandUseCase.Update(brand)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "failed to update brand", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "successfully updated brand")
}

// @Summary		Save Brand
// @Description	API for admin to delete brand
// @Security		BearerAuth
// @Tags			Admin Brand
// @Id				DeleteBrand
// @Param			brand_id	path	int	true	"Brand ID"
// @Router			/admin/brands/{brand_id} [delete]
// @Success		200	{object}	response.Response{domain.Brand{}}	"successfully brand deleted"
// @Failure		400	{object}	response.Response{}	"invalid input"
// @Failure		500	{object}	response.Response{}	"failed to delete brand"
func (b *brandHandler) Delete(ctx *gin.Context) {

	brandID, err := request.GetParamAsUint(ctx, "brand_id")
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindParamFailMessage, err, nil)
		return
	}

	err = b.brandUseCase.Delete(brandID)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "failed to deleted brand", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "successfully deleted brand")
}
