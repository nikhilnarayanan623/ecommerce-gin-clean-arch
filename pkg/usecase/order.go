package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/config"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	service "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
)

type OrderUseCase struct {
	orderRepo interfaces.OrderRepository
}

func NewOrderUseCase(orderRepo interfaces.OrderRepository) service.OrderUseCase {
	return &OrderUseCase{orderRepo: orderRepo}
}

// get all order statuses
func (c *OrderUseCase) GetAllOrderStatuses(ctx context.Context) (orderStatuses []domain.OrderStatus, err error) {
	orderStatuses, err = c.orderRepo.FindAllOrderStauses(ctx)
	if err != nil {
		return orderStatuses, err
	}

	return orderStatuses, nil
}

// func to get all shop order
func (c *OrderUseCase) GetAllShopOrders(ctx context.Context, pagination req.ReqPagination) (shopOrders []res.ResShopOrder, err error) {

	// first find all shopOrders
	if shopOrders, err = c.orderRepo.FindAllShopOrders(ctx, pagination); err != nil {
		return shopOrders, err
	}
	return shopOrders, nil
}

// get order items of a spicific order
func (c *OrderUseCase) GetOrderItemsByShopOrderID(ctx context.Context, shopOrderID uint, pagination req.ReqPagination) (orderItems []res.ResOrderItem, err error) {
	//validate the shopOrderId
	shopOdrer, err := c.orderRepo.FindShopOrderByShopOrderID(ctx, shopOrderID)
	if err != nil {
		return orderItems, err
	} else if shopOdrer.ID == 0 {
		return orderItems, errors.New("invalid shopOrder id")
	}
	orderItems, err = c.orderRepo.FindAllOrdersItemsByShopOrderID(ctx, shopOrderID, pagination)
	if err != nil {
		return orderItems, err
	}

	log.Printf("\n\n successfully got all order items with shop_order_id %v \n\n", shopOrderID)
	return orderItems, nil
}

// get all orders of user
func (c *OrderUseCase) GetUserShopOrder(ctx context.Context, userID uint, pagination req.ReqPagination) ([]res.ResShopOrder, error) {
	return c.orderRepo.FindAllShopOrdersByUserID(ctx, userID, pagination)
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

	//check the given changeStatus id is not approve or order placed(like if an order is pending , then won't allow it to return)
	changeOrderStatus, err := c.orderRepo.FindOrderStatus(ctx, domain.OrderStatus{ID: changeStatusID})
	if err != nil {
		return err
	} else if changeOrderStatus.Status == "" {
		return fmt.Errorf("invalid order_status_id %v", changeOrderStatus.ID)
	}

	//  using switch to compare order status and change status  in easy way
	// initially set an common error of all case and direct go that status and check corresponding status is not we want then return
	// otherwise update the status

	err = fmt.Errorf("order status %s can't change to %s ", orderStatus.Status, changeOrderStatus.Status)
	switch orderStatus.Status {
	case "order placed":
		if changeOrderStatus.Status != "order delivered" {
			return err
		}
	case "return requested":
		if changeOrderStatus.Status != "return approved" && changeOrderStatus.Status != "return cancelled" {
			return err
		}
	default: // order status not order placed or not retuen requsted then don't allow to change status
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
	// new only order placed
	//orderStatus.Status != "pending" && orderStatus.Status != "approved" &&
	if orderStatus.Status != "order placed" {
		return fmt.Errorf("order is %s \ncan't cancell the order", orderStatus.Status)
	}

	// if its not then find the cacell orderStatusID
	orderStatus.ID = 0
	orderStatus.Status = "order cancelled"
	orderStatus, err = c.orderRepo.FindOrderStatus(ctx, orderStatus)
	if err != nil {
		return err
	} else if orderStatus.ID == 0 {
		return errors.New("order cancell option is not avaialbe on database")
	}

	return c.orderRepo.UpdateShopOrderOrderStatus(ctx, shopOrder.ID, orderStatus.ID)
}

// to get pending order returns
func (c *OrderUseCase) GetAllPendingOrderReturns(ctx context.Context, pagination req.ReqPagination) (orderReturns []res.ResOrderReturn, err error) {

	return c.orderRepo.FindAllOrderReturns(ctx, true, pagination) // true for only pending
}

// to get all order return
func (c *OrderUseCase) GetAllOrderReturns(ctx context.Context, pagination req.ReqPagination) (orderReturns []res.ResOrderReturn, err error) {

	return c.orderRepo.FindAllOrderReturns(ctx, false, pagination) // false for  not only pending
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

	// check if the order staus not order placed
	if orderStatus.Status != "order delivered" {
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
func (c *OrderUseCase) UpdateReturnRequest(ctx context.Context, body req.ReqUpdatReturnOrder) error {

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
		if changeOrderStatus.Status != "order returned" {
			return err
		}

	default:
		return err
	}

	// update the order return
	err = c.orderRepo.UpdateOrderReturn(ctx, body)
	if err != nil {
		return err
	}

	// check if return request is changed to returned then update wallet of user
	if changeOrderStatus.Status == "order returned" {

		// get the wallet of user
		wallet, err := c.orderRepo.FindWalletByUserID(ctx, shopOrder.UserID)
		if err != nil {
			return err
		} else if wallet.WalletID == 0 { // if user have no wallet then create a wallet
			wallet.WalletID, err = c.orderRepo.SaveWallet(ctx, shopOrder.UserID)
			if err != nil {
				return err
			}
		}
		// create debit payment type
		creditPaymentType := domain.Credit

		// update wallet
		err = c.orderRepo.UpdateWallet(ctx, wallet.WalletID, shopOrder.OrderTotalPrice, creditPaymentType)
		if err != nil {
			return err
		}
	}

	log.Printf("successfully updated order return request for shop_order_id %v", shopOrder.ID)

	return nil
}

// ! place order
func (c *OrderUseCase) GetOrderDetails(ctx context.Context, userID uint, body req.ReqPlaceOrder) (userOrder res.ResUserOrder, err error) {

	// find the payment method_id
	paymentMethod, err := c.orderRepo.FindPaymentMethodByID(ctx, body.PaymentMethodID)
	if err != nil {
		return userOrder, err
	}
	if paymentMethod.PaymentType == "" {
		return userOrder, errors.New("invalid payment_method_id")
	}
	if paymentMethod.BlockStatus {
		return userOrder, errors.New("payment status is blocked use another payment method")
	}

	// validate the address_id
	err = c.orderRepo.ValidateAddressID(ctx, body.AddressID)
	if err != nil {
		return userOrder, err
	}

	// check the cart of user is valid for place order
	cart, err := c.orderRepo.CheckcartIsValidForOrder(ctx, userID)
	if err != nil {
		return userOrder, err
	}

	if cart.TotalPrice == 0 {
		return userOrder, errors.New("there is no product_s in cart")
	}

	userOrder.AmountToPay = cart.TotalPrice - cart.DiscountAmount
	userOrder.Discount = cart.DiscountAmount
	userOrder.CouponID = cart.AppliedCouponID

	log.Printf("successfully order created for user with user_id %v", userID)
	return userOrder, nil
}

func (c *OrderUseCase) GetStripeOrder(ctx context.Context, userID uint, userOrder res.ResUserOrder) (stipeOrder res.StripeOrder, err error) {
	// get user email and phone of user
	emailAnPhone, err := c.orderRepo.GetUserEmailAndPhone(ctx, userID)
	if err != nil {
		return stipeOrder, err
	}

	// create a clent secret for stipe

	clientSecret, err := utils.GenerateStipeClientSecret(userOrder.AmountToPay, emailAnPhone.Email)

	if err != nil {
		return stipeOrder, err
	}

	// setup the userOrder
	stipeOrder.Stripe = true
	stipeOrder.AmountToPay = userOrder.AmountToPay
	stipeOrder.ClientSecret = clientSecret
	stipeOrder.CouponID = userOrder.CouponID
	stipeOrder.PublishableKey = config.GetCofig().StripPublishKey

	return stipeOrder, nil
}

// generate razorpay order
func (c *OrderUseCase) GetRazorpayOrder(ctx context.Context, userID uint, userOrder res.ResUserOrder) (razorpayOrder res.ResRazorpayOrder, err error) {

	// get user email and phone of user
	emailAnPhone, err := c.orderRepo.GetUserEmailAndPhone(ctx, userID)
	if err != nil {
		return razorpayOrder, err
	}

	// generate razorpay order
	//razorpay amount is caluculate on pisa for india so make the actual price into paisa
	razorPayAmount := userOrder.AmountToPay * 100
	razopayOrderId, err := utils.GenerateRazorpayOrder(razorPayAmount, "test reciept")
	if err != nil {
		return razorpayOrder, err
	}

	// set all details on razopay order
	razorpayOrder.AmountToPay = userOrder.AmountToPay
	razorpayOrder.RazorpayAmount = razorPayAmount

	razorpayOrder.RazorpayKey = config.GetCofig().RazorPayKey

	razorpayOrder.UserID = userID
	razorpayOrder.RazorpayOrderID = razopayOrderId
	razorpayOrder.CouponID = userOrder.CouponID

	razorpayOrder.Email = emailAnPhone.Email
	razorpayOrder.Phone = emailAnPhone.Phone

	return razorpayOrder, nil
}

// save order as pending then after vefication change order status to order placed
func (c *OrderUseCase) SaveOrder(ctx context.Context, shopOrder domain.ShopOrder) (shopOrderID uint, err error) {

	//find order status for pending
	orderStatus, err := c.orderRepo.FindOrderStatus(ctx, domain.OrderStatus{Status: "payment pending"})
	if err != nil {
		return 0, err
	} else if orderStatus.ID == 0 {
		return 0, errors.New("order status order pending not found")
	}
	// set the pending order status
	shopOrder.OrderStatusID = orderStatus.ID

	// save shop_order
	shopOrder, err = c.orderRepo.SaveShopOrder(ctx, shopOrder)
	if err != nil {
		return 0, err
	}

	// make orderLines for cart
	ordlerLines, err := c.orderRepo.CartItemToOrderLines(ctx, shopOrder.UserID)
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

// approve the order from payment pending to
func (c *OrderUseCase) ApproveOrderAndClearCart(ctx context.Context, userID, shopOrderID, couponID uint) error {

	//find order status for order order placed
	orderStatus, err := c.orderRepo.FindOrderStatus(ctx, domain.OrderStatus{Status: "order placed"})
	if err != nil {
		return err
	} else if orderStatus.ID == 0 {
		return errors.New("order status order placed not found")
	}

	err = c.orderRepo.UpdateShopOrderOrderStatus(ctx, shopOrderID, orderStatus.ID)
	if err != nil {
		return err
	}

	//update the coupon status as used
	if couponID != 0 {
		err = c.orderRepo.UpdateCouponUsedForUser(ctx, userID, couponID)
		if err != nil {
			return err
		}
	}

	// delete ordered items from cart
	return c.orderRepo.DeleteOrderedCartItems(ctx, userID)
}
