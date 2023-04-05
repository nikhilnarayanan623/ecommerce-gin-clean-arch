package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
)

type CouponRepository interface {
	FindCoupon(ctx context.Context, coupon domain.Coupon) (domain.Coupon, error)
	FindAllCoupons(ctx context.Context) ([]domain.Coupon, error)
	SaveCoupon(ctx context.Context, coupon domain.Coupon) error
	UpdateCoupon(ctx context.Context, coupon domain.Coupon) error

	// user_coupon
	FindUserCouponByCouponCode(ctx context.Context, couponCode string) (domain.UserCoupon, error)
	FindAllUserCouponsByUserID(ctx context.Context, userID uint) ([]res.ResUserCoupon, error)
	SaveUserCoupon(ctx context.Context, userCoupon domain.UserCoupon) error
	UpdateUserCoupon(ctx context.Context, userCoupon domain.UserCoupon) error
	FindCartTotalPrice(ctx context.Context, userID uint, includeOutOfStck bool) (uint, error)
}
