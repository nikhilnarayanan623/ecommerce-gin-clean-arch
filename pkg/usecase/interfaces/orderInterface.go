package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
)

type OrderUseCase interface {

	// pyment
	AddPaymentMethod(ctx context.Context, paymentMethod domain.PaymentMethod) error
	EditPaymentMethod(ctx context.Context, paymentMethod domain.PaymentMethod) error
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
	UpdateOrderStatus(ctx context.Context, shopOrderID, changeStatusID uint) error
	CancellOrder(ctx context.Context, shopOrderID uint) error

	// return and update
	SubmitReturnRequest(ctx context.Context, returnDetails req.ReqReturn) error
	GetAllPendingOrderReturns(ctx context.Context, pagination req.ReqPagination) ([]res.ResOrderReturn, error)
	GetAllOrderReturns(ctx context.Context, pagination req.ReqPagination) ([]res.ResOrderReturn, error)
	UpdateReturnDetails(ctx context.Context, updateDetails req.UpdatOrderReturn) error

	// wallet
	GetUserWallet(ctx context.Context, userID uint) (wallet domain.Wallet, err error)
	GetUserWalletTransactions(ctx context.Context, userID uint, pagination req.ReqPagination) (transactions []domain.Transaction, err error)
}
