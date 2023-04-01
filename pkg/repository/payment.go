package repository

import (
	"context"
	"errors"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
)

func (c *OrderDatabase) FindPaymentMethodByID(ctx context.Context, paymenMethodtID uint) (domain.PaymentMethod, error) {
	var paymentMethods domain.PaymentMethod
	query := `SELECT * FROM payment_methods WHERE id = ?`

	if c.DB.Raw(query, paymenMethodtID).Scan(&paymentMethods).Error != nil {
		return paymentMethods, errors.New("faild to find payment method")
	}

	return paymentMethods, nil
}

func (c *OrderDatabase) FindAllPaymentMethods(ctx context.Context) ([]domain.PaymentMethod, error) {

	var paymentMethods []domain.PaymentMethod
	query := `SELECT * FROM payment_methods`
	if c.DB.Raw(query).Scan(&paymentMethods).Error != nil {
		return paymentMethods, errors.New("faild to find all payment methods")
	}
	return paymentMethods, nil
}
