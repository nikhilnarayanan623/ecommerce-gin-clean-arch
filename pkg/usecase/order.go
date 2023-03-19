package usecase

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/res"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	service "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
)

type OrderUseCase struct {
	userRepo interfaces.OrderRepository
}

func NewOrderUseCase(orderRepo interfaces.OrderRepository) service.OrderUseCase {
	return &OrderUseCase{userRepo: orderRepo}
}

func (c *OrderUseCase) GetOrdersListByUserID(ctx context.Context, userID uint) ([]res.ResOrder, error) {
	return c.userRepo.GetOrdersListByUserID(ctx, userID)
}

func (c *OrderUseCase) PlaceOrderByCart(ctx context.Context, shopOrder domain.ShopOrder) error {

	return c.userRepo.PlaceOrderByCart(ctx, shopOrder)
}
