package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/response"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
)

// RazorpayCheckout godoc
//	@Summary		Razorpay checkout (User)
//	@Security		BearerAuth
//	@Description	API for user to create stripe payment
//	@Security		ApiKeyAuth
//	@Tags			User Payment
//	@Id				RazorpayCheckout
//	@Param			shop_order_id	formData	string	true	"Shop Order ID"
//	@Router			/carts/place-order/razorpay-checkout [post]
//	@Success		200	{object}	response.Response{}	"successfully razorpay payment order created"
//	@Failure		500	{object}	response.Response{}	"Failed to make razorpay order"
func (c *paymentHandler) RazorpayCheckout(ctx *gin.Context) {

	shopOrderID, err := request.GetFormValuesAsUint(ctx, "shop_order_id")
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindFormValueMessage, err, nil)
		return
	}

	UserID := utils.GetUserIdFromContext(ctx)

	razorpayOrder, err := c.paymentUseCase.MakeRazorpayOrder(ctx, UserID, shopOrderID)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to make razorpay order ", err, nil)
		return
	}

	razorPayRes := response.OrderPayment{
		PaymentType:  domain.RazopayPayment,
		PaymentOrder: razorpayOrder,
	}
	ctx.JSON(http.StatusOK, razorPayRes)
	// response.SuccessResponse(ctx, http.StatusCreated, "Successfully razor pay order created", razorPayRes)
}

// RazorpayVerify godoc
//	@Summary		Razorpay verify (User)
//	@Security		BearerAuth
//	@Description	API for razorpay to callback backend for payment verification
//	@tags			User Payment
//	@id				RazorpayVerify
//	@Param			razorpay_order_id	formData	string	true	"Razorpay payment id"
//	@Param			razorpay_payment_id	formData	string	true	"Razorpay payment id"
//	@Param			razorpay_signature	formData	string	false	"Razorpay signature"
//	@Param			shop_order_id		formData	string	true	"Shop Order ID"
//	@Router			/carts/place-order/razorpay-verify [post]
//	@Success		200	{object}	response.Response{}	"Successfully razorpay payment verified"
//	@Failure		402	{object}	response.Response{}	"Payment not approved"
//	@Failure		500	{object}	response.Response{}	"Failed to Approve order"
func (c *paymentHandler) RazorpayVerify(ctx *gin.Context) {

	userID := utils.GetUserIdFromContext(ctx)

	razorpayOrderID, err2 := request.GetFormValuesAsString(ctx, "razorpay_order_id")
	razorpayPaymentID, err1 := request.GetFormValuesAsString(ctx, "razorpay_payment_id")
	razorpaySignature, err3 := request.GetFormValuesAsString(ctx, "razorpay_signature")

	shopOrderID, err4 := request.GetFormValuesAsUint(ctx, "shop_order_id")

	err := errors.Join(err1, err2, err3, err4)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindFormValueMessage, err, nil)
		return
	}

	verifyReq := request.RazorpayVerify{
		OrderID:   razorpayOrderID,
		PaymentID: razorpayPaymentID,
		Signature: razorpaySignature,
	}

	err = c.paymentUseCase.VerifyRazorPay(ctx, verifyReq)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, usecase.ErrPaymentNotApproved) {
			statusCode = http.StatusPaymentRequired
		}
		response.ErrorResponse(ctx, statusCode, "Failed to verify razorpay payment", err, nil)
		return
	}

	approveReq := request.ApproveOrder{
		ShopOrderID: shopOrderID,
		PaymentType: domain.RazopayPayment,
	}

	err = c.paymentUseCase.ApproveShopOrderAndClearCart(ctx, userID, approveReq)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to Approve order", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully razorpay payment verified", nil)
}
