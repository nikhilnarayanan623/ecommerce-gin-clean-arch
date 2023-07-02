package db

import (
	"fmt"
	"time"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
	"gorm.io/gorm"
)

// To save predefined order statuses on database if its not exist
func saveOrderStatuses(db *gorm.DB) error {

	statuses := []domain.OrderStatusType{
		domain.StatusPaymentPending,
		domain.StatusOrderPlaced,
		domain.StatusOrderCancelled,
		domain.StatusOrderDelivered,
		domain.StatusReturnRequested,
		domain.StatusReturnApproved,
		domain.StatusReturnCancelled,
		domain.StatusOrderReturned,
	}

	var (
		searchQuery = `SELECT CASE WHEN id != 0 THEN 'T' ELSE 'F' END as exist 
		FROM order_statuses WHERE status = $1`
		insertQuery = `INSERT INTO order_statuses (status) VALUES ($1)`
		exist       bool
		err         error
	)

	for _, status := range statuses {

		err = db.Raw(searchQuery, status).Scan(&exist).Error
		if err != nil {
			return fmt.Errorf("failed to check order status already exist err: %w", err)
		}

		if !exist {
			err = db.Exec(insertQuery, status).Error
			if err != nil {
				return fmt.Errorf("failed to save status %w", err)
			}
		}
		exist = false
	}
	return nil
}

// To save predefined payment methods on database if its not exist
func savePaymentMethods(db *gorm.DB) error {
	paymentMethods := []domain.PaymentMethod{
		{
			Name:          domain.CodPayment,
			MaximumAmount: domain.CodMaximumAmount,
		},
		{
			Name:          domain.RazopayPayment,
			MaximumAmount: domain.RazorPayMaximumAmount,
		},
		{
			Name:          domain.StripePayment,
			MaximumAmount: domain.StripeMaximumAmount,
		},
	}

	var (
		searchQuery = `SELECT CASE WHEN id != 0 THEN 'T' ELSE 'F' END as exist FROM payment_methods WHERE name = $1`
		insertQuery = `INSERT INTO payment_methods (name, maximum_amount) VALUES ($1, $2)`
		exist       bool
		err         error
	)

	for _, paymentMethod := range paymentMethods {

		err = db.Raw(searchQuery, paymentMethod.Name).Scan(&exist).Error
		if err != nil {
			return fmt.Errorf("failed to check payment methods already exist %w", err)
		}
		if !exist {
			err = db.Exec(insertQuery, paymentMethod.Name, paymentMethod.MaximumAmount).Error
			if err != nil {
				return fmt.Errorf("failed to save payment method %w", err)
			}
		}
		exist = false
	}
	return nil
}

func saveAdmin(db *gorm.DB, email, userName, password string) error {

	var (
		searchQuery = `SELECT CASE WHEN id != 0 THEN 'T' ELSE 'F' END as exist FROM admins WHERE email = $1`
		insertQuery = `INSERT INTO admins (email, user_name, password, created_at) VALUES ($1, $2, $3, $4)`
		exist       bool
		err         error
	)

	err = db.Raw(searchQuery, email).Scan(&exist).Error
	if err != nil {
		return fmt.Errorf("failed to check admin already exist err:%w", err)
	}

	if !exist {
		hashPass, err := utils.GetHashedPassword(password)
		if err != nil {
			return fmt.Errorf("failed to hash password err: %w", err)
		}
		createdAt := time.Now()
		err = db.Exec(insertQuery, email, userName, hashPass, createdAt).Error
		if err != nil {
			return fmt.Errorf("failed to save admin details %w", err)
		}
	}
	return nil
}
