package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
)

type CouponRepository interface {
	CheckCouponDetailsAlreadyExist(ctx context.Context, coupon domain.Coupon) (couponID uint, err error)
	FindCouponByID(ctx context.Context, couponID uint) (coupon domain.Coupon, err error)

	FindCouponByCouponCode(ctx context.Context, couponCode string) (coupon domain.Coupon, err error)
	FindCouponByName(ctx context.Context, couponName string) (coupon domain.Coupon, err error)

	FindAllCoupons(ctx context.Context, pagination req.ReqPagination) (coupons []domain.Coupon, err error)
	SaveCoupon(ctx context.Context, coupon domain.Coupon) error
	UpdateCoupon(ctx context.Context, coupon domain.Coupon) error

	// uses coupon
	FindCouponUses(ctx context.Context, userID, couopnID uint) (couponUses domain.CouponUses, err error)

	// find all coupon for user
	FindAllCouponForUser(ctx context.Context, userID uint, pagination req.ReqPagination) (coupons []res.ResUserCoupon, err error)

	//!cart
	FindCartByUserID(ctx context.Context, userID uint) (cart domain.Cart, err error)
	UpdateCart(ctx context.Context, cartId, discountAmount, couponID uint) error
	//!end
}
