package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/res"
)

type OrderRepository interface {
	//save order and update
	SaveOrderByCart(ctx context.Context, shopOrder domain.ShopOrder) error
	UpdateOrderStatus(ctx context.Context, shopOrder domain.ShopOrder, changeStatusID uint) error

	//find order
	FindShopOrderByShopOrderID(ctx context.Context, shopOrderID uint) (domain.ShopOrder, error)

	// find all order an order items
	FindAllShopOrdersByUserID(ctx context.Context, userID uint) ([]res.ResShopOrder, error)
	FindAllOrdersItemsByShopOrderID(ctx context.Context, shopOrderID uint) ([]res.ResOrder, error)
	// order status
	FindOrderStatus(ctx context.Context, orderStatus domain.OrderStatus) (domain.OrderStatus, error)
}
