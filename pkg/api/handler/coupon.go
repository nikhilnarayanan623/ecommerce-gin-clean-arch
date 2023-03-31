package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/res"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
)

type CouponHandler struct {
	couponUseCase interfaces.CouponUseCase
}

func NewCouponHandler(couponUseCase interfaces.CouponUseCase) *CouponHandler {
	return &CouponHandler{couponUseCase: couponUseCase}
}

func (c *CouponHandler) AddCoupon(ctx *gin.Context) {
	var coupon domain.Coupon
	if err := ctx.ShouldBindJSON(&coupon); err != nil {
		response := res.ErrorResponse(400, "invalid input", err.Error(), coupon)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	err := c.couponUseCase.AddCoupon(ctx, coupon)
	if err != nil {
		response := res.ErrorResponse(400, "faild to add coupon", err.Error(), coupon)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := res.SuccessResponse(200, "successfully added coupon", nil)
	ctx.JSON(http.StatusOK, response)
}

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

// update coupon
func (c *CouponHandler) UpdateCoupon(ctx *gin.Context) {

	var coupon domain.Coupon
	if err := ctx.ShouldBindJSON(&coupon); err != nil {
		response := res.ErrorResponse(400, "invalid input", err.Error(), coupon)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	err := c.couponUseCase.UpdateCoupon(ctx, coupon)
	if err != nil {
		response := res.ErrorResponse(400, "faild to update coupon", err.Error(), coupon)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := res.SuccessResponse(200, "successfully update the coupon", coupon)
	ctx.JSON(http.StatusOK, response)
}

// func create user_coupn
func (c *CouponHandler) AddUserCoupon(ctx *gin.Context) {

	// check ther probability and if no probability then return
	if !helper.CheckProbability(0.5) {
		response := res.SuccessResponse(200, "there is no coupon \nbetter luck next time", nil)
		ctx.JSON(http.StatusOK, response)
		return
	}

	userID := helper.GetUserIdFromContext(ctx)

	//save coupon for use
	userCoupon, err := c.couponUseCase.AddUserCoupon(ctx, userID)
	if err != nil {
		response := res.ErrorResponse(500, "faild to create coupon for user", err.Error(), userCoupon)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := res.SuccessResponse(http.StatusOK, "successfully created a coupon for user", userCoupon)
	ctx.JSON(200, response)

}

func (c *CouponHandler) GetAllUserCoupons(ctx *gin.Context) {

	userID := helper.GetUserIdFromContext(ctx)

	userCoupons, err := c.couponUseCase.GetAllUserCoupons(ctx, userID)
	if err != nil {
		response := res.ErrorResponse(500, "faild to get user coupons", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if userCoupons == nil {
		response := res.SuccessResponse(200, "there is no coupons for user", nil)
		ctx.JSON(http.StatusOK, response)
		return
	}

	respones := res.SuccessResponse(200, "successfully got user coupons", userCoupons)
	ctx.JSON(http.StatusOK, respones)
}

func (c *CouponHandler) ApplyUserCoupon(ctx *gin.Context) {

	var body struct {
		CouponCode string `json:"coupon_code" binding:"required"`
		TotalPrice uint   `json:"total_price" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response := res.ErrorResponse(400, "invalid input", err.Error(), body)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	userCoupon, err := c.couponUseCase.ApplyUserCoupon(ctx, body.CouponCode, body.TotalPrice)
	if err != nil {
		response := res.ErrorResponse(400, "faild to apply coupon_code", err.Error(), body)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := res.SuccessResponse(200, "successfully coupon code applied", userCoupon)
	ctx.JSON(http.StatusOK, response)
}
