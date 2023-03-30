package repository

import (
	"context"
	"errors"

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

func (c *couponDatabase) FindCoupon(ctx context.Context, coupon domain.Coupon) (domain.Coupon, error) {

	query := `SELECT * FROM coupons WHERE id = ? OR coupon_name = ?`
	if c.DB.Raw(query, coupon.ID, coupon.CouponName).Scan(&coupon).Error != nil {
		return coupon, errors.New("faild to find coupon")
	}
	return coupon, nil
}
func (c *couponDatabase) FindAllCoupons(ctx context.Context) ([]domain.Coupon, error) {

	var coupons []domain.Coupon
	query := `SELECT * FROM coupons`
	if c.DB.Raw(query).Scan(&coupons).Error != nil {
		return coupons, errors.New("faild to find coupon")
	}
	return coupons, nil
}
func (c *couponDatabase) SaveCoupon(ctx context.Context, coupon domain.Coupon) error {
	query := `INSERT INTO coupons (coupon_name, description, percentage_upto, minimum_price, image,block_status) 
	VALUES($1,$2,$3,$4,$5,$6)`
	if c.DB.Exec(query, coupon.CouponName, coupon.Description,
		coupon.PercentageUpto, coupon.MinimumPrice, coupon.Image, false).Error != nil {
		return errors.New("faild to save coupon")
	}
	return nil
}

func (c *couponDatabase) UpdateCoupon(ctx context.Context, coupon domain.Coupon) error {
	query := `UPDATE coupons SET coupon_name = $1, description = $2, percentage_upto = $3, 
	minimum_price = $4, image = $5, block_status = $6 WHERE id = $7`
	if c.DB.Exec(query, coupon.CouponName, coupon.Description, coupon.PercentageUpto,
		coupon.MinimumPrice, coupon.Image, coupon.BlockStatus, coupon.ID).Error != nil {
		return errors.New("faild to update coupon")
	}
	return nil
}

func (c *couponDatabase) FindUserCouponByCouponCode(ctx context.Context, couponCode string) (domain.UserCoupons, error) {
	var userCoupon domain.UserCoupons
	query := `SELECT * FROM user_coupons WHERE coupon_code = ?`
	if c.DB.Raw(query, couponCode).Scan(&userCoupon).Error != nil {
		return userCoupon, errors.New("faild to find user_coupon with coupon_code")
	}
	return userCoupon, nil
}
