package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
)

type OrderRepository interface {
	Transaction(callBack func(transactionRepo OrderRepository) error) error

	SaveOrderLine(ctx context.Context, orderLine domain.OrderLine) error

	UpdateShopOrderOrderStatus(ctx context.Context, shopOrderID, changeStatusID uint) error

	// shop order order
	SaveShopOrder(ctx context.Context, shopOrder domain.ShopOrder) (shopOrderID uint, err error)
	FindShopOrderByShopOrderID(ctx context.Context, shopOrderID uint) (domain.ShopOrder, error)
	FindAllShopOrders(ctx context.Context, pagination req.ReqPagination) (shopOrders []res.ResShopOrder, err error)
	FindAllShopOrdersByUserID(ctx context.Context, userID uint, pagination req.ReqPagination) ([]res.ResShopOrder, error)

	// find shop order items
	FindAllOrdersItemsByShopOrderID(ctx context.Context, shopOrderID uint, pagination req.ReqPagination) (orderItems []res.ResOrderItem, err error)

	// order status
	FindOrderStatusByID(ctx context.Context, orderStatusID uint) (domain.OrderStatus, error)
	FindOrderStatusByStatus(ctx context.Context, orderStatus string) (domain.OrderStatus, error)
	FindAllOrderStauses(ctx context.Context) ([]domain.OrderStatus, error)

	//order return
	FindOrderReturnByReturnID(ctx context.Context, orderReturnID uint) (domain.OrderReturn, error)
	FindOrderReturnByShopOrderID(ctx context.Context, shopOrderID uint) (orderReturn domain.OrderReturn, err error)
	FindAllOrderReturns(ctx context.Context, pagination req.ReqPagination) ([]res.ResOrderReturn, error)
	FindAllPendingOrderReturns(ctx context.Context, pagination req.ReqPagination) ([]res.ResOrderReturn, error)
	SaveOrderReturn(ctx context.Context, orderReturn domain.OrderReturn) error
	UpdateOrderReturn(ctx context.Context, orderReturn domain.OrderReturn) error

	// payments
	FindPaymentMethodByID(ctx context.Context, paymenMethodtID uint) (domain.PaymentMethod, error)
	FindPaymentMethodByType(ctx context.Context, paymentType string) (paymentMethod domain.PaymentMethod, err error)
	FindAllPaymentMethods(ctx context.Context) ([]domain.PaymentMethod, error)
	SavePaymentMethod(ctx context.Context, paymentMethod domain.PaymentMethod) (paymentMethodID uint, err error)
	UpdatePaymentMethod(ctx context.Context, paymentMethod domain.PaymentMethod) error

	// wallet
	FindWalletByUserID(ctx context.Context, userID uint) (wallet domain.Wallet, err error)
	SaveWallet(ctx context.Context, userID uint) (walletID uint, err error)
	UpdateWallet(ctx context.Context, walletID, upateTotalAmount uint) error
	SaveWalletTransaction(ctx context.Context, walletTrx domain.Transaction) error

	FindWalletTransactions(ctx context.Context, walletID uint, pagination req.ReqPagination) (transaction []domain.Transaction, err error)
}
