package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	"gorm.io/gorm"
)

type couponDatabase struct {
	DB *gorm.DB
}

func NewCouponRepository(db *gorm.DB) interfaces.CouponRepository {
	return &couponDatabase{DB: db}
}

// find all coupon
func (c *couponDatabase) FindCouponByID(ctx context.Context, couponID uint) (coupon domain.Coupon, err error) {
	query := `SELECT * FROM coupons WHERE coupon_id = $1`
	err = c.DB.Raw(query, couponID).Scan(&coupon).Error

	if err != nil {
		return coupon, err
	}

	return coupon, nil
}

// find coupon by code
func (c *couponDatabase) FindCouponByCouponCode(ctx context.Context, couponCode string) (coupon domain.Coupon, err error) {

	query := `SELECT * FROM coupons WHERE coupon_code = $1`

	err = c.DB.Raw(query, couponCode).Scan(&coupon).Error
	if err != nil {
		return coupon, fmt.Errorf("faild to find coupon with coupon_code %v", couponCode)
	}

	return coupon, nil
}

// find coupo by name
func (c *couponDatabase) FindCouponByName(ctx context.Context, couponName string) (coupon domain.Coupon, err error) {
	query := `SELECT * FROM coupons WHERE coupon_name = $1`
	err = c.DB.Raw(query, couponName).Scan(&coupon).Error

	if err != nil {
		return coupon, fmt.Errorf("faild to find coupon with coupon_name %v", couponName)
	}

	return coupon, nil
}

func (c *couponDatabase) FindAllCoupons(ctx context.Context) (coupons []domain.Coupon, err error) {

	query := `SELECT * FROM coupons`
	err = c.DB.Raw(query).Scan(&coupons).Error
	if err != nil {
		return coupons, errors.New("faild to find coupon")
	}
	return coupons, nil
}

// save a new coupon
func (c *couponDatabase) SaveCoupon(ctx context.Context, coupon domain.Coupon) error {
	query := `INSERT INTO coupons (coupon_name, coupon_code, description, expire_date, discount_rate, minimum_cart_price, image, block_status)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	err := c.DB.Exec(query, coupon.CouponName, coupon.CouponCode, coupon.Description, coupon.ExpireDate,
		coupon.DiscountRate, coupon.MinimumCartPrice, coupon.Image, coupon.BlockStatus,
	).Error

	if err != nil {
		return fmt.Errorf("faild to save coupon for coupon_name %v", coupon.CouponName)
	}
	return nil
}

// update coupon
func (c *couponDatabase) UpdateCoupon(ctx context.Context, coupon domain.Coupon) error {

	query := `UPDATE coupons SET coupon_name = $1, coupon_code = $2, description = $3, 
	discount_rate = $4, minimum_cart_price = $5, image = $6, block_status = $7`

	err := c.DB.Exec(query, coupon.CouponName, coupon.CouponCode, coupon.Description,
		coupon.DiscountRate, coupon.MinimumCartPrice, coupon.Image, coupon.BlockStatus,
	).Error
	if err != nil {
		return fmt.Errorf("faild to update coupon for coupon_name %v", coupon.CouponName)
	}
	return nil
}

// find couponUses which is also uses for checking a user is a coupon is used or not
func (c *couponDatabase) FindCouponUses(ctx context.Context, userID, couopnID uint) (couponUses domain.CouponUses, err error) {
	query := `SELECT * FROM  coupon_uses WHERE user_id = $1 AND coupon_id = $2`
	err = c.DB.Raw(query, userID, couopnID).Scan(&couponUses).Error
	if err != nil {
		return couponUses, err
	}
	return couponUses, nil
}

// save a couponUses
func (c *couponDatabase) SaveCouponUses(ctx context.Context, couponUses domain.CouponUses) error {
	query := `INSERT INTO coupon_uses ( user_id, coupon_id, used_at) VALUES ($1, $2, $3)`
	err := c.DB.Exec(query, couponUses.UserID, couponUses.CouponID, couponUses.UsedAt).Error

	if err != nil {
		return fmt.Errorf("faild save coupon for user_id %v with coupon_id %v", couponUses.UserID, couponUses.CouponID)
	}

	return nil
}

// !apply coupon cart functions
func (c *couponDatabase) FindCartByUserID(ctx context.Context, userID uint) (cart domain.Cart, err error) {

	query := `SELECT * FROM carts WHERE user_id = ?`
	if c.DB.Raw(query, userID).Scan(&cart).Error != nil {
		return cart, errors.New("faild to get cartItem of user")
	}
	return cart, nil
}

func (c *couponDatabase) UpdateCart(ctx context.Context, cartId, discountAmount uint, couponCode string) error {

	query := `UPDATE carts SET discount_amount = $1, applied_coupon_code = $2 WHERE cart_id = $3`
	err := c.DB.Exec(query, discountAmount, couponCode, cartId).Error
	if err != nil {
		return fmt.Errorf("faild to udpate discount price on cart for coupon")
	}
	return nil
}

//! end
