package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
)

type OrderRepository interface {

	//!
	SaveShopOrder(ctx context.Context, shopOrder domain.ShopOrder) (domain.ShopOrder, error)

	CheckcartIsValidForOrder(ctx context.Context, userID uint) (cart domain.Cart, err error)

	//FindUserCoupon(ctx context.Context, couponCode string) (domain.UserCoupon, error)
	UpdateCouponUsedForUser(ctx context.Context, userID, couponID uint) error
	ValidateAddressID(ctx context.Context, addressID uint) error

	CartItemToOrderLines(ctx context.Context, userID uint) ([]domain.OrderLine, error)
	SaveOrderLine(ctx context.Context, orderLine domain.OrderLine) error
	DeleteOrderedCartItems(ctx context.Context, userID uint) error
	//!

	//save order and update
	//SaveOrderByCart(ctx context.Context, shopOrder domain.ShopOrder) error //!
	UpdateShopOrderOrderStatus(ctx context.Context, shopOrderID, changeStatusID uint) error

	// shop order order
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
	FindAllOrderReturns(ctx context.Context, onlyPending bool) ([]res.ResOrderReturn, error)
	SaveOrderReturn(ctx context.Context, orderReturn domain.OrderReturn) error
	UpdateOrderReturn(ctx context.Context, body req.ReqUpdatReturnReq) error

	// payments
	FindPaymentMethodByID(ctx context.Context, paymenMethodtID uint) (domain.PaymentMethod, error)
	FindAllPaymentMethods(ctx context.Context) ([]domain.PaymentMethod, error)
}
