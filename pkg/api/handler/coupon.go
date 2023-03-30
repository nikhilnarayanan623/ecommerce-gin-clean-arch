package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
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
