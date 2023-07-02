package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/response"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/request"
)

type CouponRepository interface {
	CheckCouponDetailsAlreadyExist(ctx context.Context, coupon domain.Coupon) (couponID uint, err error)
	FindCouponByID(ctx context.Context, couponID uint) (coupon domain.Coupon, err error)

	FindCouponByCouponCode(ctx context.Context, couponCode string) (coupon domain.Coupon, err error)
	FindCouponByName(ctx context.Context, couponName string) (coupon domain.Coupon, err error)

	FindAllCoupons(ctx context.Context, pagination request.Pagination) (coupons []domain.Coupon, err error)
	SaveCoupon(ctx context.Context, coupon domain.Coupon) error
	UpdateCoupon(ctx context.Context, coupon domain.Coupon) error

	// uses coupon
	FindCouponUsesByCouponAndUserID(ctx context.Context, userID, couopnID uint) (couponUses domain.CouponUses, err error)
	SaveCouponUses(ctx context.Context, couponUses domain.CouponUses) error

	// find all coupon for user
	FindAllCouponForUser(ctx context.Context, userID uint, pagination request.Pagination) (coupons []response.UserCoupon, err error)

}
