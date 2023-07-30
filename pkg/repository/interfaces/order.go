package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/response"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/request"
)

type OrderRepository interface {
	Transaction(callBack func(transactionRepo OrderRepository) error) error

	SaveOrderLine(ctx context.Context, orderLine domain.OrderLine) error

	UpdateShopOrderOrderStatus(ctx context.Context, shopOrderID, changeStatusID uint) error
	UpdateShopOrderStatusAndSavePaymentMethod(ctx context.Context, shopOrderID, orderStatusID, paymentID uint) error

	// shop order order
	SaveShopOrder(ctx context.Context, shopOrder domain.ShopOrder) (shopOrderID uint, err error)
	FindShopOrderByShopOrderID(ctx context.Context, shopOrderID uint) (domain.ShopOrder, error)
	FindAllShopOrders(ctx context.Context, pagination request.Pagination) (shopOrders []response.ShopOrder, err error)
	FindAllShopOrdersByUserID(ctx context.Context, userID uint, pagination request.Pagination) ([]response.ShopOrder, error)

	// find shop order items
	FindAllOrdersItemsByShopOrderID(ctx context.Context,
		shopOrderID uint, pagination request.Pagination) (orderItems []response.OrderItem, err error)

	// order status
	FindOrderStatusByShopOrderID(ctx context.Context, shopOrderID uint) (domain.OrderStatus, error)
	FindOrderStatusByID(ctx context.Context, orderStatusID uint) (domain.OrderStatus, error)
	FindOrderStatusByStatus(ctx context.Context, orderStatus domain.OrderStatusType) (domain.OrderStatus, error)
	FindAllOrderStatuses(ctx context.Context) ([]domain.OrderStatus, error)

	//order return
	FindOrderReturnByReturnID(ctx context.Context, orderReturnID uint) (domain.OrderReturn, error)
	FindOrderReturnByShopOrderID(ctx context.Context, shopOrderID uint) (orderReturn domain.OrderReturn, err error)
	FindAllOrderReturns(ctx context.Context, pagination request.Pagination) ([]response.OrderReturn, error)
	FindAllPendingOrderReturns(ctx context.Context, pagination request.Pagination) ([]response.OrderReturn, error)
	SaveOrderReturn(ctx context.Context, orderReturn domain.OrderReturn) error
	UpdateOrderReturn(ctx context.Context, orderReturn domain.OrderReturn) error

	// wallet
	FindWalletByUserID(ctx context.Context, userID uint) (wallet domain.Wallet, err error)
	SaveWallet(ctx context.Context, userID uint) (walletID uint, err error)
	UpdateWallet(ctx context.Context, walletID, updateTotalAmount uint) error
	SaveWalletTransaction(ctx context.Context, walletTrx domain.Transaction) error

	FindWalletTransactions(ctx context.Context, walletID uint,
		pagination request.Pagination) (transaction []domain.Transaction, err error)
}
