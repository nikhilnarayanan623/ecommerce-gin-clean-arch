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
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/response"
)

type OrderUseCase struct {
	orderRepo   interfaces.OrderRepository
	cartRepo    interfaces.CartRepository
	userRepo    interfaces.UserRepository
	couponRepo  interfaces.CouponRepository
	paymentRepo interfaces.PaymentRepository
}

func NewOrderUseCase(orderRepo interfaces.OrderRepository, cartRepo interfaces.CartRepository,
	userRepo interfaces.UserRepository, couponRepo interfaces.CouponRepository,
	paymentRepo interfaces.PaymentRepository) service.OrderUseCase {
	return &OrderUseCase{
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		userRepo:    userRepo,
		couponRepo:  couponRepo,
		paymentRepo: paymentRepo,
	}
}

// get all order statuses
func (c *OrderUseCase) FindAllOrderStatuses(ctx context.Context) ([]domain.OrderStatus, error) {

	orderStatuses, err := c.orderRepo.FindAllOrderStatuses(ctx)
	if err != nil {
		return nil, utils.PrependMessageToError(err, "failed to find all order statuses")
	}

	return orderStatuses, nil
}

// func to Find all shop order
func (c *OrderUseCase) FindAllShopOrders(ctx context.Context, pagination request.Pagination) ([]response.ShopOrder, error) {

	shopOrders, err := c.orderRepo.FindAllShopOrders(ctx, pagination)
	if err != nil {
		return nil, utils.PrependMessageToError(err, "failed to find all shop orders")
	}

	for i, order := range shopOrders {

		if address, err := c.userRepo.FindAddressByID(ctx, order.AddressID); err != nil {
			return nil, utils.PrependMessageToError(err, "failed to find address for order")
		} else {
			shopOrders[i].Address = address
		}
	}

	return shopOrders, nil
}

func (c *OrderUseCase) FindOrderItems(ctx context.Context, shopOrderID uint,
	pagination request.Pagination) (orderItems []response.OrderItem, err error) {

	orderItems, err = c.orderRepo.FindAllOrdersItemsByShopOrderID(ctx, shopOrderID, pagination)
	if err != nil {
		return nil, utils.PrependMessageToError(err, "failed to find order items using shop order id")
	}

	return orderItems, nil
}

// Find all orders of user
func (c *OrderUseCase) FindUserShopOrder(ctx context.Context, userID uint,
	pagination request.Pagination) ([]response.ShopOrder, error) {

	shopOrders, err := c.orderRepo.FindAllShopOrdersByUserID(ctx, userID, pagination)
	if err != nil {
		return nil, utils.PrependMessageToError(err, "failed to find all shop orders by user id")
	}
	return shopOrders, nil
}

// update order
func (c *OrderUseCase) UpdateOrderStatus(ctx context.Context, shopOrderID, changeStatusID uint) error {

	shopOrder, err := c.orderRepo.FindShopOrderByShopOrderID(ctx, shopOrderID)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to find shop order")
	}

	currentOrderStatus, err := c.orderRepo.FindOrderStatusByID(ctx, shopOrder.OrderStatusID)
	if err != nil {
		return err
	}

	orderStatusChangeTo, err := c.orderRepo.FindOrderStatusByID(ctx, changeStatusID)
	if err != nil {
		return err
	}

	switch currentOrderStatus.Status {

	case domain.StatusOrderPlaced: // if order status is placed then change status should be order delivered
		if orderStatusChangeTo.Status != domain.StatusOrderDelivered {
			return fmt.Errorf("order status is 'order placed' \nchange status should be 'order delivered'")
		}
	default:
		return fmt.Errorf("order status %s can't change to %s ", currentOrderStatus.Status, orderStatusChangeTo.Status)
	}

	err = c.orderRepo.UpdateShopOrderOrderStatus(ctx, shopOrder.ID, changeStatusID)
	if err != nil {
		return fmt.Errorf("failed to change order status %v", err.Error())
	}
	return nil
}

func (c *OrderUseCase) CancelOrder(ctx context.Context, shopOrderID uint) error {

	shopOrder, err := c.orderRepo.FindShopOrderByShopOrderID(ctx, shopOrderID)
	if err != nil {
		return err
	}

	currentOrderStatus, err := c.orderRepo.FindOrderStatusByID(ctx, shopOrder.OrderStatusID)
	if err != nil {
		return err
	}

	if currentOrderStatus.Status != domain.StatusOrderPlaced {
		return fmt.Errorf("order is ' %s ' \ncan't cancel the order", currentOrderStatus.Status)
	}

	// if its not then find the cacel orderStatusID
	cancelOrderStatus, err := c.orderRepo.FindOrderStatusByStatus(ctx, domain.StatusOrderDelivered)
	if err != nil {
		return err
	}

	err = c.orderRepo.UpdateShopOrderOrderStatus(ctx, shopOrder.ID, cancelOrderStatus.ID)
	if err != nil {
		return fmt.Errorf("failed to cancel the order %v", err.Error())
	}
	log.Printf("successfully order cancelled for shop order id %v", shopOrder.ID)
	return nil
}

// to get pending order returns
func (c *OrderUseCase) FindAllPendingOrderReturns(ctx context.Context, pagination request.Pagination) ([]response.OrderReturn, error) {

	pendingOrderReturns, err := c.orderRepo.FindAllPendingOrderReturns(ctx, pagination)
	if err != nil {
		return pendingOrderReturns, fmt.Errorf("failed to Find pendin order returns \nerror:%v", err.Error())
	}
	return pendingOrderReturns, nil
}

// to get all order return
func (c *OrderUseCase) FindAllOrderReturns(ctx context.Context, pagination request.Pagination) ([]response.OrderReturn, error) {

	orderReturns, err := c.orderRepo.FindAllOrderReturns(ctx, pagination)
	if err != nil {
		return orderReturns, fmt.Errorf("faild to Find all order returns \nerror:%v", err.Error())
	}
	return orderReturns, nil
}

func (c *OrderUseCase) SubmitReturnRequest(ctx context.Context, returnDetails request.Return) error {

	shopOrder, err := c.orderRepo.FindShopOrderByShopOrderID(ctx, returnDetails.ShopOrderID)
	if err != nil {
		return err
	}

	currentOrderStatus, err := c.orderRepo.FindOrderStatusByID(ctx, shopOrder.OrderStatusID)
	if err != nil {
		return err
	}

	if currentOrderStatus.Status != domain.StatusOrderDelivered {
		return fmt.Errorf("order is ' %s '\ncan't a make return request for this order", currentOrderStatus.Status)
	}

	orderReturn := domain.OrderReturn{
		ShopOrderID:  returnDetails.ShopOrderID,
		ReturnReason: returnDetails.ReturnReason,
		RequestDate:  time.Now(),
		RefundAmount: shopOrder.OrderTotalPrice,
	}

	err = c.orderRepo.Transaction(func(trxRepo interfaces.OrderRepository) error {

		err := trxRepo.SaveOrderReturn(ctx, orderReturn)
		if err != nil {
			return fmt.Errorf("failed to submit order return \nerror:%v", err.Error())
		}

		statusToChange, err := trxRepo.FindOrderStatusByStatus(ctx, domain.StatusReturnRequested)
		if err != nil {
			return fmt.Errorf("failed to find return request status \nerror:%v", err.Error())
		}

		err = trxRepo.UpdateShopOrderOrderStatus(ctx, shopOrder.ID, statusToChange.ID)
		if err != nil {
			return fmt.Errorf("failed to update order status \n error:%v", err.Error())
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to save order return \nerror:%v", err.Error())
	}
	log.Println("successfully order return request submitted")
	return nil
}

func (c *OrderUseCase) UpdateReturnDetails(ctx context.Context, updateDetails request.UpdateOrderReturn) error {

	orderReturn, err := c.orderRepo.FindOrderReturnByReturnID(ctx, updateDetails.OrderReturnID)
	if err != nil {
		return fmt.Errorf("failed to Find order \nerror:%v", err.Error())
	}

	shopOrder, err := c.orderRepo.FindShopOrderByShopOrderID(ctx, orderReturn.ShopOrderID)
	if err != nil {
		return fmt.Errorf("failed to Find order details \nerror:%v", err.Error())
	}

	currentOrderStatus, err := c.orderRepo.FindOrderStatusByID(ctx, shopOrder.OrderStatusID)
	if err != nil {
		return err
	}

	returnStatusChangeTo, err := c.orderRepo.FindOrderStatusByID(ctx, updateDetails.OrderStatusID)
	if err != nil {
		return err
	}

	switch currentOrderStatus.Status {

	case domain.StatusReturnRequested:
		if returnStatusChangeTo.Status == domain.StatusReturnApproved {
			if time.Since(updateDetails.ReturnDate) > 0 {
				return fmt.Errorf("given return date is invalid \nto update 'return approved' return date should be greater than cuurent time")
			}
			orderReturn.ApprovalDate = time.Now()
			orderReturn.IsApproved = true
			orderReturn.ReturnDate = updateDetails.ReturnDate
		} else if returnStatusChangeTo.Status == domain.StatusReturnCancelled {
			// nothing extra update on order return may be in future when adding new statuses
		} else {
			return errors.New("order staus is return requested \nchange status must be return approved or return cancelled")
		}

	case domain.StatusReturnApproved:
		if returnStatusChangeTo.Status != domain.StatusOrderReturned {
			return errors.New(" change status must be order returned")
		} else if time.Since(updateDetails.ReturnDate) <= 0 {
			return fmt.Errorf("given return date is invalid \nto update 'order returned' return should be lessthan current time")
		} else {
			orderReturn.ReturnDate = updateDetails.ReturnDate
		}

	default:
		return fmt.Errorf("order status %s can't change to %s ", currentOrderStatus.Status, returnStatusChangeTo.Status)
	}

	orderReturn.AdminComment = updateDetails.AdminComment
	err = c.orderRepo.Transaction(func(trxRepo interfaces.OrderRepository) error {

		err := trxRepo.UpdateOrderReturn(ctx, orderReturn)
		if err != nil {
			return fmt.Errorf("failed to update orders return \nerror:%v", err.Error())
		}

		err = c.orderRepo.UpdateShopOrderOrderStatus(ctx, shopOrder.ID, returnStatusChangeTo.ID)
		if err != nil {
			return fmt.Errorf("failed to update order status \nerror:%v", err.Error())
		}

		// if order changing to order return then return the order amount to use wallet
		if returnStatusChangeTo.Status == domain.StatusOrderReturned {
			wallet, err := trxRepo.FindWalletByUserID(ctx, shopOrder.UserID)
			if err != nil {
				return fmt.Errorf("failed to get user wallet for refund amount \nerror:%v", err.Error())
			} else if wallet.ID == 0 {
				wallet.ID, err = c.orderRepo.SaveWallet(ctx, shopOrder.UserID)
				if err != nil {
					return fmt.Errorf("failed to create a wallet for user")
				}
			}

			newWalletTotal := wallet.TotalAmount + shopOrder.OrderTotalPrice
			err = c.orderRepo.UpdateWallet(ctx, wallet.ID, newWalletTotal)
			if err != nil {
				return fmt.Errorf("failed to update return amount to user wallet \nerror:%v", err.Error())
			}

			err = c.orderRepo.SaveWalletTransaction(ctx, domain.Transaction{
				WalletID:        wallet.ID,
				TransactionDate: time.Now(),
				TransactionType: domain.Credit,
				Amount:          shopOrder.OrderTotalPrice,
			})

			if err != nil {
				return fmt.Errorf("failed to save wallet transaction \nerror:%v", err.Error())
			}
		}
		return nil

	})

	if err != nil {
		return fmt.Errorf("failed to update order return \nerror:%v", err.Error())
	}

	log.Printf("successfully updated order return request for shop_order_id %v", shopOrder.ID)
	return nil
}

// ! place order

func (c *OrderUseCase) MakeStripeOrder(ctx context.Context, userID uint,
	userOrder response.UserOrder) (stipeOrder response.StripeOrder, err error) {

	userDetails, err := c.userRepo.FindUserByUserID(ctx, userID)
	if err != nil {
		return stipeOrder, err
	}
	clientSecret, err := utils.GenerateStipeClientSecret(userOrder.AmountToPay, userDetails.Email)
	if err != nil {
		return stipeOrder, err
	}

	stipeOrder.Stripe = true
	stipeOrder.AmountToPay = userOrder.AmountToPay
	stipeOrder.ClientSecret = clientSecret
	stipeOrder.CouponID = userOrder.CouponID
	stipeOrder.PublishableKey = config.GetConfig().StripPublishKey

	return stipeOrder, nil
}

// generate razorpay order
func (c *OrderUseCase) MakeRazorpayOrder(ctx context.Context, userID,
	shopOrderID uint) (razorpayOrder response.RazorpayOrder, err error) {

	shopOrder, err := c.orderRepo.FindShopOrderByShopOrderID(ctx, shopOrderID)
	if err != nil {
		return razorpayOrder, utils.PrependMessageToError(err, "failed to find shop order")
	}

	userDetails, err := c.userRepo.FindUserByUserID(ctx, userID)
	if err != nil {
		return razorpayOrder, err
	}

	//razorpay amount is caluculate on pisa for india so make the actual price into paisa
	razorPayAmount := shopOrder.OrderTotalPrice * 100
	razopayOrderId, err := utils.GenerateRazorpayOrder(razorPayAmount, "test receipt")
	if err != nil {
		return razorpayOrder, err
	}

	// set all details on razopay order
	razorpayOrder.AmountToPay = shopOrder.OrderTotalPrice
	razorpayOrder.RazorpayAmount = razorPayAmount

	razorpayOrder.RazorpayKey = config.GetConfig().RazorPayKey
	razorpayOrder.UserID = userID
	razorpayOrder.RazorpayOrderID = razopayOrderId

	razorpayOrder.Email = userDetails.Email
	razorpayOrder.Phone = userDetails.Phone

	return razorpayOrder, nil
}

// Place order
func (c *OrderUseCase) SaveOrder(ctx context.Context, userID uint,
	orderRequest request.PlaceOrder) (uint, error) {

	// check the cart of user is valid for place order
	valid, err := c.cartRepo.IsCartValidForOrder(ctx, userID)
	if err != nil {
		return 0, utils.PrependMessageToError(err, "failed to check cart is valid for order")
	}

	if !valid {
		return 0, ErrCartIsNotValidForOrder
	}

	cart, err := c.cartRepo.FindCartByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}

	if cart.TotalPrice == 0 {
		return 0, ErrEmptyCart
	}

	pendingOrderStatus, err := c.orderRepo.FindOrderStatusByStatus(ctx, domain.StatusPaymentPending)
	if err != nil {
		return 0, utils.PrependMessageToError(err, "failed to find pending order status")
	}

	payment, err := c.paymentRepo.FindPaymentMethodByType(ctx, orderRequest.PaymentType)
	if err != nil {
		return 0, utils.PrependMessageToError(err, "failed to find payment method details")
	}

	shopOrder := domain.ShopOrder{
		UserID:          userID,
		AddressID:       orderRequest.AddressID,
		OrderTotalPrice: cart.TotalPrice - cart.DiscountAmount,
		Discount:        cart.DiscountAmount,
		OrderDate:       time.Now(),
		OrderStatusID:   pendingOrderStatus.ID,
		PaymentMethodID: payment.ID,
	}

	err = c.orderRepo.Transaction(func(trxRepo interfaces.OrderRepository) error {

		shopOrder.ID, err = trxRepo.SaveShopOrder(ctx, shopOrder)
		if err != nil {
			return fmt.Errorf("failed to save order \nerror:%v", err.Error())
		}

		cart, err := c.cartRepo.FindCartByUserID(ctx, shopOrder.UserID)
		if err != nil {
			return fmt.Errorf("failed to get user cart \nerror:%v", err.Error())
		}

		cartItems, err := c.cartRepo.FindAllCartItemsByCartID(ctx, cart.ID)
		if err != nil {
			return fmt.Errorf("failed to find all cart items \nerror:%v", err.Error())
		}

		var OrderPrice uint
		// save all order lines
		for _, cartItem := range cartItems {

			if cartItem.DiscountPrice != 0 {
				OrderPrice = cartItem.DiscountPrice
			} else {
				OrderPrice = cartItem.Price
			}

			orderLine := domain.OrderLine{
				ProductItemID: cartItem.ProductItemId,
				ShopOrderID:   shopOrder.ID,
				Qty:           cartItem.Qty,
				Price:         OrderPrice,
			}
			if err := trxRepo.SaveOrderLine(ctx, orderLine); err != nil {
				return fmt.Errorf("failed to save order line \nerror:%v", err.Error())
			}
		}
		return nil
	})
	if err != nil {
		return 0, utils.PrependMessageToError(err, "failed to save order")
	}

	return shopOrder.ID, nil
}

func (c *OrderUseCase) ApproveShopOrderAndClearCart(ctx context.Context, userID, shopOrderID uint) error {

	orderStatus, err := c.orderRepo.FindOrderStatusByShopOrderID(ctx, shopOrderID)
	if err != nil {
		return fmt.Errorf("failed to get current order status \nerror:%v", err.Error())
	}
	if orderStatus.Status != domain.StatusPaymentPending {
		return fmt.Errorf("order status not payment pending can't approve the order ")
	}

	orderPlacedStatus, err := c.orderRepo.FindOrderStatusByStatus(ctx, domain.StatusOrderPlaced)
	if err != nil {
		return err
	}

	err = c.orderRepo.Transaction(func(trxRepo interfaces.OrderRepository) error {

		trxRepo.UpdateShopOrderOrderStatus(ctx, shopOrderID, orderPlacedStatus.ID)

		cart, err := c.cartRepo.FindCartByUserID(ctx, userID)
		if err != nil {
			return fmt.Errorf("failed to get user cart \nerror:%v", err.Error())
		} else if cart.ID == 0 {
			return ErrEmptyCart
		}

		// if user applied a coupon on cart then save coupon uses for user
		if cart.AppliedCouponID != 0 {
			err = c.couponRepo.SaveCouponUses(ctx, domain.CouponUses{
				UserID:   userID,
				CouponID: cart.AppliedCouponID,
			})

			if err != nil {
				return fmt.Errorf("failed to update coupon is applied for user \nerror:%v", err.Error())
			}
		}

		err = c.cartRepo.DeleteAllCartItemsByCartID(ctx, cart.ID)
		if err != nil {
			return fmt.Errorf("failed to approve order \nerror:%v", err.Error())
		}
		return nil
	})
	return err
}
