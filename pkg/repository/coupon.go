package repository

import (
	"context"
	"errors"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/res"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	"gorm.io/gorm"
)

type couponDatabase struct {
	DB *gorm.DB
}

func NewCouponRepository(db *gorm.DB) interfaces.CouponRepository {
	return &couponDatabase{DB: db}
}

// find coupon
func (c *couponDatabase) FindCoupon(ctx context.Context, coupon domain.Coupon) (domain.Coupon, error) {

	query := `SELECT * FROM coupons WHERE id = ? OR coupon_name = ?`
	if c.DB.Raw(query, coupon.ID, coupon.CouponName).Scan(&coupon).Error != nil {
		return coupon, errors.New("faild to find coupon")
	}
	return coupon, nil
}

// find all coupon
func (c *couponDatabase) FindAllCoupons(ctx context.Context) ([]domain.Coupon, error) {

	var coupons []domain.Coupon
	query := `SELECT * FROM coupons`
	if c.DB.Raw(query).Scan(&coupons).Error != nil {
		return coupons, errors.New("faild to find coupon")
	}
	return coupons, nil
}

// save a new coupon
func (c *couponDatabase) SaveCoupon(ctx context.Context, coupon domain.Coupon) error {
	query := `INSERT INTO coupons (coupon_name, description, percentage_upto, minimum_price, image,block_status) 
	VALUES($1,$2,$3,$4,$5,$6)`
	if c.DB.Exec(query, coupon.CouponName, coupon.Description,
		coupon.PercentageUpto, coupon.MinimumPrice, coupon.Image, false).Error != nil {
		return errors.New("faild to save coupon")
	}
	return nil
}

// update coupon
func (c *couponDatabase) UpdateCoupon(ctx context.Context, coupon domain.Coupon) error {
	query := `UPDATE coupons SET coupon_name = $1, description = $2, percentage_upto = $3, 
	minimum_price = $4, image = $5, block_status = $6 WHERE id = $7`
	if c.DB.Exec(query, coupon.CouponName, coupon.Description, coupon.PercentageUpto,
		coupon.MinimumPrice, coupon.Image, coupon.BlockStatus, coupon.ID).Error != nil {
		return errors.New("faild to update coupon")
	}
	return nil
}

// find user_coupon by coupon_code
func (c *couponDatabase) FindUserCouponByCouponCode(ctx context.Context, couponCode string) (domain.UserCoupon, error) {
	var userCoupon domain.UserCoupon
	query := `SELECT * FROM user_coupons WHERE coupon_code = ?`
	if c.DB.Raw(query, couponCode).Scan(&userCoupon).Error != nil {
		return userCoupon, errors.New("faild to find user_coupon with coupon_code")
	}
	return userCoupon, nil
}

// find all user_coupons of user
func (c *couponDatabase) FindAllUserCouponsByUserID(ctx context.Context, userID uint) ([]res.ResUserCoupon, error) {
	var userCoupon []res.ResUserCoupon

	query := `SELECT uc.id, uc.coupon_code, c.coupon_name,c.percentage_upto, c.minimum_price ,c.description, c.image, uc.expire_date 
	 FROM user_coupons uc INNER JOIN coupons c ON uc.coupon_id = c.id 
	 WHERE uc.user_id = ?`
	if c.DB.Raw(query, userID).Scan(&userCoupon).Error != nil {
		return userCoupon, errors.New("faild to get user_coupons")
	}
	return userCoupon, nil
}

// save user_coupon
func (c *couponDatabase) SaveUserCoupon(ctx context.Context, userCoupon domain.UserCoupon) error {
	query := `INSERT INTO user_coupons (user_id, coupon_id, coupon_code, expire_date, used) 
	VALUES ($1,$2,$3,$4,$5)`
	if c.DB.Exec(query, userCoupon.UserID, userCoupon.CouponID, userCoupon.CouponCode, userCoupon.ExpireDate, false).Error != nil {
		return errors.New("faild to save user_coupons")
	}
	return nil
}

// update use_coupon
func (c *couponDatabase) UpdateUserCoupon(ctx context.Context, userCoupon domain.UserCoupon) error {
	query := `UPDATE user_coupons SET discount_amount = $1, used = $2, last_applied = $3`
	if c.DB.Exec(query, userCoupon.DiscountAmount, userCoupon.Used, userCoupon.LastApplied).Error != nil {
		return errors.New("faild to update user_coupons")
	}
	return nil
}