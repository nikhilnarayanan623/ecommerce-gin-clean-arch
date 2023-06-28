package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	service "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/response"
)

type OrderUseCase struct {
	orderRepo interfaces.OrderRepository
	cartRepo  interfaces.CartRepository
	userRepo  interfaces.UserRepository
}

func NewOrderUseCase(orderRepo interfaces.OrderRepository, cartRepo interfaces.CartRepository,
	userRepo interfaces.UserRepository,
	paymentRepo interfaces.PaymentRepository) service.OrderUseCase {
	return &OrderUseCase{
		orderRepo: orderRepo,
		cartRepo:  cartRepo,
		userRepo:  userRepo,
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

// Save order
func (c *OrderUseCase) SaveOrder(ctx context.Context, userID, addressID uint) (uint, error) {

	cart, err := c.cartRepo.FindCartByUserID(ctx, userID)
	if err != nil {
		return 0, utils.PrependMessageToError(err, "failed to get user cart")
	}

	if cart.TotalPrice == 0 {
		return 0, ErrEmptyCart
	}

	// check the cart of user is valid for place order
	valid, err := c.cartRepo.IsCartValidForOrder(ctx, userID)
	if err != nil {
		return 0, utils.PrependMessageToError(err, "failed to check cart is valid for order")
	}

	if !valid {
		return 0, ErrOutOfStockOnCart
	}

	pendingOrderStatus, err := c.orderRepo.FindOrderStatusByStatus(ctx, domain.StatusPaymentPending)
	if err != nil {
		return 0, utils.PrependMessageToError(err, "failed to find pending order status")
	}

	orderTotal := cart.TotalPrice - cart.DiscountAmount

	shopOrder := domain.ShopOrder{
		UserID:          userID,
		AddressID:       addressID,
		OrderTotalPrice: orderTotal,
		Discount:        cart.DiscountAmount,
		OrderStatusID:   pendingOrderStatus.ID,
	}

	err = c.orderRepo.Transaction(func(trxRepo interfaces.OrderRepository) error {

		shopOrder.ID, err = trxRepo.SaveShopOrder(ctx, shopOrder)
		if err != nil {
			return utils.PrependMessageToError(err, "failed to save shop order on database")
		}

		cartItems, err := c.cartRepo.FindAllCartItemsByCartID(ctx, cart.ID)
		if err != nil {
			return utils.PrependMessageToError(err, "failed to find all cart items")
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
			err = trxRepo.SaveOrderLine(ctx, orderLine)
			if err != nil {
				return utils.PrependMessageToError(err, "failed to save order line on database")
			}
		}
		return nil
	})
	if err != nil {
		return 0, utils.PrependMessageToError(err, "failed to complete save order")
	}

	return shopOrder.ID, nil
}

// Find all orders of a user
func (c *OrderUseCase) FindUserShopOrder(ctx context.Context, userID uint,
	pagination request.Pagination) ([]response.ShopOrder, error) {

	shopOrders, err := c.orderRepo.FindAllShopOrdersByUserID(ctx, userID, pagination)
	if err != nil {
		return nil, utils.PrependMessageToError(err, "failed to find all shop orders by user id")
	}

	for i, order := range shopOrders {

		address, err := c.userRepo.FindAddressByID(ctx, order.AddressID)
		if err != nil {
			return nil, utils.PrependMessageToError(err, "failed to get order address")
		}
		shopOrders[i].Address = address
	}

	return shopOrders, nil
}

// func to Find all shop order
func (c *OrderUseCase) FindAllShopOrders(ctx context.Context, pagination request.Pagination) ([]response.ShopOrder, error) {

	shopOrders, err := c.orderRepo.FindAllShopOrders(ctx, pagination)
	if err != nil {
		return nil, utils.PrependMessageToError(err, "failed to find all shop orders")
	}

	for i, order := range shopOrders {

		address, err := c.userRepo.FindAddressByID(ctx, order.AddressID)
		if err != nil {
			return nil, utils.PrependMessageToError(err, "failed to get order address")
		}
		shopOrders[i].Address = address
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
	cancelOrderStatus, err := c.orderRepo.FindOrderStatusByStatus(ctx, domain.StatusOrderCancelled)
	if err != nil {
		return err
	}

	err = c.orderRepo.UpdateShopOrderOrderStatus(ctx, shopOrder.ID, cancelOrderStatus.ID)
	if err != nil {
		return fmt.Errorf("failed to cancel the order %v", err.Error())
	}

	return nil
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
		return orderReturns, fmt.Errorf("failed to Find all order returns \nerror:%v", err.Error())
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
		}
		if time.Since(updateDetails.ReturnDate) <= 0 {
			return fmt.Errorf("given return date is invalid \nto update 'order returned' return should be less than current time")
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

		err = trxRepo.UpdateShopOrderOrderStatus(ctx, shopOrder.ID, returnStatusChangeTo.ID)
		if err != nil {
			return fmt.Errorf("failed to update order status \nerror:%v", err.Error())
		}

		// if order changing to order return then return the order amount to use wallet
		if returnStatusChangeTo.Status == domain.StatusOrderReturned {
			// get user wallet
			wallet, err := trxRepo.FindWalletByUserID(ctx, shopOrder.UserID)
			if err != nil {
				return fmt.Errorf("failed to get user wallet for refund amount \nerror:%v", err.Error())
			}
			// if user have no wallet then create a new wallet for user
			if wallet.ID == 0 {
				wallet.ID, err = c.orderRepo.SaveWallet(ctx, shopOrder.UserID)
				if err != nil {
					return fmt.Errorf("failed to create a wallet for user")
				}
			}

			// calculate wallet amount and update
			newWalletTotal := wallet.TotalAmount + shopOrder.OrderTotalPrice
			err = trxRepo.UpdateWallet(ctx, wallet.ID, newWalletTotal)
			if err != nil {
				return fmt.Errorf("failed to update return amount to user wallet \nerror:%v", err.Error())
			}

			// wallet transaction
			transaction := domain.Transaction{
				WalletID:        wallet.ID,
				TransactionDate: time.Now(),
				TransactionType: domain.Credit,
				Amount:          shopOrder.OrderTotalPrice,
			}
			err = trxRepo.SaveWalletTransaction(ctx, transaction)

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
