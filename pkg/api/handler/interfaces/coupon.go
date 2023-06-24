package interfaces

import "github.com/gin-gonic/gin"

type CouponHandler interface {
	SaveCoupon(ctx *gin.Context)
	FindAllCoupons(ctx *gin.Context)
	FindAllCouponsForUser(ctx *gin.Context)
	UpdateCoupon(ctx *gin.Context)
	ApplyCouponToCart(ctx *gin.Context)
}
