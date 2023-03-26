package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/res"
)

type OrderRepository interface {

	// checkout for order
	CheckOutCart(ctx context.Context, userId uint) (res.ResCheckOut, error)
	//save order and update
	SaveOrderByCart(ctx context.Context, shopOrder domain.ShopOrder) error
	UpdateShopOrderOrderStatus(ctx context.Context, shopOrderID, changeStatusID uint) error

	//find shop order order
	FindAllShopOrders(ctx context.Context) ([]res.ResShopOrder, error)
	FindShopOrderByShopOrderID(ctx context.Context, shopOrderID uint) (domain.ShopOrder, error)
	FindAllShopOrdersByUserID(ctx context.Context, userID uint) ([]res.ResShopOrder, error)

	// find shop order items
	FindAllOrdersItemsByShopOrderID(ctx context.Context, shopOrderID uint) ([]res.ResOrder, error)
	// order status
	FindOrderStatus(ctx context.Context, orderStatus domain.OrderStatus) (domain.OrderStatus, error)
	FindAllOrderStauses(ctx context.Context) ([]domain.OrderStatus, error)

	//order return
	FindOrderReturn(ctx context.Context, orderReturn domain.OrderReturn) (domain.OrderReturn, error)
	FindAllOrderReturns(ctx context.Context, onlyPending bool) ([]domain.OrderReturn, error)
	SaveOrderReturn(ctx context.Context, orderReturn domain.OrderReturn) error
	UpdateOrderReturn(ctx context.Context, body req.ReqUpdatReturnReq) error
}
