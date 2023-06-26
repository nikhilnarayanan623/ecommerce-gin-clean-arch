package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	interfaces "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase"
	usecaseInterface "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/response"
)

type OrderHandler struct {
	orderUseCase usecaseInterface.OrderUseCase
}

func NewOrderHandler(orderUseCase usecaseInterface.OrderUseCase) interfaces.OrderHandler {
	return &OrderHandler{
		orderUseCase: orderUseCase,
	}
}

// FindAllOrderStatuses godoc
// @Summary Get all order statuses (Admin)
// @Description API for admin to get all available order statuses
// @Security ApiKeyAuth
// @Id FindAllOrderStatuses
// @Tags Admin Orders
// @Router /admin/orders/statuses [get]
// @Success 200 {object} response.Response{} "Successfully retrieved all order statuses"
// @Success 204 {object} response.Response{} "No order statuses found"
// @Failure 500 {object} response.Response{}  "Failed to find all order statuses"
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

	response.SuccessResponse(ctx, 200, "Successfully retrieved all order statuses", orderStatuses)
}

// PlaceOrderOnCOD godoc
// @summary Place order  for COD (User)
// @Description API for user to place order for cash on delivery
// @security ApiKeyAuth
// @tags User Orders
// @id PlaceOrderOnCOD
// @Param address_id formData string true "Address ID"
// @Router /carts/place-order/ [post]
// @Success 200 {object} response.Response{} "successfully order placed"
// @Success 204 {object} response.Response{} "Cart is empty"
// @Failure 400 {object} response.Response{}  "invalid input"
// @Failure 409 {object} response.Response{}  "Can't place order out of stock product on cart"
// @Failure 500 {object} response.Response{}  "Failed to save order"
func (c *OrderHandler) PlaceOrderOnCOD(ctx *gin.Context) {

	addressID, err := request.GetFormValuesAsUint(ctx, "address_id")
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindFormValueMessage, err, nil)
		return
	}

	body := request.PlaceOrder{
		AddressID:   addressID,
		PaymentType: domain.CODPayment,
	}

	userID := utils.GetUserIdFromContext(ctx)

	shopOrderID, err := c.orderUseCase.SaveOrder(ctx, userID, body)

	if err != nil {
		var statusCode int

		switch true {
		case errors.Is(err, usecase.ErrEmptyCart):
			statusCode = http.StatusNoContent
		case errors.Is(err, usecase.ErrOutOfStockOnCart):
			statusCode = http.StatusConflict
		default:
			statusCode = http.StatusInternalServerError
		}
		response.ErrorResponse(ctx, statusCode, "Failed to save order", err, nil)
		return
	}

	// approve the order and clear the user cart
	err = c.orderUseCase.ApproveShopOrderAndClearCart(ctx, userID, shopOrderID)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to approve order and clear cart", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully order placed for cod")
}

// FindUserOrder godoc
// @summary Get user orders (User)
// @description API to get order for user user orders
// @id FindUserOrder
// @tags User Orders
// @Param page_number query int false "Page Number"
// @Param count query int false "Count Of Order"
// @Router /orders [get]
// @Success 200 {object} response.Response{} "Successfully retrieved all user orders"
// @Success 204 {object} response.Response{} "No shop orders for user"
// @Failure 500 {object} response.Response{} "Failed to retrieve all user orders"
func (c *OrderHandler) FindUserOrder(ctx *gin.Context) {

	userId := utils.GetUserIdFromContext(ctx)
	pagination := request.GetPagination(ctx)

	orders, err := c.orderUseCase.FindUserShopOrder(ctx, userId, pagination)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve all user shop orders", err, nil)
		return
	}

	if orders == nil {
		response.SuccessResponse(ctx, http.StatusNoContent, "No shop orders found", nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully retrieved all user orders", orders)
}

// FindAllShopOrders godoc
// @Summary Get all orders (Admin)
// @Description API for admin to get all orders
// @Id FindAllShopOrders
// @Tags Admin Orders
// @Param page_number query int false "Page Number"
// @Param count query int false "Count"
// @Router /admin/orders/all [get]
// @Success 200 {object} response.Response{} "Successfully retrieved all shop orders"
// @Success 204 {object} response.Response{} "No shop order found"
// @Failure 500 {object} response.Response{} "Failed to find all shop orders"
func (c *OrderHandler) FindAllShopOrders(ctx *gin.Context) {

	pagination := request.GetPagination(ctx)

	shopOrders, err := c.orderUseCase.FindAllShopOrders(ctx, pagination)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find all shop orders", err, nil)
		return
	}

	if len(shopOrders) == 0 {
		response.SuccessResponse(ctx, http.StatusNoContent, "No shop order found", nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully retrieved all shop orders", shopOrders)
}

// FindAllOrderItemsUser godoc
// @Summary Get all order items (User)
// @Description API for user to get all order items of a specific order
// @Id FindAllOrderItemsUser
// @Tags User Orders
// @Param shop_order_id path int true "Shop Order ID"
// @Param page_number query int false "Page Number"
// @Param count query int false "Count Of Order"
// @Router /orders/{shop_order_id}/items  [get]
// @Success 200 {object} response.Response{} "Successfully found order items"
// @Failure 500 {object} response.Response{} "Failed to find order items"
func (c *OrderHandler) FindAllOrderItemsUser() func(ctx *gin.Context) {
	return c.findAllOrderItems()
}

// FindAllOrderItemsAdmin godoc
// @Summary Get all order items (Admin)
// @Description API for user to get all order items of a specific order
// @Id FindAllOrderItemsAdmin
// @Tags Admin Orders
// @Param shop_order_id path int true "Shop Order ID"
// @Param page_number query int false "Page Number"
// @Param count query int false "Count"
// @Router /admin/orders/{shop_order_id}/items [get]
// @Success 200 {object} response.Response{} "Successfully found order items"
// @Success 204 {object} response.Response{} "No order items found"
// @Failure 500 {object} response.Response{} "Failed to find order items"
func (c *OrderHandler) FindAllOrderItemsAdmin() func(ctx *gin.Context) {
	return c.findAllOrderItems()
}

func (c *OrderHandler) findAllOrderItems() func(ctx *gin.Context) {

	return func(ctx *gin.Context) {
		shopOrderID, err := request.GetParamAsUint(ctx, "shop_order_id")
		if err != nil {
			response.ErrorResponse(ctx, http.StatusBadRequest, BindParamFailMessage, err, nil)
		}
		pagination := request.GetPagination(ctx)

		orderItems, err := c.orderUseCase.FindOrderItems(ctx, shopOrderID, pagination)

		if err != nil {
			response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find order items", err, nil)
			return
		}

		if orderItems == nil {
			response.SuccessResponse(ctx, http.StatusNoContent, "No order items found", nil)
			return
		}

		response.SuccessResponse(ctx, http.StatusOK, "Successfully found order items", orderItems)
	}
}

// CancelOrder godoc
// @Summary Cancel order (User)
// @Description Api for user to cancel a order
// @Id CancelOrder
// @Tags User Orders
// @Param shop_order_id path int true "Shop Order ID"
// @Router /orders/{shop_order_id}/cancel [post]
// @Success 200 {object} response.Response{} "Successfully order cancelled"
// @Failure 400 {object} response.Response{} "Invalid inputs"
// @Failure 500 {object} response.Response{} "Failed to cancel order"
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

// UpdateOrderStatus godoc
// @Summary Change order status (Admin)
// @Description API for admin to change order status
// @Id UpdateOrderStatus
// @Tags Admin Orders
// @Param input body request.UpdateOrder{} true "input field"
// @Router /admin/orders/ [put]
// @Success 200 {object} response.Response{} "Successfully order status updated"
// @Failure 400 {object} response.Response{} "invalid input"
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

// SubmitReturnRequest godoc
// @Summary Return request
// @Description API for user to request a return for delivered order
// @Id SubmitReturnRequest
// @Tags User Orders
// @Param input body request.Return true "Input Fields"
// @Router /orders/return [post]
// @Success 200 {object} response.Response{} "Successfully return request submitted for order"
// @Failure 400 {object} response.Response{} "invalid input"
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
// @Summary Get all order returns (Admin)
// @Description API for admin to get all order returns
// @Id FindAllOrderReturns
// @Tags Admin Orders
// @Param page_number query int false "Page Number"
// @Param count query int false "Count Of Order"
// @Router /admin/orders/returns [get]
// @Success 200 {object} response.Response{} "Successfully found all order returns"
// @Failure 500 {object} response.Response{} "Failed to find all order returns"
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
// @Summary Get all pending returns (Admin)
// @Description API for admin to get all pending returns
// @Id FindAllPendingReturns
// @Tags Admin Orders
// @Param page_number query int false "Page Number"
// @Param count query int false "Count Of Order"
// @Router /admin/orders/returns/pending [get]
// @Success 200 {object} response.Response{} "Successfully found all pending orders return requests"
// @Failure 500 {object} response.Response{} "Failed to find all pending order return requests"
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
// @summary Change return request status (Admin)
// @description API for admin to change status of return requested orders
// @id UpdateReturnRequest
// @tags Admin Orders
// @Param input body request.UpdateOrderReturn{} true "Input Fields"
// @Router /admin/orders/returns/pending [put]
// @Success 200 {object} response.Response{} "successfully order_response updated"
// @Failure 500 {object} response.Response{} "invalid input"
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
