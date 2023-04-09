package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
)

type OrderHandler struct {
	orderUseCase interfaces.OrderUseCase
}

func NewOrderHandler(orderUseCase interfaces.OrderUseCase) *OrderHandler {
	return &OrderHandler{orderUseCase: orderUseCase}
}

// GetAllOrderStatuses godoc
// @summary api for admin to see all order statues for changing order's statuses
// @security ApiKeyAuth
// @tags Admin Orders
// @id GetAllOrderStatuses
// @Router /admin/orders/statuses [get]
// @Success 200 {object} res.Response{} "successfully got all order statueses"
// @Failure 500 {object} res.Response{}  "faild to get order statuses"
func (c *OrderHandler) GetAllOrderStatuses(ctx *gin.Context) {

	orderStatuses, err := c.orderUseCase.GetAllOrderStatuses(ctx)

	if err != nil {
		response := res.ErrorResponse(500, "faild to get order statuses", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := res.SuccessResponse(200, "successfully got all order statueses", orderStatuses)
	ctx.JSON(http.StatusOK, response)
}

func (c *OrderHandler) CartOrderPayementSelectPage(ctx *gin.Context) {

	Payments, err := c.orderUseCase.GetAllPaymentMethods(ctx)
	if err != nil {
		ctx.HTML(200, "paymentForm.html", nil)
		return
	}

	ctx.HTML(200, "paymentForm.html", Payments)
}

// PlaceOrderCartCOD godoc
// @summary api for user to place an order on cart with COD
// @security ApiKeyAuth
// @tags User Cart
// @id PlaceOrderCartCOD
// @Param        inputs   body     req.ReqPlaceOrder{}   true  "Input Field"
// @Router /carts/place-order/cod [post]
// @Success 200 {object} res.Response{} "successfully order placed in COD"
// @Failure 400 {object} res.Response{}  "invalid input"
// @Failure 500 {object} res.Response{}  "faild to save shop order"
func (c *OrderHandler) PlaceOrderCartCOD(ctx *gin.Context) {

	var body req.ReqPlaceOrder
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response := res.ErrorResponse(400, "invalid input", err.Error(), body)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	userID := utils.GetUserIdFromContext(ctx)

	//get the payment method of given payment_id and validate it COD or not
	paymentMethod, err := c.orderUseCase.GetPaymentMethodByID(ctx, body.PaymentMethodID)
	if err != nil {
		response := res.ErrorResponse(400, "faild to place order on COD", err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	// check payment type is  razorpay or not
	if paymentMethod.PaymentType != "COD" {
		respones := res.ErrorResponse(400, "can't place order order", "selected payment_method_id is not for COD ", nil)
		ctx.AbortWithStatusJSON(400, respones)
		return
	}

	// get order details of user
	userOrder, err := c.orderUseCase.GetOrderDetails(ctx, userID, body)
	if err != nil {
		response := res.ErrorResponse(400, "faild to place order", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	// save shopOrder
	// make a shopOrder
	shopOrder := domain.ShopOrder{
		UserID:          userID,
		PaymentMethodID: body.PaymentMethodID,
		AddressID:       body.AddressID,
		OrderTotalPrice: userOrder.AmountToPay,
		Discount:        userOrder.Discount,
		OrderDate:       time.Now(),
	}

	// save order details
	shopOrderID, err := c.orderUseCase.SaveOrder(ctx, shopOrder)
	if err != nil {
		shopOrder.ID = shopOrderID
		response := res.ErrorResponse(500, "faild to save shop order", err.Error(), shopOrder)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	// approve the order and clear the user cart
	err = c.orderUseCase.ApproveOrderAndClearCart(ctx, userID, shopOrderID, userOrder.CouponID)

	if err != nil {
		respnose := res.ErrorResponse(500, "faild to update approve order and clear cart", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, respnose)
		return
	}

	response := res.SuccessResponse(200, "successfully order placed in COD", nil)
	ctx.JSON(http.StatusOK, response)
}

// GetUserOrder godoc
// @summary api for showing User Orders list
// @description user can see all User Orders history
// @id GetUserOrder
// @tags User Orders
// @Param page_number query int false "Page Number"
// @Param count query int false "Count Of Order"
// @Router /orders [get]
// @Success 200 {object} res.Response{} "successfully got shop order list of user"
// @Failure 500 {object} res.Response{} "faild to get user shop order list"
func (c *OrderHandler) GetUserOrder(ctx *gin.Context) {

	userId := utils.GetUserIdFromContext(ctx)
	count, err1 := utils.StringToUint(ctx.Query("count"))
	pageNumber, err2 := utils.StringToUint(ctx.Query("page_number"))

	err1 = errors.Join(err1, err2)
	if err1 != nil {
		response := res.ErrorResponse(400, "invalid inputs", err1.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	pagination := req.ReqPagination{
		PageNumber: pageNumber,
		Count:      count,
	}

	orders, err := c.orderUseCase.GetUserShopOrder(ctx, userId, pagination)

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

// GetOrderItemsByShopOrderItems-User godoc
// @summary api for admin to se order items of an order
// @id GetOrderItemsByShopOrderItems-Admin
// @tags User Orders
// @Param shop_order_id query int false "Shop Order ID"
// @Param page_number query int false "Page Number"
// @Param count query int false "Count Of Order"
// @Router /admin/orders/items [get]
// @Success 200 {object} res.Response{} "successfully got order items"
// @Failure 500 {object} res.Response{} "faild to get order list of user"

// GetOrderItemsByShopOrderItems-User godoc
// @summary api for show order items of a specific order
// @id GetOrderItemsByShopOrderItems-User
// @tags User Orders
// @Param shop_order_id query int false "Shop Order ID"
// @Param page_number query int false "Page Number"
// @Param count query int false "Count Of Order"
// @Router /orders/items [get]
// @Success 200 {object} res.Response{} "successfully got order items"
// @Failure 500 {object} res.Response{} "faild to get order list of user"
func (c *OrderHandler) GetOrderItemsByShopOrderItems(ctx *gin.Context) {

	shopOrderID, err1 := utils.StringToUint(ctx.Query("shop_order_id"))
	count, err2 := utils.StringToUint(ctx.Query("count"))
	pageNumber, err3 := utils.StringToUint(ctx.Query("page_number"))

	err1 = errors.Join(err1, err2, err3)

	if err1 != nil {
		response := res.ErrorResponse(400, "invalid inputs", err1.Error(), nil)
		ctx.JSON(400, response)
		return
	}

	pagination := req.ReqPagination{
		PageNumber: pageNumber,
		Count:      count,
	}

	orderItems, err := c.orderUseCase.GetOrderItemsByShopOrderID(ctx, shopOrderID, pagination)

	if err != nil {
		response := res.ErrorResponse(500, "faild to get order list of user", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if orderItems == nil {
		response := res.SuccessResponse(200, "User Orders list is empty", nil)
		ctx.JSON(http.StatusOK, response)
		return
	}

	response := res.SuccessResponse(200, "successfully got order items", orderItems)
	ctx.JSON(http.StatusOK, response)
}

// UdateOrderStatus godoc
// @summary api for admin to change the status of order
// @description admin can change User Orders status
// @id UdateOrderStatus
// @tags Admin Orders
// @Param input body req.ReqUpdateOrder true "input field"
// @Router /admin/orders/ [put]
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
// @tags User Orders
// @Params shop_order_id path int true "shop_order_id"
// @Router /orders [post]
// @Success 200 {object} res.Response{} "Successfully order cancelled"
// @Failure 400 {object} res.Response{} "invalid input on param"
func (c *OrderHandler) CancellOrder(ctx *gin.Context) {

	shopOrderID, err := utils.StringToUint(ctx.Param("shop_order_id"))
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
// @summary api for admin to show all order
// @description admin can see all orders in application
// @id GetAllShopOrders
// @tags Admin Orders
// @Param page_number query int false "Page Number"
// @Param count query int false "Count Of Order"
// @Router /admin/orders [get]
// @Success 200 {object} res.Response{} "successfully order list got"
// @Failure 500 {object} res.Response{} "faild to get shop order data"
func (c *OrderHandler) GetAllShopOrders(ctx *gin.Context) {

	pageNumber, err1 := utils.StringToUint(ctx.Query("page_number"))
	count, err2 := utils.StringToUint(ctx.Query("count"))

	err1 = errors.Join(err1, err2)

	if err1 != nil {
		response := res.ErrorResponse(400, "invalid inputs", err1.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	pagination := req.ReqPagination{
		PageNumber: pageNumber,
		Count:      count,
	}

	shopOrders, err := c.orderUseCase.GetAllShopOrders(ctx, pagination)
	if err != nil {
		response := res.ErrorResponse(500, "faild to get shop order data", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if shopOrders == nil {
		response := res.SuccessResponse(200, "order list is empty", nil)
		ctx.JSON(200, response)
		return
	}

	response := res.SuccessResponse(200, "successfully order list got", shopOrders)
	ctx.JSON(http.StatusOK, response)
}

// SubmitReturnRequest godoc
// @summary api for user to request a return for an order
// @description user can request return for placed orders
// @id SubmitReturnRequest
// @tags User Orders
// @Param input body req.ReqReturn true "Input Fields"
// @Router /orders/return [post]
// @Success 200 {object} res.Response{} "successfully submited return request for order"
// @Failure 400 {object} res.Response{} "invalid input"
func (c OrderHandler) SubmitReturnRequest(ctx *gin.Context) {
	var body req.ReqReturn
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response := res.ErrorResponse(400, "invalid input", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	err := c.orderUseCase.SubmitReturnRequest(ctx, body)
	if err != nil {
		response := res.ErrorResponse(400, "faild to place return request", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := res.SuccessResponse(200, "successfully submited return request for order", nil)
	ctx.JSON(http.StatusOK, response)
}

// GetAllOrderReturns godoc
// @summary api for admin to see all order reutns
// @id GetAllOrderReturns
// @tags Admin Orders
// @Param page_number query int false "Page Number"
// @Param count query int false "Count Of Order"
// @Router /admin/orders/returns [get]
// @Success 200 {object} res.Response{} "successfully got all order returns"
// @Failure 500 {object} res.Response{} "faild to get order returns"
func (c *OrderHandler) GetAllOrderReturns(ctx *gin.Context) {

	pageNumber, err1 := utils.StringToUint(ctx.Query("page_number"))
	count, err2 := utils.StringToUint(ctx.Query("count"))

	err1 = errors.Join(err1, err2)

	if err1 != nil {
		response := res.ErrorResponse(400, "invalid inputs", err1.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	pagination := req.ReqPagination{
		PageNumber: pageNumber,
		Count:      count,
	}

	orderReturns, err := c.orderUseCase.GetAllOrderReturns(ctx, pagination)
	if err != nil {
		response := res.ErrorResponse(500, "faild to get order returns", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if orderReturns == nil {
		response := res.SuccessResponse(200, "there is no order returns", nil)
		ctx.JSON(http.StatusOK, response)
		return
	}

	response := res.SuccessResponse(200, "successfully got all order returns", orderReturns)
	ctx.JSON(http.StatusOK, response)
}

// GetAllPendingReturns godoc
// @summary api for admin to show pending return request and update it
// @description admin can see the pending return request and accept it or not
// @id GetAllPendingReturns
// @tags Admin Orders
// @Param page_number query int false "Page Number"
// @Param count query int false "Count Of Order"
// @Router /admin/orders/returns/pending [get]
// @Success 200 {object} res.Response{} "successfully got  pending orders return request"
// @Failure 500 {object} res.Response{} "faild to get pending order return requests"
func (c *OrderHandler) GetAllPendingReturns(ctx *gin.Context) {

	pageNumber, err1 := utils.StringToUint(ctx.Query("page_number"))
	count, err2 := utils.StringToUint(ctx.Query("count"))

	err1 = errors.Join(err1, err2)

	if err1 != nil {
		response := res.ErrorResponse(400, "invalid inputs", err1.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	pagination := req.ReqPagination{
		PageNumber: pageNumber,
		Count:      count,
	}

	orderReturns, err := c.orderUseCase.GetAllPendingOrderReturns(ctx, pagination)
	if err != nil {
		response := res.ErrorResponse(500, "faild to get pending order return requests", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if orderReturns == nil {
		response := res.SuccessResponse(200, "there is no pendinng orders return request", nil)
		ctx.JSON(200, response)
		return
	}

	response := res.SuccessResponse(200, "successfully got  pending orders return request", orderReturns)
	ctx.JSON(http.StatusOK, response)
}

// UpdategReturnRequest godoc
// @summary api for admin to supdate the order_return request from user
// @description admin can approve, cancell etc. updation on User Orders_return
// @id UpdategReturnRequest
// @tags Admin Orders
// @Param input body req.ReqUpdatReturnOrder{} true "Input Fiields"
// @Router /admin/orders/returns/pending [put]
// @Success 200 {object} res.Response{} "successfully order_response updated"
// @Failure 500 {object} res.Response{} "invalid input"
func (c *OrderHandler) UpdateReturnRequest(ctx *gin.Context) {

	var body req.ReqUpdatReturnOrder
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response := res.ErrorResponse(400, "invalid input", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	err := c.orderUseCase.UpdateReturnRequest(ctx, body)
	if err != nil {
		response := res.ErrorResponse(400, "faild to update order_return", err.Error(), body)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := res.SuccessResponse(200, "successfully order_response updated")
	ctx.JSON(http.StatusOK, response)
}
