package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	handlerInterface "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
)

type CouponHandler struct {
	couponUseCase interfaces.CouponUseCase
}

func NewCouponHandler(couponUseCase interfaces.CouponUseCase) handlerInterface.CouponHandler {
	return &CouponHandler{couponUseCase: couponUseCase}
}

// AddCoupon godoc
// @summary api for admin to add coupon
// @security ApiKeyAuth
// @tags Admin Coupon
// @id AddCoupon
// @Param        inputs   body     req.Coupon{}   true  "Input Field"
// @Router /admin/coupons [post]
// @Success 200 {object} res.Response{} "successfully coupon added"
// @Failure 400 {object} res.Response{}  "invalid input"
func (c *CouponHandler) AddCoupon(ctx *gin.Context) {

	var body req.Coupon

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response := res.ErrorResponse(400, "invalid input", err.Error(), body)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	var coupon domain.Coupon

	copier.Copy(&coupon, &body)

	err := c.couponUseCase.AddCoupon(ctx, coupon)
	if err != nil {
		response := res.ErrorResponse(400, "faild to add coupon", err.Error(), coupon)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := res.SuccessResponse(200, "successfully coupon added", body)
	ctx.JSON(http.StatusOK, response)
}

// GetAllCoupons godoc
// @summary api for admin to see all coupons
// @security ApiKeyAuth
// @tags Admin Coupon
// @id GetAllCoupons
// @Param page_number query int false "Page Number"
// @Param count query int false "Count Of Order"
// @Router /admin/coupons [get]
// @Success 200 {object} res.Response{} "successfully go all the coupons
// @Failure 500 {object} res.Response{}  "faild to get all coupons"
func (c *CouponHandler) GetAllCoupons(ctx *gin.Context) {

	count, err1 := utils.StringToUint(ctx.Query("count"))
	pageNumber, err2 := utils.StringToUint(ctx.Query("page_number"))

	err1 = errors.Join(err1, err2)
	if err1 != nil {
		response := res.ErrorResponse(400, "invalid inputs", err1.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	pagination := req.Pagination{
		PageNumber: pageNumber,
		Count:      count,
	}

	coupons, err := c.couponUseCase.GetAllCoupons(ctx, pagination)
	if err != nil {
		response := res.ErrorResponse(500, "faild to get all coupons", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if coupons == nil {
		response := res.SuccessResponse(200, "there is no coupons available", nil)
		ctx.JSON(http.StatusOK, response)
		return
	}

	response := res.SuccessResponse(200, "successsfully got all coupons", coupons)
	ctx.JSON(http.StatusOK, response)
}

// GetAllCoupons godoc
// @summary api for user to see all coupons
// @security ApiKeyAuth
// @tags User Coupon
// @id GetAllCouponsForUser
// @Param page_number query int false "Page Number"
// @Param count query int false "Count Of Order"
// @Router /coupons [get]
// @Success 200 {object} res.Response{} "successfully go all the coupons
// @Failure 500 {object} res.Response{}  "faild to get all coupons"
func (c *CouponHandler) GetAllCouponsForUser(ctx *gin.Context) {

	userID := utils.GetUserIdFromContext(ctx)
	count, err1 := utils.StringToUint(ctx.Query("count"))
	pageNumber, err2 := utils.StringToUint(ctx.Query("page_number"))

	err1 = errors.Join(err1, err2)
	if err1 != nil {
		response := res.ErrorResponse(400, "invalid inputs", err1.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	pagination := req.Pagination{
		PageNumber: pageNumber,
		Count:      count,
	}

	coupons, err := c.couponUseCase.GetCouponsForUser(ctx, userID, pagination)

	if err != nil {
		response := res.ErrorResponse(400, "faild to get copons for user", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	if coupons == nil {
		response := res.SuccessResponse(200, "there is no copons for user to show", nil)
		ctx.JSON(http.StatusOK, response)
		return
	}

	response := res.SuccessResponse(200, "succesfully go all users coupons", coupons)
	ctx.JSON(http.StatusOK, response)
}

// UpdateCoupon godoc
// @summary api for admin to update the coupon
// @security ApiKeyAuth
// @tags Admin Coupon
// @id UpdateCoupon
// @Param        inputs   body     req.EditCoupon{}   true  "Input Field"
// @Router /admin/coupons [put]
// @Success 200 {object} res.Response{} "successfully update the coupon"
// @Failure 400 {object} res.Response{}  "invalid input"
func (c *CouponHandler) UpdateCoupon(ctx *gin.Context) {

	var body req.EditCoupon

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response := res.ErrorResponse(400, "invalid input", err.Error(), body)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	var coupon domain.Coupon

	copier.Copy(&coupon, &body)

	err := c.couponUseCase.UpdateCoupon(ctx, coupon)
	if err != nil {
		response := res.ErrorResponse(400, "faild to update coupon", err.Error(), coupon)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := res.SuccessResponse(200, "successfully update the coupon", coupon)
	ctx.JSON(http.StatusOK, response)
}

// ApplyUserCoupon godoc
// @summary api user to apply on cart on checkout time
// @security ApiKeyAuth
// @tags User Cart
// @id ApplyCouponToCart
// @Param        inputs   body     req.ApplyCoupon{}   true  "Input Field"
// @Router /carts/coupons [patch]
// @Success 200 {object} res.Response{} "successfully updated the coupon code"
// @Failure 400 {object} res.Response{}  "invalid input"
func (c *CouponHandler) ApplyCouponToCart(ctx *gin.Context) {

	var body req.ApplyCoupon

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response := res.ErrorResponse(400, "invalid input", err.Error(), body)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	userID := utils.GetUserIdFromContext(ctx)

	discountPrice, err := c.couponUseCase.ApplyCouponToCart(ctx, userID, body.CouponCode)
	if err != nil {
		respone := res.ErrorResponse(400, "faild to apply the coupon code", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, respone)
		return
	}

	response := res.SuccessResponse(200, "successfully applied the coupon code on cart", gin.H{"discount_amount": discountPrice})
	ctx.JSON(http.StatusOK, response)
}
