package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	usecase "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/response"
)

type CouponHandler struct {
	couponUseCase usecase.CouponUseCase
}

func NewCouponHandler(couponUseCase usecase.CouponUseCase) interfaces.CouponHandler {
	return &CouponHandler{couponUseCase: couponUseCase}
}

// SaveCoupon godoc
// @summary api for admin to add coupon
// @security ApiKeyAuth
// @tags Admin Coupon
// @id SaveCoupon
// @Param        inputs   body     request.Coupon{}   true  "Input Fields"
// @Router /admin/coupons [post]
// @Success 200 {object} response.Response{} "successfully coupon added"
// @Failure 400 {object} response.Response{}  "invalid input"
func (c *CouponHandler) SaveCoupon(ctx *gin.Context) {

	var body request.Coupon

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	var coupon domain.Coupon

	copier.Copy(&coupon, &body)

	err := c.couponUseCase.AddCoupon(ctx, coupon)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, "Failed to add coupon", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully coupon added")
}

// FindAllCoupons godoc
// @summary api for admin to see all coupons
// @security ApiKeyAuth
// @tags Admin Coupon
// @id FindAllCoupons
// @Param page_number query int false "Page Number"
// @Param count query int false "Count"
// @Router /admin/coupons [get]
// @Success 200 {object} response.Response{} "successfully go all the coupons
// @Failure 500 {object} response.Response{}  "faild to get all coupons"
func (c *CouponHandler) FindAllCoupons(ctx *gin.Context) {

	pagination := request.GetPagination(ctx)

	coupons, err := c.couponUseCase.GetAllCoupons(ctx, pagination)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get all coupons", err, nil)
		return
	}

	if len(coupons) == 0 {
		response.SuccessResponse(ctx, http.StatusOK, "No Coupons found", nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully found coupons")
}

// GetAllCoupons godoc
// @summary api for user to see all coupons
// @security ApiKeyAuth
// @tags User Coupon
// @id GetAllCouponsForUser
// @Param page_number query int false "Page Number"
// @Param count query int false "Count Of Order"
// @Router /coupons [get]
// @Success 200 {object} response.Response{} ""Successfully found all coupons for user"
// @Failure 500 {object} response.Response{}  "Failed to find all user"
func (c *CouponHandler) FindAllCouponsForUser(ctx *gin.Context) {

	userID := utils.GetUserIdFromContext(ctx)
	pagination := request.GetPagination(ctx)

	coupons, err := c.couponUseCase.GetCouponsForUser(ctx, userID, pagination)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, "Failed to find all user", err, nil)
		return
	}

	if len(coupons) == 0 {
		response.SuccessResponse(ctx, http.StatusOK, "No coupons found", nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully found all coupons for user", coupons)
}

// UpdateCoupon godoc
// @summary api for admin to update the coupon
// @security ApiKeyAuth
// @tags Admin Coupon
// @id UpdateCoupon
// @Param        inputs   body     request.EditCoupon{}   true  "Input Field"
// @Router /admin/coupons [put]
// @Success 200 {object} response.Response{} "Successfully updated the coupon"
// @Failure 400 {object} response.Response{}  "invalid input"
func (c *CouponHandler) UpdateCoupon(ctx *gin.Context) {

	var body request.EditCoupon

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	var coupon domain.Coupon

	copier.Copy(&coupon, &body)

	err := c.couponUseCase.UpdateCoupon(ctx, coupon)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, "Failed to update coupon", err, coupon)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully updated the coupon", coupon)
}

// ApplyCouponToCart godoc
// @summary api user to apply on cart on checkout time
// @security ApiKeyAuth
// @tags User Cart
// @id ApplyCouponToCart
// @Param        inputs   body     request.ApplyCoupon{}   true  "Input Field"
// @Router /carts/coupons [patch]
// @Success 200 {object} response.Response{} "Successfully coupon applied to user cart"
// @Failure 400 {object} response.Response{}  "invalid input"
func (c *CouponHandler) ApplyCouponToCart(ctx *gin.Context) {

	var body request.ApplyCoupon

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	userID := utils.GetUserIdFromContext(ctx)

	discountPrice, err := c.couponUseCase.ApplyCouponToCart(ctx, userID, body.CouponCode)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, "Failed to apply the coupon code", err, nil)
		return
	}

	data := gin.H{"discount_amount": discountPrice}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully coupon applied to user cart", data)
}
