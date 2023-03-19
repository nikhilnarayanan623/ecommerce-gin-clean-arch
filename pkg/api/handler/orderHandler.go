package handler

import (
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

func (c *OrderHandler) ListUserOrder(ctx *gin.Context) {
	userId := helper.GetUserIdFromContext(ctx)

	orders, err := c.orderUseCase.GetOrdersListByUserID(ctx, userId)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"StatusCode": 500,
			"msg":        "faild to get user order list",
			"error":      err.Error(),
		})
		return
	}

	if orders == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"StatusCode": 200,
			"msg":        "User Order list is empty",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"msg":        "Successfully Order list got",
		"orderList":  orders,
	})
}
