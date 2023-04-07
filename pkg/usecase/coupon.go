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
)

type couponUseCase struct {
	couponRepo interfaces.CouponRepository
}

func NewCouponUseCase(couponRepo interfaces.CouponRepository) service.CouponUseCase {
	return &couponUseCase{couponRepo: couponRepo}
}

func (c *couponUseCase) AddCoupon(ctx context.Context, coupon domain.Coupon) error {
	// first check coupon already exist with this coupon name
	checkCoupon, err := c.couponRepo.FindCouponByName(ctx, coupon.CouponName)
	if err != nil {
		return err
	} else if checkCoupon.CouponID != 0 {
		return fmt.Errorf("there already a coupon exist with coupon_name %v", coupon.CouponName)
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
func (c *couponUseCase) GetAllCoupons(ctx context.Context) ([]domain.Coupon, error) {

	return c.couponRepo.FindAllCoupons(ctx)
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

	// get the cart of user
	cart, err := c.couponRepo.FindCartByUserID(ctx, userID)
	if err != nil {
		return discountAmount, err
	} else if cart.CartID == 0 {
		return discountAmount, fmt.Errorf("there is no cart_items avialable for user with user_id %d", userID)
	}

	// then check the cart have already a coupon applied
	if cart.AppliedCouponCode != "" {
		return discountAmount, fmt.Errorf("cart have already a coupon applied %s", cart.AppliedCouponCode)
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
	err = c.couponRepo.UpdateCart(ctx, cart.CartID, discountAmount, couponCode)
	if err != nil {
		return discountAmount, err
	}

	log.Printf("successfully updated the cart price with dicount price %d", discountAmount)
	return discountAmount, nil
}

// // save coupon code for user
// func (c *couponUseCase) AddUserCoupon(ctx context.Context, userID uint) (domain.UserCoupon, error) {

// 	// first all coupons
// 	coupons, err := c.couponRepo.FindAllCoupons(ctx)
// 	if err != nil {
// 		return domain.UserCoupon{}, err
// 	} else if coupons == nil {
// 		return domain.UserCoupon{}, errors.New("there is no coupons available")
// 	}

// 	// then slelect a random coupon and set its id to user_coupons
// 	randomCouponID := coupons[utils.SelectRandomNumber(0, len(coupons))].ID
// 	// create a random coupon for user_coupon
// 	randomCouponCode := utils.CreateRandomCouponCode(10)
// 	// select a random date
// 	randomExpireDate := time.Now().AddDate(0, 0, utils.SelectRandomNumber(10, 30))

// 	// create a useCoupon with this details

// 	userCoupon := domain.UserCoupon{
// 		UserID:     userID,
// 		CouponID:   randomCouponID,
// 		CouponCode: randomCouponCode,
// 		ExpireDate: randomExpireDate,
// 	}

// 	// return user_coupn with save coupon
// 	return userCoupon, c.couponRepo.SaveUserCoupon(ctx, userCoupon)
// }

// // get all user_coupons
// func (c *couponUseCase) GetAllUserCoupons(ctx context.Context, userID uint) ([]res.ResUserCoupon, error) {
// 	return c.couponRepo.FindAllUserCouponsByUserID(ctx, userID)
// }

// // apply a user_coupon
// func (c *couponUseCase) ApplyUserCoupon(ctx context.Context, userID uint, couponCode string) (domain.UserCoupon, error) {

// 	// first get the coupon using coupon_code and validate it
// 	userCoupon, err := c.couponRepo.FindUserCouponByCouponCode(ctx, couponCode)
// 	if err != nil {
// 		return userCoupon, err
// 	} else if userCoupon.ID == 0 {
// 		return userCoupon, errors.New("invalid coupon_code")
// 	}

// 	// check coupon is used or not
// 	if userCoupon.Used {
// 		return userCoupon, errors.New("coupon is already used")
// 	}

// 	// check coupon expire time and last_applied
// 	if time.Since(userCoupon.ExpireDate) > 0 {
// 		return userCoupon, errors.New("can't apply coupon \ncoupon expired")
// 	} else if time.Since(userCoupon.LastApplied.AddDate(0, 0, 1)) < 0 { // check the coupon is within a date or not
// 		return userCoupon, errors.New("can't apply coupon \ncoupon already applied on cart \ntry after a day")
// 	}

// 	// find the coupon and check its minmum price its requird or not
// 	coupon, err := c.couponRepo.FindCoupon(ctx, domain.Coupon{ID: userCoupon.CouponID})
// 	if err != nil {
// 		return userCoupon, err
// 	}

// 	cartPrice, err := c.couponRepo.FindCartTotalPrice(ctx, userID, false) // find total price of cart exclude the the outOfStock

// 	if err != nil {
// 		return userCoupon, err
// 	} else if cartPrice == 0 {
// 		return userCoupon, errors.New("cart is empty or cart items out of stock")
// 	} else if cartPrice < coupon.MinimumPrice {
// 		return userCoupon, fmt.Errorf("can't apply coupon \nrequired order_toatal price not met with coupon minimum_price %d", coupon.MinimumPrice)
// 	}

// 	// calucalte a random discount price with coupun_percentage_upto form 5
// 	randomDisountRate := utils.SelectRandomNumber(5, int(coupon.PercentageUpto))
// 	fmt.Println("discount rate ", randomDisountRate, coupon.PercentageUpto)
// 	discountPrice := (cartPrice * uint(randomDisountRate)) / 100
// 	fmt.Println("discount price", discountPrice)

// 	// set the disount_price and last applied time to use_coupn and update it
// 	userCoupon.DiscountAmount = discountPrice
// 	userCoupon.LastApplied = time.Now()
// 	userCoupon.CartPrice = cartPrice

// 	return userCoupon, c.couponRepo.UpdateUserCoupon(ctx, userCoupon)
// }
