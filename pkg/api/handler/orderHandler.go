package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/res"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
)

type OrderHandler struct {
	orderUseCase interfaces.OrderUseCase
	userUseCae   interfaces.UserUseCase
}

func NewOrderHandler(orderUseCase interfaces.OrderUseCase) *OrderHandler {
	return &OrderHandler{orderUseCase: orderUseCase}
}

// CheckOutCart godoc
// @summary api for cart checkout
// @description user can checkout user cart items
// @security ApiKeyAuth
// @id CheckOutCart
// @tags Carts
// @Router /carts/checkout [get]
// @Success 200 {object} res.Response{} "successfully got checkout data"
// @Failure 401 {object} res.Response{} "cart is empty so user can't call this api"
// @Failure 500 {object} res.Response{} "faild to get checkout items"
func (c *OrderHandler) CheckOutCart(ctx *gin.Context) {

	userId := helper.GetUserIdFromContext(ctx)

	resCheckOut, err := c.orderUseCase.CheckOutCart(ctx, userId)

	if err != nil {
		response := res.ErrorResponse(500, "faild to get checkout items", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if resCheckOut.ProductItems == nil {
		response := res.ErrorResponse(401, "cart is empty so user can't call this api", "", nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	responser := res.SuccessResponse(200, "successfully got checkout data", resCheckOut)
	ctx.JSON(http.StatusOK, responser)
}

// PlaceOrderByCart godoc
// @summary api for place order of all items in user cart
// @description user can place after checkout
// @id PlaceOrderByCart
// @tags Carts
// @Router /carts/place-order/:address_id [post]
// @Params address_id path int true "address_id"
// @Success 200 {object} res.Response{} "successfully placed your order for COD"
// @Failure 400 {object} res.Response{} "faild to place to order"
func (c *OrderHandler) PlaceOrderByCart(ctx *gin.Context) {

	userId := helper.GetUserIdFromContext(ctx)

	addressID, err := helper.StringToUint(ctx.Param("address_id"))

	if err != nil {
		response := res.ErrorResponse(400, "invalid input for params", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	var shopOrder = domain.ShopOrder{
		UserID:    userId,
		AddressID: addressID,
		COD:       true,
	}

	if err := c.orderUseCase.PlaceOrderByCart(ctx, shopOrder); err != nil {
		response := res.ErrorResponse(400, "faild to place order", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := res.SuccessResponse(200, "Successfully placed your order for COD", nil)
	ctx.JSON(http.StatusOK, response)

}

// GetUserOrder godoc
// @summary api for showing user order list
// @description user can see all user order history
// @id GetUserOrder
// @tags Orders
// @Router /orders [get]
// @Success 200 {object} res.Response{} "successfully got shop order list of user"
// @Failure 500 {object} res.Response{} "faild to get user shop order list"
func (c *OrderHandler) GetUserOrder(ctx *gin.Context) {

	userId := helper.GetUserIdFromContext(ctx)

	orders, err := c.orderUseCase.GetUserShopOrder(ctx, userId)

	if err != nil {
		response := res.ErrorResponse(500, "faild to get user shop order list", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if orders == nil {
		response := res.SuccessResponse(200, "user have no order history", nil)
		ctx.JSON(http.StatusOK, response)
		return
	}

	response := res.SuccessResponse(200, "successfully got shop order list of user", orders)
	ctx.JSON(http.StatusOK, response)

}

// GetOrderItemsForUser godoc
// @summary api for show order items of a specific order
// @description user can place after checkout
// @id GetOrderItemsForUser
// @tags Orders
// @Params shop_order_id path int true "shop_order_id"
// @Router /orders/items [get]
// @Success 200 {object} res.Response{} "successfully got order items"
// @Failure 500 {object} res.Response{} "faild to get order list of user"
func (c *OrderHandler) GetOrderItemsForUser(ctx *gin.Context) {

	shopOrderID, _ := helper.StringToUint(ctx.Param("shop_order_id"))

	orderItems, err := c.orderUseCase.GetOrderItemsByShopOrderID(ctx, shopOrderID)

	if err != nil {
		response := res.ErrorResponse(500, "faild to get order list of user", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if orderItems == nil {
		response := res.SuccessResponse(200, "user order list is empty", nil)
		ctx.JSON(http.StatusOK, response)
		return
	}

	response := res.SuccessResponse(200, "successfully got order items", orderItems)
	ctx.JSON(http.StatusOK, response)
}

// UdateOrderStatus godoc
// @summary api for admin to change the status of order
// @description admin can change user order status
// @id UdateOrderStatus
// @tags Orders
// @Param input body req.ReqUpdateOrder true "input field"
// @Router /orders/ [put]
// @Success 200 {object} res.Response{} "successfully got order items"
// @Failure 400 {object} res.Response{} "invalid input"
func (c *OrderHandler) UdateOrderStatus(ctx *gin.Context) {

	var body req.ReqUpdateOrder

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response := res.ErrorResponse(400, "invalid input", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	//update the order
	if err := c.orderUseCase.ChangeOrderStatus(ctx, body.ShopOrderID, body.OrderStatusID); err != nil {
		respose := res.ErrorResponse(400, "faild to update order status", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, respose)
		return
	}

	response := res.SuccessResponse(200, "successfully order status updated", nil)
	ctx.JSON(http.StatusOK, response)
}

// CancellOrder godoc
// @summary api for user to cancell the order
// @description user can cancell the order if it's not placed
// @id CancellOrder
// @tags Orders
// @Params shop_order_id path int true "shop_order_id"
// @Router /orders [put]
// @Success 200 {object} res.Response{} "Successfully order cancelled"
// @Failure 400 {object} res.Response{} "invalid input on param"
func (c *OrderHandler) CancellOrder(ctx *gin.Context) {

	shopOrderID, err := helper.StringToUint(ctx.Param("shop_order_id"))
	if err != nil {
		respnose := res.ErrorResponse(400, "invalid input on param", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, respnose)
		return
	}
	err = c.orderUseCase.CancellOrder(ctx, shopOrderID)
	if err != nil {
		respnose := res.ErrorResponse(400, "faild to cancell order", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, respnose)
		return
	}

	response := res.SuccessResponse(200, "successfully order cancelled", nil)
	ctx.JSON(http.StatusOK, response)
}

// GetAllShopOrders godoc
// @summary api for admin to change the status of order
// @description admin can change user order status
// @id GetAllShopOrders
// @tags Orders
// @Router /admin/orders [get]
// @Success 200 {object} res.Response{} "Successfully order cancelled"
// @Failure 400 {object} res.Response{} "invalid input"
func (c *OrderHandler) GetAllShopOrders(ctx *gin.Context) {

	orders, err := c.orderUseCase.GetAllShopOrders(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"StatusCode": 500,
			"msg":        "faild to get all order list",
			"error":      err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"msg":        "Successfully order list got",
		"order list": orders,
	})
}

// ReturnRequest godoc
// @summary api for user to request a return for an order
// @description user can request return for placed orders
// @id ReturnRequest
// @tags Orders
// @Router /orders/return [put]
// @Success 200 {object} res.Response{} "successfully submited return request for order"
// @Failure 400 {object} res.Response{} "invalid input"
func (c OrderHandler) ReturnRequest(ctx *gin.Context) {
	var body req.ReqReturn
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response := res.ErrorResponse(400, "invalid input", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	err := c.orderUseCase.ReturnRequest(ctx, body)
	if err != nil {
		response := res.ErrorResponse(400, "faild to place return request", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := res.SuccessResponse(200, "successfully submited return request for order", nil)
	ctx.JSON(http.StatusOK, response)
}

func (c *OrderHandler) ShowAllPendingReturns(ctx *gin.Context) {

	orderReturns, err := c.orderUseCase.GetAllPendingOrderReturn(ctx)
	if err != nil {
		response := res.ErrorResponse(500, "faild to get pending order return requests", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if orderReturns == nil {
		response := res.SuccessResponse(200, "there is no pendinng order return request", nil)
		ctx.JSON(200, response)
		return
	}

	var responseStruct []res.ResOrderReturn

	copier.Copy(&responseStruct, &orderReturns)

	response := res.SuccessResponse(200, "successfully got  pending request", responseStruct)
	ctx.JSON(http.StatusOK, response)
}
