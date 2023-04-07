package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
)

type CouponHandler struct {
	couponUseCase interfaces.CouponUseCase
}

func NewCouponHandler(couponUseCase interfaces.CouponUseCase) *CouponHandler {
	return &CouponHandler{couponUseCase: couponUseCase}
}

// AddCoupon godoc
// @summary api for admin to add coupon
// @security ApiKeyAuth
// @tags Admin Coupon
// @id AddCoupon
// @Param        inputs   body     domain.Coupon{}   true  "Input Field"
// @Router /admin/coupons [post]
// @Success 200 {object} res.Response{} "successfully added coupon"
// @Failure 400 {object} res.Response{}  "invalid input"
func (c *CouponHandler) AddCoupon(ctx *gin.Context) {

	var body req.ReqCoupon

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

	response := res.SuccessResponse(200, "successfully added coupon", nil)
	ctx.JSON(http.StatusOK, response)
}

// GetAllCoupons godoc
// @summary api for admin to see all coupons
// @security ApiKeyAuth
// @tags Admin Coupon
// @id GetAllCoupons
// @Router /admin/coupons [get]
// @Success 200 {object} res.Response{} "successfully update the coupon"
// @Failure 400 {object} res.Response{}  "invalid input"
func (c *CouponHandler) GetAllCoupons(ctx *gin.Context) {
	coupons, err := c.couponUseCase.GetAllCoupons(ctx)
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

// UpdateCoupon godoc
// @summary api for admin to update the coupon
// @security ApiKeyAuth
// @tags Admin Coupon
// @id UpdateCoupon
// @Param        inputs   body     domain.Coupon{}   true  "Input Field"
// @Router /admin/coupons [put]
// @Success 200 {object} res.Response{} "successfully update the coupon"
// @Failure 400 {object} res.Response{}  "invalid input"
func (c *CouponHandler) UpdateCoupon(ctx *gin.Context) {

	var body req.ReqCoupon

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

// !
func (c *CouponHandler) ApplyCouponToCart(ctx *gin.Context) {

	userID := utils.GetUserIdFromContext(ctx)
	var body req.ReqApplyCoupon
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response := res.ErrorResponse(400, "invalid input", err.Error(), body)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	discountPrice, err := c.couponUseCase.ApplyCouponToCart(ctx, userID, body.CouponCode)
	if err != nil {
		respone := res.ErrorResponse(400, "faild to apply the coupon code", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, respone)
		return
	}

	response := res.SuccessResponse(200, "successfully updated the coupon code", gin.H{"discount_price": discountPrice})
	ctx.JSON(http.StatusOK, response)
}

// func (c *CouponHandler) CheckUserCouponChance(ctx *gin.Context) {

// 	// check ther probability and if no probability then return
// 	if !utils.CheckProbability(0.7) {
// 		response := res.SuccessResponse(200, "there is no coupon better luck next time", nil)
// 		ctx.JSON(http.StatusOK, response)
// 		return
// 	}

// 	userID := utils.GetUserIdFromContext(ctx)

// 	//save coupon for use
// 	userCoupon, err := c.couponUseCase.AddUserCoupon(ctx, userID)
// 	if err != nil {
// 		response := res.ErrorResponse(500, "faild to create coupon for user", err.Error(), userCoupon)
// 		ctx.JSON(http.StatusInternalServerError, response)
// 		return
// 	}

// 	response := res.SuccessResponse(http.StatusOK, "successfully created a coupon for user", userCoupon)
// 	ctx.JSON(200, response)

// }

// // GetAllUserCoupons godoc
// // @summary api user to see all coupons that use got
// // @security ApiKeyAuth
// // @tags User Coupons
// // @id GetAllUserCoupons
// // @Router /coupons [get]
// // @Success 200 {object} res.Response{} "successfully copon code applied for cart"
// // @Failure 400 {object} res.Response{}  "invalid input"
// func (c *CouponHandler) GetAllUserCoupons(ctx *gin.Context) {

// 	userID := utils.GetUserIdFromContext(ctx)

// 	userCoupons, err := c.couponUseCase.GetAllUserCoupons(ctx, userID)
// 	if err != nil {
// 		response := res.ErrorResponse(500, "faild to get user coupons", err.Error(), nil)
// 		ctx.JSON(http.StatusInternalServerError, response)
// 		return
// 	}

// 	if userCoupons == nil {
// 		response := res.SuccessResponse(200, "there is no coupons for user", nil)
// 		ctx.JSON(http.StatusOK, response)
// 		return
// 	}

// 	respones := res.SuccessResponse(200, "successfully got user coupons", userCoupons)
// 	ctx.JSON(http.StatusOK, respones)
// }

// // ApplyUserCoupon godoc
// // @summary api user to apply on cart on checkout time
// // @security ApiKeyAuth
// // @tags User Cart
// // @id ApplyUserCoupon
// // @Param        inputs   body     req.ReqApplyCoupon{}   true  "Input Field"
// // @Router /carts/apply-coupon [patch]
// // @Success 200 {object} res.Response{} "successfully copon code applied for cart"
// // @Failure 400 {object} res.Response{}  "invalid input"
// func (c *CouponHandler) ApplyUserCoupon(ctx *gin.Context) {

// 	var body req.ReqApplyCoupon

// 	userID := utils.GetUserIdFromContext(ctx)

// 	if err := ctx.ShouldBindJSON(&body); err != nil {
// 		response := res.ErrorResponse(400, "invalid input", err.Error(), body)
// 		ctx.JSON(http.StatusBadRequest, response)
// 		return
// 	}

// 	userCoupon, err := c.couponUseCase.ApplyUserCoupon(ctx, userID, body.CouponCode)

// 	if err != nil {
// 		response := res.ErrorResponse(400, "faild to apply coupon_code", err.Error(), body)
// 		ctx.JSON(http.StatusBadRequest, response)
// 		return
// 	}

// 	response := res.SuccessResponse(200, "successfully coupon code applied", userCoupon)
// 	ctx.JSON(http.StatusOK, response)
// }
