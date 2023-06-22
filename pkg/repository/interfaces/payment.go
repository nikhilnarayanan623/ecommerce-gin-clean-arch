package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/request"
)

type PaymentRepository interface {
	FindPaymentMethodByID(ctx context.Context, paymenMethodtID uint) (paymentMethods domain.PaymentMethod, err error)
	FindPaymentMethodByType(ctx context.Context, paymentType string) (paymentMethod domain.PaymentMethod, err error)
	FindAllPaymentMethods(ctx context.Context) ([]domain.PaymentMethod, error)
	SavePaymentMethod(ctx context.Context, paymentMethod domain.PaymentMethod) (paymentMethodID uint, err error)
	UpdatePaymentMethod(ctx context.Context, paymentMethod request.PaymentMethodUpdate) error
}
