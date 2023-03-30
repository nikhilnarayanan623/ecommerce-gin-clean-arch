package usecase

import (
	"context"
	"errors"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	service "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
)

type couponUseCase struct {
	couponRepo interfaces.CouponRepository
}

func NewCouponUseCase(couponRepo interfaces.CouponRepository) service.CouponUseCase {
	return &couponUseCase{couponRepo: couponRepo}
}

func (c *couponUseCase) AddCoupon(ctx context.Context, coupon domain.Coupon) error {
	// first check coupon already exist with this coupon name
	coupon, err := c.couponRepo.FindCoupon(ctx, coupon)
	if err != nil {
		return err
	} else if coupon.ID != 0 {
		return errors.New("ther is another coupon already exist with this coupon_name")
	}

	return c.couponRepo.SaveCoupon(ctx, coupon)
}
func (c *couponUseCase) GetAllCoupons(ctx context.Context) ([]domain.Coupon, error) {

	return c.couponRepo.FindAllCoupons(ctx)
}

func (c *couponUseCase) UpdateCoupon(ctx context.Context, coupon domain.Coupon) error {

	// first check the coupon_id is valid or not
	if checkCoupon, err := c.couponRepo.FindCoupon(ctx, domain.Coupon{ID: coupon.ID}); err != nil {
		return err
	} else if checkCoupon.CouponName == "" {
		return errors.New("invalid coupon_id \nthere is no coupon available with this coupn_id")
	}

	// update the coupon
	return c.couponRepo.UpdateCoupon(ctx, coupon)
}
