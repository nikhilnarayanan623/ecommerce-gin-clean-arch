package usecase

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
)

func (c *OrderUseCase) GetAllPaymentMethods(ctx context.Context) ([]domain.PaymentMethod, error) {
	return c.orderRepo.FindAllPaymentMethods(ctx)
}

func (c *OrderUseCase) GetPaymentMethodByID(ctx context.Context, paymentMethodID uint) (domain.PaymentMethod, error) {
	return c.orderRepo.FindPaymentMethodByID(ctx, paymentMethodID)
}
