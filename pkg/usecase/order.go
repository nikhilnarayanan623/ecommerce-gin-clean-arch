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
func (c *OrderUseCase) FindAllOrderStatuses(ctx context.Context) (orderStatuses []domain.OrderStatus, err error) {
	orderStatuses, err = c.orderRepo.FindAllOrderStauses(ctx)
	if err != nil {
		return orderStatuses, err
	}

	return orderStatuses, nil
}

// func to Find all shop order
func (c *OrderUseCase) FindAllShopOrders(ctx context.Context, pagination request.Pagination) (shopOrders []response.ShopOrder, err error) {

	// first find all shopOrders
	if shopOrders, err = c.orderRepo.FindAllShopOrders(ctx, pagination); err != nil {
		return shopOrders, err
	}

	for i, order := range shopOrders {

		if address, err := c.userRepo.FindAddressByID(ctx, order.AddressID); err != nil {
			return shopOrders, fmt.Errorf("faild to Find address")
		} else {
			shopOrders[i].Address = address
		}
	}

	return shopOrders, nil
}

// Find order items of a spicific order
func (c *OrderUseCase) FindOrderItemsByShopOrderID(ctx context.Context, shopOrderID uint,
	pagination request.Pagination) (orderItems []response.OrderItem, err error) {
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

// Find all orders of user
func (c *OrderUseCase) FindUserShopOrder(ctx context.Context, userID uint, pagination request.Pagination) ([]response.ShopOrder, error) {
	return c.orderRepo.FindAllShopOrdersByUserID(ctx, userID, pagination)
}

// update order
func (c *OrderUseCase) UpdateOrderStatus(ctx context.Context, shopOrderID, changeStatusID uint) error {

	shopOrder, err := c.orderRepo.FindShopOrderByShopOrderID(ctx, shopOrderID)
	if err != nil {
		return err
	} else if shopOrder.ID == 0 {
		return errors.New("invalid shopOrderID")
	}

	currentOrderStatus, err := c.orderRepo.FindOrderStatusByID(ctx, shopOrder.OrderStatusID)
	if err != nil {
		return err
	}

	orderStatusChangeTo, err := c.orderRepo.FindOrderStatusByID(ctx, changeStatusID)
	if err != nil {
		return err
	} else if orderStatusChangeTo.Status == "" {
		return fmt.Errorf("invalid order_status_id %v", orderStatusChangeTo.ID)
	}

	switch currentOrderStatus.Status { // switch to add more status in future if need add new status on switch and validate
	case "order placed":
		if orderStatusChangeTo.Status != "order delivered" {
			return fmt.Errorf("order status is 'order placed' \nchange status should be 'order delivered'")
		}
	default:
		return fmt.Errorf("order status %s can't change to %s ", currentOrderStatus.Status, orderStatusChangeTo.Status)
	}

	err = c.orderRepo.UpdateShopOrderOrderStatus(ctx, shopOrder.ID, changeStatusID)
	if err != nil {
		return fmt.Errorf("faild to chnage order status %v", err.Error())
	}
	return nil
}

func (c *OrderUseCase) CancelOrder(ctx context.Context, shopOrderID uint) error {

	shopOrder, err := c.orderRepo.FindShopOrderByShopOrderID(ctx, shopOrderID)
	if err != nil {
		return err
	} else if shopOrder.ID == 0 {
		return errors.New("invalid shopOrderID")
	}

	currentOrderStatus, err := c.orderRepo.FindOrderStatusByID(ctx, shopOrder.OrderStatusID)
	if err != nil {
		return err
	}

	if currentOrderStatus.Status != "order placed" {
		return fmt.Errorf("order is ' %s ' \ncan't cancell the order", currentOrderStatus.Status)
	}

	// if its not then find the cacell orderStatusID
	cancellOrderStatus, err := c.orderRepo.FindOrderStatusByStatus(ctx, "order cancelled")
	if err != nil {
		return err
	} else if cancellOrderStatus.ID == 0 {
		return errors.New("order cancell option is not avaialbe on database")
	}

	err = c.orderRepo.UpdateShopOrderOrderStatus(ctx, shopOrder.ID, cancellOrderStatus.ID)
	if err != nil {
		return fmt.Errorf("faild to cancel the order %v", err.Error())
	}
	log.Printf("successfullu order cancelled for shop order id %v", shopOrder.ID)
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
	} else if shopOrder.ID == 0 {
		return errors.New("invalid shop_order_id")
	}

	currentOrderStatus, err := c.orderRepo.FindOrderStatusByID(ctx, shopOrder.OrderStatusID)
	if err != nil {
		return err
	}

	if currentOrderStatus.Status != "order delivered" {
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
			return fmt.Errorf("faild to submit order return \nerror:%v", err.Error())
		}

		statusToChange, err := trxRepo.FindOrderStatusByStatus(ctx, "return requested")
		if err != nil {
			return fmt.Errorf("faild to find return request status \nerror:%v", err.Error())
		} else if statusToChange.ID == 0 {
			return fmt.Errorf("'return requested' status not found")
		}

		err = trxRepo.UpdateShopOrderOrderStatus(ctx, shopOrder.ID, statusToChange.ID)
		if err != nil {
			return fmt.Errorf("faild to update order status \n error:%v", err.Error())
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("faild to save order return \nerror:%v", err.Error())
	}
	log.Println("successfully order rerturn request submited")
	return nil
}

func (c *OrderUseCase) UpdateReturnDetails(ctx context.Context, updateDetails request.UpdateOrderReturn) error {

	orderReturn, err := c.orderRepo.FindOrderReturnByReturnID(ctx, updateDetails.OrderReturnID)
	if err != nil {
		return fmt.Errorf("faild to Find order \nerror:%v", err.Error())
	} else if orderReturn.ShopOrderID == 0 {
		return errors.New("invalid order_return_id")
	}

	shopOrder, err := c.orderRepo.FindShopOrderByShopOrderID(ctx, orderReturn.ShopOrderID)
	if err != nil {
		return fmt.Errorf("faild to Find order details \nerror:%v", err.Error())
	}

	currentOrderStatus, err := c.orderRepo.FindOrderStatusByID(ctx, shopOrder.OrderStatusID)
	if err != nil {
		return err
	}

	returnStatusChangeTo, err := c.orderRepo.FindOrderStatusByID(ctx, updateDetails.OrderStatusID)
	if err != nil {
		return err
	} else if returnStatusChangeTo.Status == "" {
		return errors.New("invalid order_status_id")
	}

	switch currentOrderStatus.Status {
	case "return requested":
		if returnStatusChangeTo.Status == "return approved" {
			if time.Since(updateDetails.ReturnDate) > 0 {
				return fmt.Errorf("given return date is invalid \nto update 'return approved' return date should be greater than cuurent time")
			}
			orderReturn.ApprovalDate = time.Now()
			orderReturn.IsApproved = true
			orderReturn.ReturnDate = updateDetails.ReturnDate
		} else if returnStatusChangeTo.Status == "return cancelled" {
			// nothing extra update on order return may be in future when adding new statuses
		} else {
			return errors.New("order staus is return requested \nchange status must be return approved or return cancelled")
		}

	case "return approved":
		if returnStatusChangeTo.Status != "order returned" {
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
			return fmt.Errorf("faild to update orders return \nerror:%v", err.Error())
		}

		err = c.orderRepo.UpdateShopOrderOrderStatus(ctx, shopOrder.ID, returnStatusChangeTo.ID)
		if err != nil {
			return fmt.Errorf("faild to update order status \nerror:%v", err.Error())
		}

		if returnStatusChangeTo.Status == "order returned" {
			wallet, err := trxRepo.FindWalletByUserID(ctx, shopOrder.UserID)
			if err != nil {
				return fmt.Errorf("faild to get user wallet for refund amount \nerror:%v", err.Error())
			} else if wallet.ID == 0 {
				wallet.ID, err = c.orderRepo.SaveWallet(ctx, shopOrder.UserID)
				if err != nil {
					return fmt.Errorf("faild to create a wallet for user")
				}
			}

			newWalletTotal := wallet.TotalAmount + shopOrder.OrderTotalPrice
			err = c.orderRepo.UpdateWallet(ctx, wallet.ID, newWalletTotal)
			if err != nil {
				return fmt.Errorf("faild to update return amount to user wallet \nerror:%v", err.Error())
			}

			err = c.orderRepo.SaveWalletTransaction(ctx, domain.Transaction{
				WalletID:        wallet.ID,
				TransactionDate: time.Now(),
				TransactionType: domain.Credit,
				Amount:          shopOrder.OrderTotalPrice,
			})

			if err != nil {
				return fmt.Errorf("faild to save wallet transaction \nerror:%v", err.Error())
			}
		}
		return nil

	})

	if err != nil {
		return fmt.Errorf("faild to update order return \nerror:%v", err.Error())
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
	shopOrderID, paymentMethodID uint) (razorpayOrder response.RazorpayOrder, err error) {

	shopOrder, err := c.orderRepo.FindShopOrderByShopOrderID(ctx, shopOrderID)
	if err != nil {
		return razorpayOrder, fmt.Errorf("faild to get shop order \nerror:%v", err.Error())
	} else if shopOrder.ID == 0 {
		return razorpayOrder, fmt.Errorf("invalid shop_order_id")
	}

	userDetails, err := c.userRepo.FindUserByUserID(ctx, userID)
	if err != nil {
		return razorpayOrder, err
	}

	//razorpay amount is caluculate on pisa for india so make the actual price into paisa
	razorPayAmount := shopOrder.OrderTotalPrice * 100
	razopayOrderId, err := utils.GenerateRazorpayOrder(razorPayAmount, "test reciept")
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

func (c *OrderUseCase) PlaceOrder(ctx context.Context, userID uint,
	placeOrder request.PlaceOrder) (shopOrder domain.ShopOrder, err error) {
	address, err := c.userRepo.FindAddressByID(ctx, placeOrder.AddressID)
	if err != nil {
		return shopOrder, err
	} else if address.ID == 0 {
		return shopOrder, fmt.Errorf("invalid addess id")
	}

	// check the cart of user is valid for place order
	valid, err := c.cartRepo.IsCartValidForOrder(ctx, userID)
	if err != nil {
		return shopOrder, err
	}

	if !valid {
		return shopOrder, fmt.Errorf("cart is not valid for order")
	}

	cart, err := c.cartRepo.FindCartByUserID(ctx, userID)
	if err != nil {
		return shopOrder, err
	}

	if cart.TotalPrice == 0 {
		return shopOrder, errors.New("there is no products in cart")
	}

	pendingOrderStatus, err := c.orderRepo.FindOrderStatusByStatus(ctx, "payment pending")
	if err != nil {
		return shopOrder, err
	} else if pendingOrderStatus.ID == 0 {
		return shopOrder, errors.New("order status order pending not found")
	}

	shopOrder.UserID = userID
	shopOrder.AddressID = address.ID
	shopOrder.OrderTotalPrice = cart.TotalPrice - cart.DiscountAmount
	shopOrder.Discount = cart.DiscountAmount
	shopOrder.OrderDate = time.Now()
	shopOrder.OrderStatusID = pendingOrderStatus.ID

	err = c.orderRepo.Transaction(func(trxRepo interfaces.OrderRepository) error {

		shopOrder.ID, err = trxRepo.SaveShopOrder(ctx, shopOrder)
		if err != nil {
			return fmt.Errorf("faild to save order \nerror:%v", err.Error())
		}

		cart, err := c.cartRepo.FindCartByUserID(ctx, shopOrder.UserID)
		if err != nil {
			return fmt.Errorf("faild to get user cart \nerror:%v", err.Error())
		}

		cartItems, err := c.cartRepo.FindAllCartItemsByCartID(ctx, cart.ID)
		if err != nil {
			return fmt.Errorf("faild to find all cart items \nerror:%v", err.Error())
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
				return fmt.Errorf("faild to save order line \nerror:%v", err.Error())
			}
		}
		return nil
	})
	return shopOrder, nil
}

func (c *OrderUseCase) ApproveShopOrderAndClearCart(ctx context.Context, userID, shopOrderID, paymentID uint) error {

	if err := c.orderRepo.IsShpoOrderIDIsValid(ctx, shopOrderID); err != nil {
		return fmt.Errorf("invalid shop_order_id \nerror:%v", err.Error())
	}

	if orderStatus, err := c.orderRepo.FindOrderStatusByShopOrderID(ctx, shopOrderID); err != nil {
		return fmt.Errorf("faild to get current order status \nerror:%v", err.Error())
	} else if orderStatus.Status != "payment pending" {
		return fmt.Errorf("order status not payment pending can't approve the order ")
	}

	orderPlacedStatus, err := c.orderRepo.FindOrderStatusByStatus(ctx, "order placed")
	if err != nil {
		return err
	} else if orderPlacedStatus.ID == 0 {
		return errors.New("order status order placed not found")
	}

	err = c.orderRepo.Transaction(func(trxRepo interfaces.OrderRepository) error {

		err = trxRepo.UpdateShopOrderStatusAndPaymentID(ctx, shopOrderID, orderPlacedStatus.ID, paymentID)
		if err != nil {
			return fmt.Errorf("faild to approve order error:%v", err.Error())
		}

		cart, err := c.cartRepo.FindCartByUserID(ctx, userID)
		if err != nil {
			return fmt.Errorf("faild to get user cart \nerror:%v", err.Error())
		} else if cart.ID == 0 {
			return fmt.Errorf("user have no cart")
		}

		if cart.AppliedCouponID != 0 { // if user applied a coupon on cart
			err = c.couponRepo.SaveCouponUses(ctx, domain.CouponUses{
				UserID:   userID,
				CouponID: cart.AppliedCouponID,
			})

			if err != nil {
				return fmt.Errorf("faild to update coupon is applied for user \nerror:%v", err.Error())
			}
		}

		err = c.cartRepo.DeleteAllCartItemsByCartID(ctx, cart.ID)
		if err != nil {
			return fmt.Errorf("faild to approve order \nerror:%v", err.Error())
		}
		return nil
	})
	fmt.Println("herererere test")
	return err
}
