package usecase

import (
	"context"
	"errors"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
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

// place an order for user cart
func (c *OrderUseCase) PlaceOrderByCart(ctx context.Context, shopOrder domain.ShopOrder) error {

	return c.orderRepo.SaveOrderByCart(ctx, shopOrder)
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

	// check order status is placed or cancelled
	if orderStatus.Status == "placed" {
		return errors.New("order already placed can't change its status")
	} else if orderStatus.Status == "cancelled" {
		return errors.New("order already cancelled can't change its status")
	}

	//at last update the order status
	return c.orderRepo.UpdateOrderStatus(ctx, shopOrder, changeStatusID)
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

	// check the status is placed or cancelled
	if orderStatus.Status == "cancelled" || orderStatus.Status == "placed" {
		return errors.New("order is already " + orderStatus.Status)
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

	return c.orderRepo.UpdateOrderStatus(ctx, shopOrder, orderStatus.ID)
}
