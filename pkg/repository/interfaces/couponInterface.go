package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
)

type CouponRepository interface {
	FindCouponByID(ctx context.Context, couponID uint) (coupon domain.Coupon, err error)
	FindCouponByCouponCode(ctx context.Context, couponCode string) (coupon domain.Coupon, err error)
	FindCouponByName(ctx context.Context, couponName string) (coupon domain.Coupon, err error)

	FindAllCoupons(ctx context.Context) (coupons []domain.Coupon, err error)
	SaveCoupon(ctx context.Context, coupon domain.Coupon) error
	UpdateCoupon(ctx context.Context, coupon domain.Coupon) error

	// uses coupon
	FindCouponUses(ctx context.Context, userID, couopnID uint) (couponUses domain.CouponUses, err error)

	//!cart
	FindCartByUserID(ctx context.Context, userID uint) (cart domain.Cart, err error)
	UpdateCart(ctx context.Context, cartId, totalPrice uint, couponCode string) error
	//!end
}
