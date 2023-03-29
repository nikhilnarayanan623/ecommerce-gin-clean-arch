package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/res"
)

type OrderUseCase interface {

	GetAllShopOrders(ctx context.Context) (res.ResShopOrdersPage, error)

	PlaceOrderByCart(ctx context.Context, userID domain.ShopOrder) error
	GetUserShopOrder(ctx context.Context, userID uint) ([]res.ResShopOrder, error)
	GetOrderItemsByShopOrderID(ctx context.Context, shopOrderID uint) ([]res.ResOrder, error)

	ChangeOrderStatus(ctx context.Context, shopOrderID, changeStatusID uint) error
	CancellOrder(ctx context.Context, shopOrderID uint) error

	SubmitReturnRequest(ctx context.Context, body req.ReqReturn) error
	GetAllPendingOrderReturns(ctx context.Context) ([]domain.OrderReturn, error)
	GetAllOrderReturns(ctx context.Context) ([]domain.OrderReturn, error)
	UpdateReturnRequest(ctx context.Context, body req.ReqUpdatReturnReq) error
}
