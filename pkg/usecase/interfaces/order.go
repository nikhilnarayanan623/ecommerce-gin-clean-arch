package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/response"
)

type OrderUseCase interface {

	//
	PlaceOrder(ctx context.Context, userID uint, placeOrder request.PlaceOrder) (shopOrder domain.ShopOrder, err error)
	ApproveShopOrderAndClearCart(ctx context.Context, userID, shopOrderID, paymentID uint) error

	// razorpay
	MakeRazorpayOrder(ctx context.Context, userID, shopOrderID, paymentMethodID uint) (razorpayOrder response.RazorpayOrder, err error)
	// stipe
	MakeStripeOrder(ctx context.Context, userID uint, userOrder response.UserOrder) (stipeOrder response.StripeOrder, err error)

	// Find order and order items
	FindAllShopOrders(ctx context.Context, pagination request.Pagination) (shopOrders []response.ShopOrder, err error)
	FindUserShopOrder(ctx context.Context, userID uint, pagination request.Pagination) ([]response.ShopOrder, error)
	FindOrderItemsByShopOrderID(ctx context.Context, shopOrderID uint, pagination request.Pagination) ([]response.OrderItem, error)

	// cancel order and change order status
	FindAllOrderStatuses(ctx context.Context) (orderStatuses []domain.OrderStatus, err error)
	UpdateOrderStatus(ctx context.Context, shopOrderID, changeStatusID uint) error
	CancelOrder(ctx context.Context, shopOrderID uint) error

	// return and update
	SubmitReturnRequest(ctx context.Context, returnDetails request.Return) error
	FindAllPendingOrderReturns(ctx context.Context, pagination request.Pagination) ([]response.OrderReturn, error)
	FindAllOrderReturns(ctx context.Context, pagination request.Pagination) ([]response.OrderReturn, error)
	UpdateReturnDetails(ctx context.Context, updateDetails request.UpdateOrderReturn) error

	// wallet
	FindUserWallet(ctx context.Context, userID uint) (wallet domain.Wallet, err error)
	FindUserWalletTransactions(ctx context.Context, userID uint, pagination request.Pagination) (transactions []domain.Transaction, err error)
}
