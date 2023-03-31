package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/res"
)

type CouponUseCase interface {
	// coupon
	AddCoupon(ctx context.Context, coupon domain.Coupon) error
	GetAllCoupons(ctx context.Context) ([]domain.Coupon, error)
	UpdateCoupon(ctx context.Context, coupon domain.Coupon) error

	// user_coupon
	AddUserCoupon(ctx context.Context, userID uint) (domain.UserCoupon, error)
	GetAllUserCoupons(ctx context.Context, userID uint) ([]res.ResUserCoupon, error)
	ApplyUserCoupon(ctx context.Context, couponCode string, toatalPrice uint) (domain.UserCoupon, error)
}
