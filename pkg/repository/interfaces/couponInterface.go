package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
)

type CouponRepository interface {
	FindCoupon(ctx context.Context, coupon domain.Coupon) (domain.Coupon, error)
	FindAllCoupons(ctx context.Context) ([]domain.Coupon, error)
	SaveCoupon(ctx context.Context, coupon domain.Coupon) error
	UpdateCoupon(ctx context.Context, coupon domain.Coupon) error
}
