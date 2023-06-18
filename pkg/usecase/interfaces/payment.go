package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/request"
)

type PaymentUseCase interface {
	FindAllPaymentMethods(ctx context.Context) ([]domain.PaymentMethod, error)
	FindPaymentMethodByID(ctx context.Context, paymentMethodID uint) (domain.PaymentMethod, error)
	SavePaymentMethod(ctx context.Context, paymentMethod domain.PaymentMethod) error
	UpdatePaymentMethod(ctx context.Context, paymentMethod request.PaymentMethodUpdate) error
}
