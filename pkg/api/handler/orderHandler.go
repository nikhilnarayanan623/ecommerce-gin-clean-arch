package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
)

type OrderHandler struct {
	orderUseCase interfaces.OrderUseCase
}

func NewOrderHandler(orderUseCase interfaces.OrderUseCase) *OrderHandler {
	return &OrderHandler{orderUseCase: orderUseCase}
}

// PlaceOrderByCart godoc
// @summary api for place order for all cartItem
func (c *OrderHandler) PlaceOrderByCart(ctx *gin.Context) {

	userId := helper.GetUserIdFromContext(ctx)

	addressID, err := helper.StringToUint(ctx.Param("address_id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "faild to conver param",
		})
		return
	}

	var shopOrder = domain.ShopOrder{
		UserID:    userId,
		AddressID: addressID,
		COD:       true,
	}

	if err := c.orderUseCase.PlaceOrderByCart(ctx, shopOrder); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "faild to place to order",
			"error":      err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"msg":        "Successfully order placed",
	})

}

func (c *OrderHandler) GetOrderItemsForUser(ctx *gin.Context) {

	shopOrderID, err := helper.StringToUint(ctx.Param("shop_order_id"))
	fmt.Println("hererererererere", shopOrderID)

	orderItems, err := c.orderUseCase.GetOrderItemsByShopOrderID(ctx, shopOrderID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"StatusCode": 500,
			"msg":        "faild to get user order list",
			"error":      err.Error(),
		})
		return
	}

	if orderItems == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"StatusCode": 200,
			"msg":        "User Order list is empty",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode":  200,
		"msg":         "Successfully Order list got",
		"order items": orderItems,
	})
}

func (c *OrderHandler) GetOrdersOfUser(ctx *gin.Context) {

	userId := helper.GetUserIdFromContext(ctx)

	orders, err := c.orderUseCase.GetUserShopOrder(ctx, userId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"StatusCode": 500,
			"msg":        "faild to get user shop order",
			"error":      err.Error(),
		})
		return
	}

	if orders == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"StatusCode": 200,
			"msg":        "User Have no order history",
		})
		return
	}
	// // copy to response
	// var respose []res.ResShopOrder
	// copier.Copy(&respose, orders)

	// for i, time := range orders {

	// 	respose[i].OrderDate = time.OrderDate.Format("2006-January-02 15:04")
	// }

	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"msg":        "Successfully Shop Order List got",
		"orderList":  orders,
	})

}

func (c *OrderHandler) UdateOrderStatus(ctx *gin.Context) {

	var body struct {
		ShopOrderID   uint `json:"shop_order_id" binding:"required"`
		OrderStatusID uint `json:"order_status_id"`
	}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 200,
			"msg":        "faild to bind json",
			"error":      err.Error(),
		})
		return
	}

	//update the order
	if err := c.orderUseCase.ChangeOrderStatus(ctx, body.ShopOrderID, body.OrderStatusID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "faild to upate order status",
			"error":      err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"msg":        "Successfully order status updated",
	})
}

func (c *OrderHandler) CancellOrder(ctx *gin.Context) {

	shopOrderID, err := helper.StringToUint(ctx.Param("shop_order_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "faild on params",
		})
		return
	}

	if err := c.orderUseCase.CancellOrder(ctx, shopOrderID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "faild to cancell order",
			"error":      err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"msg":        "Successfully order cancelled",
	})
}
