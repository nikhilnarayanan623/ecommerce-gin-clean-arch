package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
)

// RazorpayCheckout godoc
// @summary api for create razorpay payment order for shop order
// @security ApiKeyAuth
// @tags User Cart
// @id RazorpayCheckout
// @Param payment_method_id formData uint true "Payment Method ID"
// @Param shop_order_id formData uint true "ShopOrder ID"
// @Router /carts/place-order/razorpay-checkout [post]
// @Success 200 {object} res.Response{} "successfully razorpay payment order created"
// @Failure 400 {object} res.Response{}  "faild to create razorpay payment order"
func (c *OrderHandler) RazorpayCheckout(ctx *gin.Context) {

	UserID := utils.GetUserIdFromContext(ctx)
	paymentMethodID, err1 := utils.StringToUint(ctx.Request.PostFormValue("payment_method_id"))
	shopOrderID, err2 := utils.StringToUint(ctx.Request.PostFormValue("shop_order_id"))

	err := errors.Join(err1, err2)
	if err != nil {
		fmt.Println(err)
		response := res.ErrorResponse(400, "invalid inputs", err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// make razorpay order and make razorpay order respones
	razorpayOrder, err := c.orderUseCase.GetRazorpayOrder(ctx, UserID, shopOrderID, paymentMethodID)
	if err != nil {
		response := res.ErrorResponse(500, "faild to create razorpay order ", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	// make a respone of order and and razorpay for fron-end validation
	response := gin.H{
		"Razorpay": true,
		"Order":    razorpayOrder,
	}

	ctx.JSON(200, response)
}

// razorpay verification

// RazorpayVerify godoc
// @summary api user for verify razorpay payment
// @security ApiKeyAuth
// @tags User Cart
// @id RazorpayVerify
// @Param razorpay_order_id formData string true "razorpay_order_id"
// @Param razorpay_payment_id formData string true "razorpay_payment_id"
// @Param razorpay_signature formData string false "razorpay_signature"
// @Param shop_order_id formData string true "shop_order_id"
// @Param payment_method_id formData uint true "payment_method_id"
// @Router /carts/place-order/razorpay-verify [post]
// @Success 200 {object} res.Response{} "faild to veify payment"
// @Failure 400 {object} res.Response{}  "successfully payment completed and order approved"
func (c *OrderHandler) RazorpayVerify(ctx *gin.Context) {

	userID := utils.GetUserIdFromContext(ctx)

	// take value as form value from ajax call
	razorpayPaymentID := ctx.Request.PostFormValue("razorpay_payment_id")
	razorpayOrderID := ctx.Request.PostFormValue("razorpay_order_id")
	razorpaySignature := ctx.Request.PostFormValue("razorpay_signature")
	shopOrderID, err1 := utils.StringToUint(ctx.Request.PostFormValue("shop_order_id"))
	paymentMethodID, err2 := utils.StringToUint(ctx.Request.PostFormValue("payment_method_id"))
	err := errors.Join(err1, err2)
	if err != nil {
		response := res.ErrorResponse(400, "can't make order", "shop_order id is or coupon id is not int valid type", nil)
		ctx.JSON(400, response)
		return
	}

	//verify the razorpay payment
	err = utils.VeifyRazorpayPayment(razorpayOrderID, razorpayPaymentID, razorpaySignature)
	if err != nil {
		response := res.ErrorResponse(400, "faild to verfiy razorpay payment", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	err = c.orderUseCase.ApproveShopOrderAndClearCart(ctx, userID, shopOrderID, paymentMethodID)
	if err != nil {
		response := res.ErrorResponse(400, "faild to place order faild on approve and clear cart", err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := res.SuccessResponse(200, "successfully payment completed and order approved", nil)
	ctx.JSON(200, response)
}
