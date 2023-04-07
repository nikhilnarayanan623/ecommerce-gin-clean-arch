package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
)

type OrderUseCase interface {

	// pyment
	GetAllPaymentMethods(ctx context.Context) ([]domain.PaymentMethod, error)
	GetPaymentMethodByID(ctx context.Context, paymentMethodID uint) (domain.PaymentMethod, error)

	// order placement
	GetOrderDetails(ctx context.Context, userID uint, body req.ReqPlaceOrder) (userOrder res.UserOrderCOD, err error)
	SaveOrder(ctx context.Context, checkoutValues domain.ShopOrder) (shopOrderID uint, err error)
	ApproveOrderAndClearCart(ctx context.Context, userID, shopOrderID, couponID uint) error
	// end

	// get order and orde items
	GetAllShopOrders(ctx context.Context) (res.ResShopOrdersPage, error)
	GetUserShopOrder(ctx context.Context, userID uint) ([]res.ResShopOrder, error)
	GetOrderItemsByShopOrderID(ctx context.Context, shopOrderID uint) ([]res.ResOrder, error)

	// cancell order and change order status
	ChangeOrderStatus(ctx context.Context, shopOrderID, changeStatusID uint) error
	CancellOrder(ctx context.Context, shopOrderID uint) error

	// return order and updte
	SubmitReturnRequest(ctx context.Context, body req.ReqReturn) error
	GetAllPendingOrderReturns(ctx context.Context) ([]res.ResOrderReturn, error)
	GetAllOrderReturns(ctx context.Context) ([]res.ResOrderReturn, error)
	UpdateReturnRequest(ctx context.Context, body req.ReqUpdatReturnReq) error
}
