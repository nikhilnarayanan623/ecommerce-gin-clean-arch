package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	service "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
)

type couponUseCase struct {
	couponRepo interfaces.CouponRepository
	cartRepo   interfaces.CartRepository
}

func NewCouponUseCase(couponRepo interfaces.CouponRepository, cartRepo interfaces.CartRepository) service.CouponUseCase {
	return &couponUseCase{
		couponRepo: couponRepo,
		cartRepo:   cartRepo,
	}
}

func (c *couponUseCase) AddCoupon(ctx context.Context, coupon domain.Coupon) error {
	// first check coupon already exist with this coupon name
	checkCoupon, err := c.couponRepo.FindCouponByName(ctx, coupon.CouponName)
	if err != nil {
		return err
	} else if checkCoupon.CouponID != 0 {
		return fmt.Errorf("there already a coupon exist with coupon_name %v", coupon.CouponName)
	}
	// validate the coupn expire date
	if time.Since(coupon.ExpireDate) > 0 {
		return fmt.Errorf("given coupon expire date already exceeded %v", coupon.ExpireDate)
	}

	// check the given expire time is valid or not

	if time.Since(coupon.ExpireDate) > 0 {
		return fmt.Errorf("given expire date is already over \ngiven time %v", coupon.ExpireDate)
	}

	// create a random coupon code
	coupon.CouponCode = utils.GenerateCouponCode(10)

	// create a coupon
	err = c.couponRepo.SaveCoupon(ctx, coupon)
	if err != nil {
		return err
	}

	return nil
}
func (c *couponUseCase) GetAllCoupons(ctx context.Context, pagination req.Pagination) (coupons []domain.Coupon, err error) {

	coupons, err = c.couponRepo.FindAllCoupons(ctx, pagination)
	if err != nil {
		return coupons, err
	}

	log.Printf("successfully got all coupons \n\n")
	return coupons, nil
}

// get all coupon for user
func (c *couponUseCase) GetCouponsForUser(ctx context.Context, userID uint, pagination req.Pagination) (coupons []res.UserCoupon, err error) {

	coupons, err = c.couponRepo.FindAllCouponForUser(ctx, userID, pagination)

	if err != nil {
		return coupons, err
	}

	log.Printf("successfully go coupons for user of user_id %v", userID)

	return coupons, nil
}

func (c *couponUseCase) GetCouponByCouponCode(ctx context.Context, couponCode string) (coupon domain.Coupon, err error) {
	coupon, err = c.couponRepo.FindCouponByCouponCode(ctx, couponCode)

	if err != nil {
		return coupon, err
	} else if coupon.CouponID == 0 {
		return coupon, fmt.Errorf("invalid coupon code %s", couponCode)
	}
	return coupon, nil
}

func (c *couponUseCase) UpdateCoupon(ctx context.Context, coupon domain.Coupon) error {

	// first check the coupon_id is valid or not
	checkCoupon, err := c.couponRepo.FindCouponByID(ctx, coupon.CouponID)
	if err != nil {
		return err
	} else if checkCoupon.CouponID == 0 {
		return fmt.Errorf("invalid coupon_id %v", coupon.CouponID)
	}

	// check any coupon already exist wtih this details
	couponID, err := c.couponRepo.CheckCouponDetailsAlreadyExist(ctx, coupon)

	if err != nil {
		return err
	} else if couponID != 0 {
		return fmt.Errorf("another coupon already exist with this details with coupon_id %v", couponID)
	}

	if time.Since(coupon.ExpireDate) > 0 {
		return fmt.Errorf("given expire date is already over \ngiven time %v", coupon.ExpireDate)
	}

	// then update the coupon
	err = c.couponRepo.UpdateCoupon(ctx, coupon)
	if err != nil {
		return err
	}

	return nil
}

// apply coupon
func (c *couponUseCase) ApplyCouponToCart(ctx context.Context, userID uint, couponCode string) (discountAmount uint, err error) {

	// get the coupon with given coupon code
	coupon, err := c.couponRepo.FindCouponByCouponCode(ctx, couponCode)
	if err != nil {
		return discountAmount, err
	} else if coupon.CouponID == 0 {
		return discountAmount, fmt.Errorf("invalid coupon_code %s", couponCode)
	}

	// check the coupon is user already used or not
	couponUses, err := c.couponRepo.FindCouponUsesByCouponAndUserID(ctx, userID, coupon.CouponID)
	if err != nil {
		return discountAmount, err
	} else if couponUses.CouponUsesID != 0 {
		return discountAmount, fmt.Errorf("user already applied this coupon at %v", couponUses.UsedAt)
	}

	// get the cart of user
	cart, err := c.cartRepo.FindCartByUserID(ctx, userID)
	if err != nil {
		return discountAmount, err
	} else if cart.CartID == 0 {
		return discountAmount, fmt.Errorf("there is no cart_items avialable for user with user_id %d", userID)
	}

	// then check the cart have already a coupon applied
	if cart.AppliedCouponID != 0 {
		return discountAmount, fmt.Errorf("cart have already a coupon applied with coupon_id %d", cart.AppliedCouponID)
	}

	// validate the coupon expire date and cart price
	if time.Since(coupon.ExpireDate) > 0 {
		return discountAmount, fmt.Errorf("can't apply coupn \ncoupn expired")
	}
	if cart.TotalPrice < coupon.MinimumCartPrice {
		return discountAmount, fmt.Errorf("can't apply coupn \ncoupn minimum cart_amount %d not met with user cart total price %d",
			coupon.MinimumCartPrice, cart.TotalPrice)
	}

	// calculate a discount for cart
	discountAmount = (cart.TotalPrice * coupon.DiscountRate) / 100
	// update the cart
	err = c.cartRepo.UpdateCart(ctx, cart.CartID, discountAmount, coupon.CouponID)
	if err != nil {
		return discountAmount, err
	}

	log.Printf("successfully updated the cart price with dicount price %d", discountAmount)
	return discountAmount, nil
}
