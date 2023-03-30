package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
)

type CouponUseCase interface {
	AddCoupon(ctx context.Context, coupon domain.Coupon) error
	GetAllCoupons(ctx context.Context) ([]domain.Coupon, error)
	UpdateCoupon(ctx context.Context, coupon domain.Coupon) error
}
