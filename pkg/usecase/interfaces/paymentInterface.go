package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
)

type PaymentUseCase interface {
	GetAllPaymentMethods(ctx context.Context) ([]domain.PaymentMethod, error)
	GetPaymentMethodByID(ctx context.Context, paymentMethodID uint) (domain.PaymentMethod, error)
	AddPaymentMethod(ctx context.Context, paymentMethod domain.PaymentMethod) error
	EditPaymentMethod(ctx context.Context, paymentMethod req.PaymentMethodUpdate) error
}
