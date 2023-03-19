package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/res"
)

type OrderUseCase interface {
	PlaceOrderByCart(ctx context.Context, shopOrder domain.ShopOrder) error
	GetOrdersListByUserID(ctx context.Context, userID uint) ([]res.ResOrder, error)
}
