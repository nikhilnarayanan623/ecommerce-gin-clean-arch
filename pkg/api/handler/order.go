package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	interfaces "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/interfaces"
	usecase "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/response"
)

type OrderHandler struct {
	orderUseCase usecase.OrderUseCase
}

func NewOrderHandler(orderUseCase usecase.OrderUseCase) interfaces.OrderHandler {
	return &OrderHandler{
		orderUseCase: orderUseCase,
	}
}

// GetAllOrderStatuses godoc
// @summary api for admin to see all order statues for changing order's statuses
// @security ApiKeyAuth
// @tags Admin Orders
// @id GetAllOrderStatuses
// @Router /admin/orders/statuses [get]
// @Success 200 {object} response.Response{} "Successfully found all order statuses"
// @Failure 500 {object} response.Response{}  "failed to get order statuses"
func (c *OrderHandler) FindAllOrderStatuses(ctx *gin.Context) {

	orderStatuses, err := c.orderUseCase.FindAllOrderStatuses(ctx)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, "Failed to find all order statuses", err, nil)
		return
	}

	if orderStatuses == nil {
		response.SuccessResponse(ctx, 200, "No order statuses found")
		return
	}

	response.SuccessResponse(ctx, 200, "Successfully found all order statuses", orderStatuses)
}

// PlaceOrder godoc
// @summary api for user to place an order on cart with COD
// @security ApiKeyAuth
// @tags User Cart
// @id PlaceOrder
// @Param        inputs   body     req.OrderPayment{}   true  "Input Field"
// @Router /carts/place-order/ [post]
// @Success 200 {object} response.Response{} "successfully order placed"
// @Failure 400 {object} response.Response{}  "invalid input"
// @Failure 500 {object} response.Response{}  "failed to save shop order"
func (c *OrderHandler) PlaceOrder(ctx *gin.Context) {

	var body request.PlaceOrder
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	userID := utils.GetUserIdFromContext(ctx)

	shopOrder, err := c.orderUseCase.PlaceOrder(ctx, userID, body)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to save order", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully order placed for payment pending", shopOrder)
}

// ApproveOrderCOD godoc
// @summary api for user to place an order on cart with COD
// @security ApiKeyAuth
// @tags User Cart
// @id ApproveOrderCOD
// @Param       inputs   body     request.OrderPayment{}   true  "Input Field"
// @Router /carts/place-order/cod [post]
// @Success 200 {object} res.Response{} "successfully order placed in COD"
// @Failure 400 {object} res.Response{}  "invalid input"
// @Failure 500 {object} res.Response{}  "failed to save shop order"
func (c *OrderHandler) ApproveOrderCOD(ctx *gin.Context) {

	var body request.OrderPayment

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	userID := utils.GetUserIdFromContext(ctx)

	// approve the order and clear the user cart
	err := c.orderUseCase.ApproveShopOrderAndClearCart(ctx, userID, body.ShopOrderID, body.PaymentMethodID)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to approve order and clear cart", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully order placed for COD", nil)
}

// GetUserOrder godoc
// @summary api for showing User Orders list
// @description user can see all User Orders history
// @id GetUserOrder
// @tags User Orders
// @Param page_number query int false "Page Number"
// @Param count query int false "Count Of Order"
// @Router /orders [get]
// @Success 200 {object} res.Response{} "Successfully found all shop orders"
// @Failure 500 {object} res.Response{} "Failed to find all user shop orders"
func (c *OrderHandler) FindUserOrder(ctx *gin.Context) {

	userId := utils.GetUserIdFromContext(ctx)
	pagination := request.GetPagination(ctx)

	orders, err := c.orderUseCase.FindUserShopOrder(ctx, userId, pagination)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find all user shop orders", err, nil)
		return
	}

	if orders == nil {
		response.SuccessResponse(ctx, http.StatusOK, "No shop orders found", nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully found all user shop orders", orders)
}

// FindAllOrderItems-User godoc
// @summary api for admin to se order items of an order
// @id FindAllOrderItems-Admin
// @tags User Orders
// @Param shop_order_id query int false "Shop Order ID"
// @Param page_number query int false "Page Number"
// @Param count query int false "Count Of Order"
// @Router /admin/orders/items [get]
// @Success 200 {object} res.Response{} "successfully got order items"
// @Failure 500 {object} res.Response{} "faild to get order list of user"
// FindAllOrderItems-User godoc
// @summary api for show order items of a specific order
// @id FindAllOrderItems-User
// @tags User Orders
// @Param shop_order_id query int false "Shop Order ID"
// @Param page_number query int false "Page Number"
// @Param count query int false "Count"
// @Router /orders/items [get]
// @Success 200 {object} res.Response{} "Successfully found order items"
// @Failure 500 {object} res.Response{} "Failed to find order items"
func (c *OrderHandler) FindAllOrderItems(ctx *gin.Context) {

	shopOrderID, err := request.GetParamAsUint(ctx, "shop_order_id")
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindParamFailMessage, err, nil)
	}
	pagination := request.GetPagination(ctx)

	orderItems, err := c.orderUseCase.FindOrderItemsByShopOrderID(ctx, shopOrderID, pagination)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find order items", err, nil)
		return
	}

	if orderItems == nil {
		response.SuccessResponse(ctx, http.StatusOK, "No order items found", nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully found order items", orderItems)
}

// UpdateOrderStatus godoc
// @summary api for admin to change the status of order
// @description admin can change User Orders status
// @id UpdateOrderStatus
// @tags Admin Orders
// @Param input body request.UpdateOrder{} true "input field"
// @Router /admin/orders/ [put]
// @Success 200 {object} res.Response{} "Successfully order status updated"
// @Failure 400 {object} res.Response{} "invalid input"
func (c *OrderHandler) UpdateOrderStatus(ctx *gin.Context) {

	var body request.UpdateOrder

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	err := c.orderUseCase.UpdateOrderStatus(ctx, body.ShopOrderID, body.OrderStatusID)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, "Failed to update order status", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully order status updated", nil)
}

// CancelOrder godoc
// @summary api for user to cancel the order
// @description user can cancel the order if it's not placed
// @id CancelOrder
// @tags User Orders
// @Params shop_order_id path int true "shop_order_id"
// @Router /orders [post]
// @Success 200 {object} res.Response{} "Successfully order cancelled"
// @Failure 400 {object} res.Response{} "invalid input on param"
func (c *OrderHandler) CancelOrder(ctx *gin.Context) {

	shopOrderID, err := request.GetParamAsUint(ctx, "shop_order_id")
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindParamFailMessage, err, nil)
	}

	err = c.orderUseCase.CancelOrder(ctx, shopOrderID)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, "Failed to cancel order", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "successfully order cancelled", nil)
}

// FindAllShopOrders godoc
// @summary api for admin to show all order
// @description admin can see all orders in application
// @id FindAllShopOrders
// @tags Admin Orders
// @Param page_number query int false "Page Number"
// @Param count query int false "Count Of Order"
// @Router /admin/orders [get]
// @Success 200 {object} res.Response{} "Successfully found all shop orders"
// @Failure 500 {object} res.Response{} "Failed to find all shop orders"
func (c *OrderHandler) FindAllShopOrders(ctx *gin.Context) {

	pagination := request.GetPagination(ctx)

	shopOrders, err := c.orderUseCase.FindAllShopOrders(ctx, pagination)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find all shop orders", err, nil)
		return
	}

	if len(shopOrders) == 0 {
		response.SuccessResponse(ctx, http.StatusOK, "No shop order found", nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully found all shop orders", shopOrders)
}

// SubmitReturnRequest godoc
// @summary api for user to request a return for an order
// @description user can request return for placed orders
// @id SubmitReturnRequest
// @tags User Orders
// @Param input body request.Return true "Input Fields"
// @Router /orders/return [post]
// @Success 200 {object} res.Response{} "Successfully return request submitted for order"
// @Failure 400 {object} res.Response{} "invalid input"
func (c OrderHandler) SubmitReturnRequest(ctx *gin.Context) {

	var body request.Return
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	err := c.orderUseCase.SubmitReturnRequest(ctx, body)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, "Failed to submit return request", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully return request submitted for order", nil)
}

// FindAllOrderReturns godoc
// @summary api for admin to see all order returns
// @id FindAllOrderReturns
// @tags Admin Orders
// @Param page_number query int false "Page Number"
// @Param count query int false "Count Of Order"
// @Router /admin/orders/returns [get]
// @Success 200 {object} response.Response{} "Successfully found all order returns"
// @Failure 500 {object} res.Response{} "Failed to find all order returns"
func (c *OrderHandler) FindAllOrderReturns(ctx *gin.Context) {

	pagination := request.GetPagination(ctx)

	orderReturns, err := c.orderUseCase.FindAllOrderReturns(ctx, pagination)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find all order returns", err, nil)
		return
	}

	if len(orderReturns) == 0 {
		response.SuccessResponse(ctx, http.StatusOK, "No order returns found", nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully found all order returns", orderReturns)
}

// FindAllPendingReturns godoc
// @summary api for admin to show pending return request and update it
// @description admin can see the pending return request and accept it or not
// @id FindAllPendingReturns
// @tags Admin Orders
// @Param page_number query int false "Page Number"
// @Param count query int false "Count Of Order"
// @Router /admin/orders/returns/pending [get]
// @Success 200 {object} res.Response{} "Successfully found all pending orders return requests"
// @Failure 500 {object} res.Response{} "Failed to find all pending order return requests"
func (c *OrderHandler) FindAllPendingReturns(ctx *gin.Context) {

	pagination := request.GetPagination(ctx)

	orderReturns, err := c.orderUseCase.FindAllPendingOrderReturns(ctx, pagination)
	if err != nil {
		response.ErrorResponse(ctx, 500, "Failed to find all pending order return requests", err, nil)
		return
	}

	if len(orderReturns) == 0 {
		response.SuccessResponse(ctx, 200, "No pending order returns requests found", nil)
		return
	}

	response.SuccessResponse(ctx, 200, "Successfully found all pending orders return requests", orderReturns)
}

// UpdateReturnRequest godoc
// @summary api for admin to update the order_return request from user
// @description admin can approve, cancel etc. updating on User Orders_return
// @id UpdateReturnRequest
// @tags Admin Orders
// @Param input body req.UpdateOrderReturn{} true "Input Fields"
// @Router /admin/orders/returns/pending [put]
// @Success 200 {object} res.Response{} "successfully order_response updated"
// @Failure 500 {object} res.Response{} "invalid input"
func (c *OrderHandler) UpdateReturnRequest(ctx *gin.Context) {

	var body request.UpdateOrderReturn

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	err := c.orderUseCase.UpdateReturnDetails(ctx, body)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, "Failed to update order return", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "successfully order return updated")
}
