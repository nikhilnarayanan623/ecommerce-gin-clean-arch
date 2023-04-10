package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
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
// @tags User Cart
// @id RazorpayPage
// @Param payment_method_id formData uint true "Payment Method ID"
// @Param address_id formData uint true "Address ID"
// @Router /carts/place-order/razorpay-checkout [post]
// @Success 200 {object} res.Response{} "place order"
// @Failure 400 {object} res.Response{}  "faill place order"
func (c *OrderHandler) RazorpayCheckout(ctx *gin.Context) {

	UserID := utils.GetUserIdFromContext(ctx)
	paymentMethodID, err1 := utils.StringToUint(ctx.Request.PostFormValue("payment_method_id"))
	addressID, err2 := utils.StringToUint(ctx.Request.PostFormValue("address_id"))

	err := errors.Join(err1, err2)
	if err != nil {
		fmt.Println(err)
		response := res.ErrorResponse(400, "invalid inputs", err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	paymentMethod, err := c.orderUseCase.GetPaymentMethodByID(ctx, paymentMethodID)
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

	body := req.ReqPlaceOrder{
		PaymentMethodID: paymentMethodID,
		AddressID:       addressID,
	}

	// checkout the order
	userOrder, err := c.orderUseCase.GetOrderDetails(ctx, UserID, body)
	if err != nil {
		response := res.ErrorResponse(400, "faild to place order on razor pay", err.Error(), body)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// make razorpay order and make razorpay order respones
	razorpayOrder, err := c.orderUseCase.GetRazorpayOrder(ctx, UserID, userOrder)
	if err != nil {
		response := res.ErrorResponse(500, "faild to create razorpay order ", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	// save user order as pending only
	shopOrder := domain.ShopOrder{
		UserID:          UserID,
		AddressID:       body.AddressID,
		OrderTotalPrice: userOrder.AmountToPay,
		Discount:        userOrder.Discount,
		PaymentMethodID: body.PaymentMethodID,
		OrderDate:       time.Now(),
	}
	shopOrderID, err := c.orderUseCase.SaveOrder(ctx, shopOrder)
	if err != nil {
		response := res.ErrorResponse(500, "faild to save order for user on place razorpay", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	// set shop_order_id on razorpay order
	razorpayOrder.ShopOrderID = shopOrderID

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
// @Param payment_method_id formData uint true "Payment Method ID"
// @Param address_id formData uint true "Address ID"
// @Param coupon_code formData string false "Coupon Code"
// @Router /carts/place-order/razorpay-verify [post]
// @Success 200 {object} res.Response{} "faild to veify payment"
// @Failure 400 {object} res.Response{}  "successfully payment completed and order approved"
func (c *OrderHandler) RazorpayVerify(ctx *gin.Context) {

	// // struct of razorpay varification
	// var (
	// 	body req.ReqRazorpayVeification
	// 	err  error
	// )

	userID := utils.GetUserIdFromContext(ctx)

	// take value as form value from ajax call
	razorpayPaymentID := ctx.Request.PostFormValue("razorpay_payment_id")
	razorpayOrderID := ctx.Request.PostFormValue("razorpay_order_id")
	razorpaySignature := ctx.Request.PostFormValue("razorpay_signature")
	shopOrderID, err1 := utils.StringToUint(ctx.Request.PostFormValue("shop_order_id"))
	couponID, err2 := utils.StringToUint(ctx.Request.PostFormValue("coupon_id"))

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

	// approve the order using order id
	fmt.Println(userID, shopOrderID, couponID)
	err = c.orderUseCase.ApproveOrderAndClearCart(ctx, userID, shopOrderID, couponID)
	if err != nil {
		response := res.ErrorResponse(400, "faild to place order faild on approve and clear cart", err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	response := res.SuccessResponse(200, "successfully payment completed and order approved", nil)
	ctx.JSON(200, response)
}
