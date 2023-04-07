package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
)

type CouponUseCase interface {
	// coupon
	AddCoupon(ctx context.Context, coupon domain.Coupon) error
	GetAllCoupons(ctx context.Context) ([]domain.Coupon, error)
	UpdateCoupon(ctx context.Context, coupon domain.Coupon) error

	GetCouponByCouponCode(ctx context.Context, couponCode string) (coupon domain.Coupon, err error)
	ApplyCouponToCart(ctx context.Context, userID uint, couponCode string) (discountPrice uint, err error)
}
