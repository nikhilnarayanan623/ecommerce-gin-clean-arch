package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/res"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	service "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
)

type OrderUseCase struct {
	orderRepo interfaces.OrderRepository
}

func NewOrderUseCase(orderRepo interfaces.OrderRepository) service.OrderUseCase {
	return &OrderUseCase{orderRepo: orderRepo}
}

// func to get all shop order
func (c *OrderUseCase) GetAllShopOrders(ctx context.Context) (res.ResShopOrdersPage, error) {
	var (
		resShopOrdersPage res.ResShopOrdersPage
		err               error
	)
	// first find all shopOrders
	if resShopOrdersPage.Orders, err = c.orderRepo.FindAllShopOrders(ctx); err != nil {
		return resShopOrdersPage, err
	}

	// then get all  orderStatus
	if resShopOrdersPage.Statuses, err = c.orderRepo.FindAllOrderStauses(ctx); err != nil {
		return resShopOrdersPage, err
	}

	return resShopOrdersPage, nil
}

// get order items of a spicific order
func (c *OrderUseCase) GetOrderItemsByShopOrderID(ctx context.Context, shopOrderID uint) ([]res.ResOrder, error) {
	//validate the shopOrderId
	shopOdrer, err := c.orderRepo.FindShopOrderByShopOrderID(ctx, shopOrderID)
	if err != nil {
		return nil, err
	} else if shopOdrer.ID == 0 {
		return nil, errors.New("invalid shopOrder id")
	}
	return c.orderRepo.FindAllOrdersItemsByShopOrderID(ctx, shopOrderID)
}

// get all orders of user
func (c *OrderUseCase) GetUserShopOrder(ctx context.Context, userID uint) ([]res.ResShopOrder, error) {
	return c.orderRepo.FindAllShopOrdersByUserID(ctx, userID)
}

// update order
func (c *OrderUseCase) ChangeOrderStatus(ctx context.Context, shopOrderID, changeStatusID uint) error {

	// find the shop order by shopOrderID
	shopOrder, err := c.orderRepo.FindShopOrderByShopOrderID(ctx, shopOrderID)
	if err != nil {
		return err
	} else if shopOrder.ID == 0 {
		return errors.New("invalid shopOrderID")
	}

	// find the order status of order using order statusID
	var orderStaus = domain.OrderStatus{ID: shopOrder.OrderStatusID}
	orderStatus, err := c.orderRepo.FindOrderStatus(ctx, orderStaus)
	if err != nil {
		return err
	}

	//check the given changeStatus id is not approve or placed(like if an order is pending , then won't allow it to return)
	changeOrderStatus, err := c.orderRepo.FindOrderStatus(ctx, domain.OrderStatus{ID: changeStatusID})
	if err != nil {
		return err
	}
	//  using switch to compare order status and change status  in easy way
	// initially set an common error of all case and direct go that status and check corresponding status is not we want then return
	// otherwise update the status

	err = fmt.Errorf("order status %s can't change to %s ", orderStatus.Status, changeOrderStatus.Status)
	switch orderStatus.Status {
	case "placed":
		if changeOrderStatus.Status != "delivered" {
			return err
		}
	case "return requested":
		if changeOrderStatus.Status != "return approved" && changeOrderStatus.Status != "return cancelled" {
			return err
		}
	default: // order status not placed or not retuen requsted then don't allow to change status
		return err
	}

	return c.orderRepo.UpdateShopOrderOrderStatus(ctx, shopOrder.ID, changeStatusID)
}

func (c *OrderUseCase) CancellOrder(ctx context.Context, shopOrderID uint) error {

	// find the shop order by shopOrderID
	shopOrder, err := c.orderRepo.FindShopOrderByShopOrderID(ctx, shopOrderID)
	if err != nil {
		return err
	} else if shopOrder.ID == 0 {
		return errors.New("invalid shopOrderID")
	}

	// find the order status of order
	var orderStatus = domain.OrderStatus{ID: shopOrder.OrderStatusID}
	orderStatus, err = c.orderRepo.FindOrderStatus(ctx, orderStatus)
	if err != nil {
		return err
	}

	// check if order is not in pending or approved then don't allow to cancell
	// new only placed
	//orderStatus.Status != "pending" && orderStatus.Status != "approved" &&
	if orderStatus.Status != "placed" {
		return fmt.Errorf("order is %s \ncan't cancell the order", orderStatus.Status)
	}

	// if its not then find the cacell orderStatusID
	orderStatus.ID = 0
	orderStatus.Status = "cancelled"
	orderStatus, err = c.orderRepo.FindOrderStatus(ctx, orderStatus)
	if err != nil {
		return err
	} else if orderStatus.ID == 0 {
		return errors.New("order cancell option is not avaialbe on database")
	}

	return c.orderRepo.UpdateShopOrderOrderStatus(ctx, shopOrder.ID, orderStatus.ID)
}

// to get pending order returns
func (c *OrderUseCase) GetAllPendingOrderReturns(ctx context.Context) ([]res.ResOrderReturn, error) {

	return c.orderRepo.FindAllOrderReturns(ctx, true) // true for only pending
}

// to get all order return
func (c *OrderUseCase) GetAllOrderReturns(ctx context.Context) ([]res.ResOrderReturn, error) {

	return c.orderRepo.FindAllOrderReturns(ctx, false) // false for  not only pending
}

// return request
func (c *OrderUseCase) SubmitReturnRequest(ctx context.Context, body req.ReqReturn) error {

	// validte the shop order id
	shopOrder, err := c.orderRepo.FindShopOrderByShopOrderID(ctx, body.ShopOrderID)
	if err != nil {
		return err
	} else if shopOrder.ID == 0 {
		return errors.New("invalid shop_order_id")
	}

	// check order return time is over

	// find the status of shop order
	orderStatus := domain.OrderStatus{ID: shopOrder.OrderStatusID}
	if orderStatus, err = c.orderRepo.FindOrderStatus(ctx, orderStatus); err != nil {
		return err
	}

	// check if the order staus not placed
	if orderStatus.Status != "delivered" {
		return fmt.Errorf("order is '%s'\ncan't a make return request for this order", orderStatus.Status)
	}

	// then create a new returnOrder for saving
	var OfferReturn = domain.OrderReturn{
		ShopOrderID:  body.ShopOrderID,
		ReturnReason: body.ReturnReason,
		RequestDate:  time.Now(),
		RefundAmount: shopOrder.OrderTotalPrice,
	}
	//save the return request
	return c.orderRepo.SaveOrderReturn(ctx, OfferReturn)
}

// admin to change the update the return request
func (c *OrderUseCase) UpdateReturnRequest(ctx context.Context, body req.ReqUpdatReturnReq) error {

	//validate the order_retun_id
	var orderReturn = domain.OrderReturn{ID: body.OrderReturnID}
	orderReturn, err := c.orderRepo.FindOrderReturn(ctx, orderReturn)
	if err != nil {
		return err
	} else if orderReturn.ShopOrderID == 0 {
		fmt.Print(orderReturn)
		return errors.New("invalid order_return_id")
	}

	// get the shopOrder
	shopOrder, err := c.orderRepo.FindShopOrderByShopOrderID(ctx, orderReturn.ShopOrderID)
	if err != nil {
		return err
	}

	// get the order status
	orderStatus, err := c.orderRepo.FindOrderStatus(ctx, domain.OrderStatus{ID: shopOrder.OrderStatusID})
	if err != nil {
		return err
	}
	// get the change order status
	changeOrderStatus, err := c.orderRepo.FindOrderStatus(ctx, domain.OrderStatus{ID: body.OrderStatusID})
	if err != nil {
		return err
	} else if changeOrderStatus.Status == "" {
		return errors.New("invalid order_status_id")
	}

	// define an error for invalid status change
	err = fmt.Errorf("order return status %s can't change to %s ", orderStatus.Status, changeOrderStatus.Status)

	switch orderStatus.Status {
	case "return requested": // if order status is requsted it can only change into given two or its an error
		if changeOrderStatus.Status != "return approved" && changeOrderStatus.Status != "return cancelled" {
			return errors.Join(err, errors.New(" change status must be return approved or return cancelled"))
		}

	case "return approved":
		if changeOrderStatus.Status != "returned" {
			return err
		}

	default:
		return err
	}

	return c.orderRepo.UpdateOrderReturn(ctx, body)
}

func (c *OrderUseCase) OrderCheckOut(ctx context.Context, body req.ReqCheckout) (res.ResOrderCheckout, error) {

	// find the payment method
	paymentMethod, err := c.orderRepo.FindPaymentMethodByID(ctx, body.PaymentMethodID)
	if err != nil {
		return res.ResOrderCheckout{}, err
	} else if paymentMethod.PaymentType == "" {
		return res.ResOrderCheckout{}, errors.New("invalid payment_method_id")
	} else if paymentMethod.BlockStatus {
		return res.ResOrderCheckout{}, errors.New("payment status is blocked use another payment method")
	}

	// find the total price of cart of user
	cartTotalPrice, err := c.orderRepo.FindCartTotalPrice(ctx, body.UserID)
	if err != nil {
		return res.ResOrderCheckout{}, err
	} else if cartTotalPrice == 0 {
		return res.ResOrderCheckout{}, errors.New("there is no product_items eligible for place order in cart \nadd product to cart")
	}

	// compare payement max_amount with total price
	if cartTotalPrice > paymentMethod.MaximumAmount {
		return res.ResOrderCheckout{}, fmt.Errorf("cart order total price is more than payment max amount %d", paymentMethod.MaximumAmount)
	}

	var discountAmount uint
	// if couponCode exist then check coupon code is valid or not
	if body.CouponCode != "" {
		userCoupon, err := c.orderRepo.FindUserCoupon(ctx, body.CouponCode)
		if err != nil {
			return res.ResOrderCheckout{}, err
		} else if userCoupon.ID == 0 { // validate coupon
			return res.ResOrderCheckout{}, errors.New("invalid coupon code \nplease enter valid coupon code")
		} else if time.Since(userCoupon.ExpireDate) > 0 { // check expire date
			return res.ResOrderCheckout{}, errors.New("can't use coupon code \ncoupon code is expired")
		} else if userCoupon.Used { // check couponCode already used
			return res.ResOrderCheckout{}, errors.New("can't user coupon code \nthis coupon already used")
		} else if userCoupon.DiscountAmount == 0 { // check coupon applied on cart
			return res.ResOrderCheckout{}, errors.New("given coupon_code is not applied apply on cart checkout")
		} else if userCoupon.CartPrice != cartTotalPrice { // check coupon applied time cart price and current cart price
			return res.ResOrderCheckout{}, errors.New("coupon applied cart price and order time cart price not matching")
		}
		discountAmount = userCoupon.DiscountAmount // discount check coupon applied or not
	}

	// validate the address_id
	if err := c.orderRepo.ValidateAddressID(ctx, body.AddressID); err != nil {
		return res.ResOrderCheckout{}, err
	}

	//create a resCheckout and return
	return res.ResOrderCheckout{
		UserID:          body.UserID,
		PaymentMethodID: paymentMethod.ID,
		PaymentType:     paymentMethod.PaymentType,
		AmountToPay:     cartTotalPrice - discountAmount, // subtract discount price on total price
		Discount:        discountAmount,
		AddressID:       body.AddressID,
		CouponCode:      body.CouponCode,
	}, nil
}

// save order as pending then after vefication change order status to placed
func (c *OrderUseCase) SaveOrder(ctx context.Context, checkoutValues res.ResOrderCheckout) (uint, error) {

	//find order status for pending
	orderStatus, err := c.orderRepo.FindOrderStatus(ctx, domain.OrderStatus{Status: "payment pending"})
	if err != nil {
		return 0, err
	} else if orderStatus.ID == 0 {
		return 0, errors.New("order status pending not found")
	}

	shopOrder := domain.ShopOrder{
		UserID:          checkoutValues.UserID,
		AddressID:       checkoutValues.AddressID,
		PaymentMethodID: checkoutValues.PaymentMethodID,
		OrderTotalPrice: checkoutValues.AmountToPay,
		OrderStatusID:   orderStatus.ID,
		Discount:        checkoutValues.Discount,
	}

	// save shop_order
	shopOrder, err = c.orderRepo.SaveShopOrder(ctx, shopOrder)
	if err != nil {
		return 0, err
	}

	// make orderLines for cart
	ordlerLines, err := c.orderRepo.CartItemToOrderLines(ctx, checkoutValues.UserID)
	if err != nil {
		return 0, err
	}

	// save all order lines
	for _, orderLine := range ordlerLines {
		// set shop_order_id
		orderLine.ShopOrderID = shopOrder.ID
		if err := c.orderRepo.SaveOrderLine(ctx, orderLine); err != nil {
			return 0, err
		}
	}

	return shopOrder.ID, nil
}

func (c *OrderUseCase) ApproveOrder(ctx context.Context, userID, shopOrderID uint, couponCode string) error {

	//find order status for placed
	orderStatus, err := c.orderRepo.FindOrderStatus(ctx, domain.OrderStatus{Status: "placed"})
	if err != nil {
		return err
	} else if orderStatus.ID == 0 {
		return errors.New("order status placed not found")
	}

	err = c.orderRepo.UpdateShopOrderOrderStatus(ctx, shopOrderID, orderStatus.ID)
	if err != nil {
		return err
	}

	//update the coupon status as used
	if err := c.orderRepo.UpdteUserCouponAsused(ctx, couponCode); err != nil {
		return err
	}

	// delete ordered items from cart
	return c.orderRepo.DeleteOrderedCartItems(ctx, userID)
}
