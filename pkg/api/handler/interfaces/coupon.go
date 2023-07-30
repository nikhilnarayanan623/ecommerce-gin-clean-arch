package interfaces

import "github.com/gin-gonic/gin"

type CouponHandler interface {
	SaveCoupon(ctx *gin.Context)
	GetAllCouponsAdmin(ctx *gin.Context)
	GetAllCouponsForUser(ctx *gin.Context)
	UpdateCoupon(ctx *gin.Context)
	ApplyCouponToCart(ctx *gin.Context)
}
