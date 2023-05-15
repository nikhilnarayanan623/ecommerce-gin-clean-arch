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
	GetRazorpayOrder(ctx context.Context, userID uint, userOrder res.UserOrder) (razorpayOrder res.RazorpayOrder, err error)

	// stipe
	GetStripeOrder(ctx context.Context, userID uint, userOrder res.UserOrder) (stipeOrder res.StripeOrder, err error)

	// order placement
	GetOrderDetails(ctx context.Context, userID uint, body req.PlaceOrder) (userOrder res.UserOrder, err error)
	SaveOrder(ctx context.Context, shopOrder domain.ShopOrder) (shopOrderID uint, err error)
	ApproveOrderAndClearCart(ctx context.Context, userID, shopOrderID, couponID uint) error
	// end

	// get order and orde items
	GetAllShopOrders(ctx context.Context, pagination req.Pagination) (shopOrders []res.ShopOrder, err error)
	GetUserShopOrder(ctx context.Context, userID uint, pagination req.Pagination) ([]res.ShopOrder, error)
	GetOrderItemsByShopOrderID(ctx context.Context, shopOrderID uint, pagination req.Pagination) ([]res.OrderItem, error)

	// cancell order and change order status
	GetAllOrderStatuses(ctx context.Context) (orderStatuses []domain.OrderStatus, err error)
	UpdateOrderStatus(ctx context.Context, shopOrderID, changeStatusID uint) error
	CancellOrder(ctx context.Context, shopOrderID uint) error

	// return and update
	SubmitReturnuest(ctx context.Context, returnDetails req.Return) error
	GetAllPendingOrderReturns(ctx context.Context, pagination req.Pagination) ([]res.OrderReturn, error)
	GetAllOrderReturns(ctx context.Context, pagination req.Pagination) ([]res.OrderReturn, error)
	UpdateReturnDetails(ctx context.Context, updateDetails req.UpdatOrderReturn) error

	// wallet
	GetUserWallet(ctx context.Context, userID uint) (wallet domain.Wallet, err error)
	GetUserWalletTransactions(ctx context.Context, userID uint, pagination req.Pagination) (transactions []domain.Transaction, err error)
}
