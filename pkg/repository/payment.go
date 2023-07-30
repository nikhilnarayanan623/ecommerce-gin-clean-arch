package repository

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
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

func (c *paymentDatabase) FindPaymentMethodByID(ctx context.Context, paymentMethodID uint) (paymentMethods domain.PaymentMethod, err error) {

	query := `SELECT * FROM payment_methods WHERE id = $1`

	err = c.db.Raw(query, paymentMethodID).Scan(&paymentMethods).Error

	return paymentMethods, err
}

// find payment_method by name
func (c *paymentDatabase) FindPaymentMethodByType(ctx context.Context,
	paymentType domain.PaymentType) (paymentMethod domain.PaymentMethod, err error) {

	query := `SELECT * FROM payment_methods WHERE name = $1`
	err = c.db.Raw(query, paymentType).Scan(&paymentMethod).Error

	return paymentMethod, err
}

func (c *paymentDatabase) FindAllPaymentMethods(ctx context.Context) (paymentMethods []domain.PaymentMethod, err error) {

	query := `SELECT * FROM payment_methods`
	err = c.db.Raw(query).Scan(&paymentMethods).Error

	return paymentMethods, err
}

func (c *paymentDatabase) SavePaymentMethod(ctx context.Context, paymentMethod domain.PaymentMethod) (paymentMethodID uint, err error) {

	query := `INSERT INTO payment_methods (name, block_status, maximum_amount) VALUES ($1, $2, $3)`
	err = c.db.Raw(query, paymentMethod.Name, paymentMethod.BlockStatus, paymentMethod.MaximumAmount).Scan(&paymentMethod).Error

	return paymentMethod.ID, err
}
func (c *paymentDatabase) UpdatePaymentMethod(ctx context.Context, paymentMethod request.PaymentMethodUpdate) error {

	query := `UPDATE payment_methods SET  block_status = $1, maximum_amount = $2 WHERE id = $3`

	err := c.db.Exec(query, paymentMethod.BlockStatus, paymentMethod.MaximumAmount, paymentMethod.ID).Error

	return err
}
