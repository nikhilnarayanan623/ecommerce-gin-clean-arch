package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
)

func (c *OrderDatabase) FindPaymentMethodByID(ctx context.Context, paymenMethodtID uint) (paymentMethods domain.PaymentMethod, err error) {

	query := `SELECT * FROM payment_methods WHERE id = $1`

	err = c.DB.Raw(query, paymenMethodtID).Scan(&paymentMethods).Error
	if err != nil {
		return paymentMethods, fmt.Errorf("faild to find payment_method by id %v \n%v", paymenMethodtID, err.Error())
	}

	return paymentMethods, nil
}

// find payment_method by payment_type
func (c *OrderDatabase) FindPaymentMethodByType(ctx context.Context, paymentType string) (paymentMethod domain.PaymentMethod, err error) {

	query := `SELECT * FROM payment_methods WHERE payment_type = $1`
	err = c.DB.Raw(query, paymentType).Scan(&paymentMethod).Error

	if err != nil {
		return paymentMethod, fmt.Errorf("faild to find payment_method by payment_type %v \n %v", paymentType, err.Error())
	}
	return paymentMethod, nil
}

func (c *OrderDatabase) FindAllPaymentMethods(ctx context.Context) ([]domain.PaymentMethod, error) {

	var paymentMethods []domain.PaymentMethod
	query := `SELECT * FROM payment_methods`
	if c.DB.Raw(query).Scan(&paymentMethods).Error != nil {
		return paymentMethods, errors.New("faild to find all payment methods")
	}
	return paymentMethods, nil
}

func (c *OrderDatabase) SavePaymentMethod(ctx context.Context, paymentMethod domain.PaymentMethod) (paymentMethodID uint, err error) {

	query := `INSERT INTO payment_methods (payment_type,block_status,maximum_amount) VALUES ($1, $2, $3)`

	err = c.DB.Raw(query, paymentMethod.PaymentType, paymentMethod.BlockStatus, paymentMethod.MaximumAmount).Scan(&paymentMethod).Error

	if err != nil {
		return paymentMethodID, fmt.Errorf("faild to save payment method on insert \n %v", err.Error())
	}
	return paymentMethod.ID, nil
}
func (c *OrderDatabase) UpdatePaymentMethod(ctx context.Context, paymentMethod domain.PaymentMethod) error {

	query := `UPDATE payment_methods SET payment_type = $1, block_status = $2, maximum_amount = $3 WHERE id = $4`

	err := c.DB.Exec(query, paymentMethod.PaymentType, paymentMethod.BlockStatus, paymentMethod.MaximumAmount, paymentMethod.ID).Error

	if err != nil {
		return fmt.Errorf("faild to update payment_method on update %v", err.Error())
	}

	return nil
}
