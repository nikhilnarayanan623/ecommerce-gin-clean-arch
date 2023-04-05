package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
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

// save coupon code for user
func (c *couponUseCase) AddUserCoupon(ctx context.Context, userID uint) (domain.UserCoupon, error) {

	// first all coupons
	coupons, err := c.couponRepo.FindAllCoupons(ctx)
	if err != nil {
		return domain.UserCoupon{}, err
	} else if coupons == nil {
		return domain.UserCoupon{}, errors.New("there is no coupons available")
	}

	// then slelect a random coupon and set its id to user_coupons
	randomCouponID := coupons[utils.SelectRandomNumber(0, len(coupons))].ID
	// create a random coupon for user_coupon
	randomCouponCode := utils.CreateRandomCouponCode(10)
	// select a random date
	randomExpireDate := time.Now().AddDate(0, 0, utils.SelectRandomNumber(10, 30))

	// create a useCoupon with this details

	userCoupon := domain.UserCoupon{
		UserID:     userID,
		CouponID:   randomCouponID,
		CouponCode: randomCouponCode,
		ExpireDate: randomExpireDate,
	}

	// return user_coupn with save coupon
	return userCoupon, c.couponRepo.SaveUserCoupon(ctx, userCoupon)
}

// get all user_coupons
func (c *couponUseCase) GetAllUserCoupons(ctx context.Context, userID uint) ([]res.ResUserCoupon, error) {
	return c.couponRepo.FindAllUserCouponsByUserID(ctx, userID)
}

// apply a user_coupon
func (c *couponUseCase) ApplyUserCoupon(ctx context.Context, userID uint, couponCode string) (domain.UserCoupon, error) {

	// first get the coupon using coupon_code and validate it
	userCoupon, err := c.couponRepo.FindUserCouponByCouponCode(ctx, couponCode)
	if err != nil {
		return userCoupon, err
	} else if userCoupon.ID == 0 {
		return userCoupon, errors.New("invalid coupon_code")
	}

	// check coupon is used or not
	if userCoupon.Used {
		return userCoupon, errors.New("coupon is already used")
	}

	// check coupon expire time and last_applied
	if time.Since(userCoupon.ExpireDate) > 0 {
		return userCoupon, errors.New("can't apply coupon \ncoupon expired")
	} else if time.Since(userCoupon.LastApplied.AddDate(0, 0, 1)) < 0 { // check the coupon is within a date or not
		return userCoupon, errors.New("can't apply coupon \ncoupon already applied on cart \ntry after a day")
	}

	// find the coupon and check its minmum price its requird or not
	coupon, err := c.couponRepo.FindCoupon(ctx, domain.Coupon{ID: userCoupon.CouponID})
	if err != nil {
		return userCoupon, err
	}

	cartPrice, err := c.couponRepo.FindCartTotalPrice(ctx, userID, false) // find total price of cart exclude the the outOfStock

	if err != nil {
		return userCoupon, err
	} else if cartPrice == 0 {
		return userCoupon, errors.New("cart is empty or cart items out of stock")
	} else if cartPrice < coupon.MinimumPrice {
		return userCoupon, fmt.Errorf("can't apply coupon \nrequired order_toatal price not met with coupon minimum_price %d", coupon.MinimumPrice)
	}

	// calucalte a random discount price with coupun_percentage_upto form 5
	randomDisountRate := utils.SelectRandomNumber(5, int(coupon.PercentageUpto))
	fmt.Println("discount rate ", randomDisountRate, coupon.PercentageUpto)
	discountPrice := (cartPrice * uint(randomDisountRate)) / 100
	fmt.Println("discount price", discountPrice)

	// set the disount_price and last applied time to use_coupn and update it
	userCoupon.DiscountAmount = discountPrice
	userCoupon.LastApplied = time.Now()
	userCoupon.CartPrice = cartPrice

	return userCoupon, c.couponRepo.UpdateUserCoupon(ctx, userCoupon)
}
