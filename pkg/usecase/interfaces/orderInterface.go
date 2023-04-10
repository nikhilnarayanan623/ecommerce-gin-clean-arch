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

	// razorpay
	GetRazorpayOrder(ctx context.Context, userID uint, userOrder res.ResUserOrder) (razorpayOrder res.ResRazorpayOrder, err error)

	// stipe
	GetStripeOrder(ctx context.Context, userID uint, userOrder res.ResUserOrder) (stipeOrder res.StripeOrder, err error)

	// order placement
	GetOrderDetails(ctx context.Context, userID uint, body req.ReqPlaceOrder) (userOrder res.ResUserOrder, err error)
	SaveOrder(ctx context.Context, shopOrder domain.ShopOrder) (shopOrderID uint, err error)
	ApproveOrderAndClearCart(ctx context.Context, userID, shopOrderID, couponID uint) error
	// end

	// get order and orde items
	GetAllShopOrders(ctx context.Context, pagination req.ReqPagination) (shopOrders []res.ResShopOrder, err error)
	GetUserShopOrder(ctx context.Context, userID uint, pagination req.ReqPagination) ([]res.ResShopOrder, error)
	GetOrderItemsByShopOrderID(ctx context.Context, shopOrderID uint, pagination req.ReqPagination) ([]res.ResOrderItem, error)

	// cancell order and change order status
	GetAllOrderStatuses(ctx context.Context) (orderStatuses []domain.OrderStatus, err error)
	ChangeOrderStatus(ctx context.Context, shopOrderID, changeStatusID uint) error
	CancellOrder(ctx context.Context, shopOrderID uint) error

	// return order and updte
	SubmitReturnRequest(ctx context.Context, body req.ReqReturn) error
	GetAllPendingOrderReturns(ctx context.Context, pagination req.ReqPagination) (orderReturns []res.ResOrderReturn, err error)
	GetAllOrderReturns(ctx context.Context, pagination req.ReqPagination) (orderReturns []res.ResOrderReturn, err error)
	UpdateReturnRequest(ctx context.Context, body req.ReqUpdatReturnOrder) error

	// wallet
	GetUserWallet(ctx context.Context, userID uint) (wallet domain.Wallet, err error)

	GetUserWalletTransactions(ctx context.Context, userID uint, pagination req.ReqPagination) (transactions []domain.Transaction, err error)
}
