package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/res"
)

type OrderUseCase interface {

	// pyment
	GetAllPaymentMethods(ctx context.Context) ([]domain.PaymentMethod, error)
	GetPaymentMethodByID(ctx context.Context, paymentMethodID uint) (domain.PaymentMethod, error)

	// order placement
	OrderCheckOut(ctx context.Context, body req.ReqCheckout) (res.ResOrderCheckout, error)
	PlaceOrderCOD(ctx context.Context, checkoutValues res.ResOrderCheckout) error
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
	GetAllPendingOrderReturns(ctx context.Context) ([]domain.OrderReturn, error)
	GetAllOrderReturns(ctx context.Context) ([]domain.OrderReturn, error)
	UpdateReturnRequest(ctx context.Context, body req.ReqUpdatReturnReq) error
}
