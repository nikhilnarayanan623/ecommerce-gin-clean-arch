package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/config"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/res"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"github.com/razorpay/razorpay-go"
)

type OrderHandler struct {
	orderUseCase interfaces.OrderUseCase
}

func NewOrderHandler(orderUseCase interfaces.OrderUseCase) *OrderHandler {
	return &OrderHandler{orderUseCase: orderUseCase}
}

func (c *OrderHandler) CartOrderPayementSelectPage(ctx *gin.Context) {

	Payments, err := c.orderUseCase.GetAllPaymentMethods(ctx)
	if err != nil {
		ctx.HTML(200, "paymentForm.html", nil)
		return
	}

	ctx.HTML(200, "paymentForm.html", Payments)
}

// PlaceOrderForCart godoc
// @summary api cart order
// @security ApiKeyAuth
// @tags Order
// @id PlaceOrderForCart
// @Param        inputs   body     req.ReqCheckout{}   true  "Input Field"
// @Router /carts/place-order/cod [post]
// @Success 200 {object} res.Response{} "place order"
// @Failure 400 {object} res.Response{}  "faill place order"
func (c *OrderHandler) PlaceOrderForCartCOD(ctx *gin.Context) {

	var body req.ReqCheckout
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response := res.ErrorResponse(400, "invalid inputs", err.Error(), body)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	body.UserID = helper.GetUserIdFromContext(ctx)

	// checkout the order
	resCheckout, err := c.orderUseCase.OrderCheckOut(ctx, body)
	if err != nil {
		response := res.ErrorResponse(400, "faild to place order on COD", err.Error(), body)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// check order status is COD or not
	if resCheckout.PaymentType != "COD" {
		respones := res.ErrorResponse(400, "can't place order order", "payement type is not COD", nil)
		ctx.AbortWithStatusJSON(400, respones)
		return
	}
	// place order on COD
	err = c.orderUseCase.PlaceOrderCOD(ctx, resCheckout)
	if err != nil {
		response := res.ErrorResponse(400, "faild to place order on COD", err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := res.SuccessResponse(200, "successfully placed order for COD")
	ctx.JSON(http.StatusOK, response)

}

// RazorpayPage godoc
// @summary api for create an razorpay order
// @security ApiKeyAuth
// @tags Order
// @id RazorpayPage
// @Param payment_method_id formData uint true "Payment Method ID"
// @Param address_id formData uint true "Address ID"
// @Param coupon_code formData string false "Coupon Code"
// @Router /carts/place-order/razorpay-checkout [post]
// @Success 200 {object} res.Response{} "place order"
// @Failure 400 {object} res.Response{}  "faill place order"
func (c *OrderHandler) RazorpayCheckout(ctx *gin.Context) {

	paymentMethodID, err1 := helper.StringToUint(ctx.Request.PostFormValue("payment_method_id"))
	addressID, err2 := helper.StringToUint(ctx.Request.PostFormValue("address_id"))
	couponCode := ctx.Request.PostFormValue("coupon_code")

	err := errors.Join(err1, err2)
	if err != nil {
		fmt.Println(err)
		response := res.ErrorResponse(400, "invalid inputs", err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	UserID := helper.GetUserIdFromContext(ctx)

	var body = req.ReqCheckout{
		UserID:          UserID,
		PaymentMethodID: paymentMethodID,
		CouponCode:      couponCode,
		AddressID:       addressID,
	}

	// checkout the order
	resCheckout, err := c.orderUseCase.OrderCheckOut(ctx, body)
	if err != nil {
		response := res.ErrorResponse(400, "faild to place order on razor pay", err.Error(), body)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// frist get razor pay key and secret
	razorPaykey := config.GetCofig().RazorPayKey
	razorPaysecret := config.GetCofig().RazorPaySecret

	// create a new client
	client := razorpay.NewClient(razorPaykey, razorPaysecret)

	razorPayAmount := resCheckout.AmountToPay * 100

	data := map[string]interface{}{
		"amount":   razorPayAmount,
		"currency": "INR",
		"receipt":  "some_receipt_id",
	}

	razorPayRes, err := client.Order.Create(data, nil)
	if err != nil {
		response := res.ErrorResponse(500, "faild to create razor pay order", err.Error(), nil)
		ctx.AbortWithStatusJSON(500, response)
		return
	}
	fmt.Println("razor pay response when orde creating", razorPayRes)

	Order := gin.H{
		"Key":           razorPaykey,
		"UserID":        resCheckout.UserID,
		"AmountToPay":   resCheckout.AmountToPay,
		"RazorpayAmout": razorPayAmount,
		"OrderID":       razorPayRes["id"],
		"Email":         "nikhil@gmail.com",
		"Phone":         "62385893260",
	}

	// make a respone of order and and razorpay for fron-end validation
	response := gin.H{
		"Razorpay": true,
		"Order":    Order,
	}

	ctx.JSON(200, response)
}

// razorpay verification
func (c *OrderHandler) RazorpayVerify(ctx *gin.Context) {
	// struct of razorpay varification
	var data req.ReqRazorpayVeification
	// take value as form value
	data.RazorpayPaymentID = ctx.Request.PostFormValue("razorpay_payment_id")
	data.RazorpayOrderID = ctx.Request.PostFormValue("razorpay_order_id")
	data.RazorpaySignature = ctx.Request.PostFormValue("razorpay_signature")

	//verify the signature
	err := helper.VeifyRazorPaySignature(data.RazorpayOrderID, data.RazorpayPaymentID, data.RazorpaySignature)
	if err != nil {
		respones := res.ErrorResponse(400, "faild to veify payment", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, respones)
		return
	}
	// get razorpay key and secret
	razorpayKey, razorPaySecret := config.GetCofig().RazorPayKey, config.GetCofig().RazorPaySecret
	// create a new razorpay client
	razorpayClient := razorpay.NewClient(razorpayKey, razorPaySecret)
	payment, err := razorpayClient.Payment.Fetch(data.RazorpayPaymentID, nil, nil)
	if err != nil {
		response := res.ErrorResponse(400, "faild to get payment details", err.Error(), nil)
		ctx.JSON(400, response)
		return
	}

	if payment["status"] != "captured" {
		response := res.ErrorResponse(400, "payment faild", "payment not got on razorpay", payment)
		ctx.JSON(400, response)
		return
	}

	response := res.SuccessResponse(200, "successfully payment completed", payment)
	ctx.JSON(200, response)
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
// @summary api for admin to show all order
// @description admin can see all orders in application
// @id GetAllShopOrders
// @tags Orders
// @Router /admin/orders [get]
// @Success 200 {object} res.Response{} "successfully order list got"
// @Failure 500 {object} res.Response{} "faild to get shop order data"
func (c *OrderHandler) GetAllShopOrders(ctx *gin.Context) {

	resShopOrdersPage, err := c.orderUseCase.GetAllShopOrders(ctx)
	if err != nil {
		response := res.ErrorResponse(500, "faild to get shop order data", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if resShopOrdersPage.Orders == nil {
		response := res.SuccessResponse(200, "order list is empty", nil)
		ctx.JSON(200, response)
		return
	}

	response := res.SuccessResponse(200, "successfully order list got", resShopOrdersPage)
	ctx.JSON(http.StatusOK, response)
}

// SubmitReturnRequest godoc
// @summary api for user to request a return for an order
// @description user can request return for placed orders
// @id SubmitReturnRequest
// @tags Orders
// @Router /orders/return [put]
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
// @tags Orders
// @Router /orders/returns [get]
// @Success 200 {object} res.Response{} "successfully got all order returns"
// @Failure 500 {object} res.Response{} "faild to get order returns"
func (c *OrderHandler) GetAllOrderReturns(ctx *gin.Context) {

	orderReturns, err := c.orderUseCase.GetAllOrderReturns(ctx)
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
// @tags Orders
// @Router /admin/orders/returns/pending [get]
// @Success 200 {object} res.Response{} "successfully got  pending orders return request"
// @Failure 500 {object} res.Response{} "faild to get pending order return requests"
func (c *OrderHandler) GetAllPendingReturns(ctx *gin.Context) {

	orderReturns, err := c.orderUseCase.GetAllPendingOrderReturns(ctx)
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
// @description admin can approve, cancell etc. updation on user order_return
// @id UpdategReturnRequest
// @tags Orders
// @Router /admin/orders/returns/penging [put]
// @Success 200 {object} res.Response{} "successfully order_response updated"
// @Failure 500 {object} res.Response{} "invalid input"
func (c *OrderHandler) UpdateReturnRequest(ctx *gin.Context) {

	var body req.ReqUpdatReturnReq
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
