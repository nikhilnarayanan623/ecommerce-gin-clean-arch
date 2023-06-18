package repository

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/request"
	"gorm.io/gorm"
)

type paymentDatabase struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) interfaces.PaymentRepository {
	return &paymentDatabase{
		db: db,
	}
}

func (c *paymentDatabase) FindPaymentMethodByID(ctx context.Context, paymenMethodtID uint) (paymentMethods domain.PaymentMethod, err error) {

	query := `SELECT * FROM payment_methods WHERE id = $1`

	err = c.db.Raw(query, paymenMethodtID).Scan(&paymentMethods).Error

	return paymentMethods, err
}

// find payment_method by payment_type
func (c *paymentDatabase) FindPaymentMethodByType(ctx context.Context, paymentType string) (paymentMethod domain.PaymentMethod, err error) {

	query := `SELECT * FROM payment_methods WHERE payment_type = $1`
	err = c.db.Raw(query, paymentType).Scan(&paymentMethod).Error

	return paymentMethod, err
}

func (c *paymentDatabase) FindAllPaymentMethods(ctx context.Context) (paymentMethods []domain.PaymentMethod, err error) {

	query := `SELECT * FROM payment_methods`
	err = c.db.Raw(query).Scan(&paymentMethods).Error

	return paymentMethods, err
}

func (c *paymentDatabase) SavePaymentMethod(ctx context.Context, paymentMethod domain.PaymentMethod) (paymentMethodID uint, err error) {

	query := `INSERT INTO payment_methods (payment_type,block_status,maximum_amount) VALUES ($1, $2, $3)`
	err = c.db.Raw(query, paymentMethod.PaymentType, paymentMethod.BlockStatus, paymentMethod.MaximumAmount).Scan(&paymentMethod).Error

	return paymentMethod.ID, err
}
func (c *paymentDatabase) UpdatePaymentMethod(ctx context.Context, paymentMethod request.PaymentMethodUpdate) error {

	query := `UPDATE payment_methods SET  block_status = $1, maximum_amount = $2 WHERE id = $3`

	err := c.db.Exec(query, paymentMethod.BlockStatus, paymentMethod.MaximumAmount, paymentMethod.ID).Error

	return err
}
