package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/response"
)

type CouponUseCase interface {
	// coupon
	AddCoupon(ctx context.Context, coupon domain.Coupon) error
	GetAllCoupons(ctx context.Context, pagination request.Pagination) (coupons []domain.Coupon, err error)
	UpdateCoupon(ctx context.Context, coupon domain.Coupon) error

	//user side coupons
	GetCouponsForUser(ctx context.Context, userID uint, pagination request.Pagination) (coupons []response.UserCoupon, err error)

	GetCouponByCouponCode(ctx context.Context, couponCode string) (coupon domain.Coupon, err error)
	ApplyCouponToCart(ctx context.Context, userID uint, couponCode string) (discountPrice uint, err error)
}
