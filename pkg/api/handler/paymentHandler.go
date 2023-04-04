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
	"github.com/razorpay/razorpay-go"
)

func (c *OrderHandler) GetAllPaymentMethods(ctx *gin.Context) {

	paymentMethods, err := c.orderUseCase.GetAllPaymentMethods(ctx)
	if err != nil {
		response := res.ErrorResponse(500, "faild to get all payment methods", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := res.SuccessResponse(200, "successfully got all payment methods", paymentMethods)
	ctx.JSON(http.StatusOK, response)
}

// RazorpayPage godoc
// @summary api for create an razorpay order
// @security ApiKeyAuth
// @tags User Order
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

	paymentMethod, err := c.orderUseCase.GetPaymentMethodByID(ctx, body.PaymentMethodID)
	if err != nil {
		response := res.ErrorResponse(400, "faild to place order on razor pay", err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	// check payment type is  razorpay or not
	if paymentMethod.PaymentType != "Razorpay" {
		respones := res.ErrorResponse(400, "can't place order order", "selected payment_method_id is not for RazorPay ", nil)
		ctx.AbortWithStatusJSON(400, respones)
		return
	}

	// checkout the order
	resCheckout, err := c.orderUseCase.OrderCheckOut(ctx, body)
	if err != nil {
		response := res.ErrorResponse(400, "faild to place order on razor pay", err.Error(), body)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// check payment type is  razorpay or not
	if resCheckout.PaymentType != "Razorpay" {
		respones := res.ErrorResponse(400, "can't place order order", "selected payment_method_id is not for RazorPay", resCheckout.PaymentType)
		ctx.AbortWithStatusJSON(400, respones)
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

	// save order as pending
	shopOrderID, err := c.orderUseCase.SaveOrder(ctx, resCheckout)
	if err != nil {
		response := res.ErrorResponse(400, "faild to place order on COD for save order", err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	// create razorpay order
	Order := gin.H{
		"Key":           razorPaykey,
		"UserID":        resCheckout.UserID,
		"AmountToPay":   resCheckout.AmountToPay,
		"RazorpayAmout": razorPayAmount,
		"OrderID":       razorPayRes["id"],
		"Email":         "nikhil@gmail.com",
		"Phone":         "62385893260",
		"ShopOrderID":   shopOrderID,
		"CouponCode":    resCheckout.CouponCode,
	}

	// make a respone of order and and razorpay for fron-end validation
	response := gin.H{
		"Razorpay": true,
		"Order":    Order,
	}

	ctx.JSON(200, response)
}

// razorpay verification

// RazorpayVerify godoc
// @summary api user for verify razorpay payment
// @security ApiKeyAuth
// @tags User Order
// @id RazorpayVerify
// @Param payment_method_id formData uint true "Payment Method ID"
// @Param address_id formData uint true "Address ID"
// @Param coupon_code formData string false "Coupon Code"
// @Router /carts/place-order/razorpay-verify [post]
// @Success 200 {object} res.Response{} "faild to veify payment"
// @Failure 400 {object} res.Response{}  "successfully payment completed and order approved"
func (c *OrderHandler) RazorpayVerify(ctx *gin.Context) {

	// struct of razorpay varification
	var (
		body req.ReqRazorpayVeification
		err  error
	)

	body.UserID = helper.GetUserIdFromContext(ctx)

	// take value as form value from ajax call
	body.RazorpayPaymentID = ctx.Request.PostFormValue("razorpay_payment_id")
	body.RazorpayOrderID = ctx.Request.PostFormValue("razorpay_order_id")
	body.RazorpaySignature = ctx.Request.PostFormValue("razorpay_signature")
	body.ShopOrderID, err = helper.StringToUint(ctx.Request.PostFormValue("shop_order_id"))
	couponCode := ctx.Request.PostFormValue("coupon_code")
	fmt.Println("coupon code ", couponCode)
	if err != nil {
		response := res.ErrorResponse(400, "can't make order", "shop_order id is not valid", nil)
		fmt.Println(response)
		ctx.JSON(400, response)
		return
	}

	//verify the razorpay signature
	err = helper.VeifyRazorPaySignature(body.RazorpayOrderID, body.RazorpayPaymentID, body.RazorpaySignature)
	if err != nil {
		response := res.ErrorResponse(400, "faild to veify payment", err.Error(), nil)
		fmt.Println(response)

		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	// get razorpay key and secret
	razorpayKey, razorPaySecret := config.GetCofig().RazorPayKey, config.GetCofig().RazorPaySecret
	// create a new razorpay client
	razorpayClient := razorpay.NewClient(razorpayKey, razorPaySecret)
	payment, err := razorpayClient.Payment.Fetch(body.RazorpayPaymentID, nil, nil)
	if err != nil {
		response := res.ErrorResponse(400, "faild to get payment details", err.Error(), nil)
		fmt.Println(response)

		ctx.JSON(400, response)
		return
	}

	if payment["status"] != "captured" {
		response := res.ErrorResponse(400, "payment faild", "payment not got on razorpay", payment)
		fmt.Println(response)

		ctx.JSON(400, response)
		return
	}

	// approve the order using order id
	err = c.orderUseCase.ApproveOrder(ctx, body.UserID, body.ShopOrderID, couponCode)
	if err != nil {
		response := res.ErrorResponse(400, "faild to place order on Razorpay for approve", err.Error(), nil)
		fmt.Println(response)

		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	response := res.SuccessResponse(200, "successfully payment completed and order approved", payment)

	ctx.Set("razorpay-response", response)
	// set the response on context for coupon code handler
	// ctx.JSON(200, response)
}
