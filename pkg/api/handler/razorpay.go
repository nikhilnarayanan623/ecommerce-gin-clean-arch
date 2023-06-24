package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/response"
)

// RazorpayCheckout godoc
// @summary api for create razorpay payment order
// @security ApiKeyAuth
// @tags User Cart
// @id RazorpayCheckout
// @Param address_id formData string true "Address ID"
// @Router /carts/place-order/razorpay-checkout [post]
// @Success 200 {object} response.Response{} "successfully razorpay payment order created"
// @Failure 400 {object} response.Response{}  "faild to create razorpay payment order"
func (c *OrderHandler) RazorpayCheckout(ctx *gin.Context) {

	addressID, err := request.GetFormValuesAsUint(ctx, "address_id")
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindFormValueMessage, err, nil)
		return
	}

	UserID := utils.GetUserIdFromContext(ctx)

	body := request.PlaceOrder{
		AddressID:   addressID,
		PaymentType: domain.RazopayPayment,
	}

	shopOrderID, err := c.orderUseCase.SaveOrder(ctx, UserID, body)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to save order", err, nil)
		return
	}

	razorpayOrder, err := c.orderUseCase.MakeRazorpayOrder(ctx, UserID, shopOrderID)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to make razorpay order ", err, nil)
		return
	}

	razorPayRes := gin.H{
		"razorpay": true,
		"order":    razorpayOrder,
	}
	ctx.JSON(http.StatusOK, razorPayRes)
	// response.SuccessResponse(ctx, http.StatusCreated, "Successfully razor pay order created", razorPayRes)
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
// @Success 200 {object} response.Response{} "faild to veify payment"
// @Failure 400 {object} response.Response{}  "successfully payment completed and order approved"
func (c *OrderHandler) RazorpayVerify(ctx *gin.Context) {

	userID := utils.GetUserIdFromContext(ctx)

	razorpayPaymentID, err1 := request.GetFormValuesAsString(ctx, "razorpay_payment_id")
	razorpayOrderID, err2 := request.GetFormValuesAsString(ctx, "razorpay_order_id")
	razorpaySignature, err3 := request.GetFormValuesAsString(ctx, "razorpay_order_id")

	shopOrderID, err4 := request.GetFormValuesAsUint(ctx, "shop_order_id")

	err := errors.Join(err1, err2, err3, err4)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindFormValueMessage, err, nil)
		return
	}

	err = utils.VerifyRazorpayPayment(razorpayOrderID, razorpayPaymentID, razorpaySignature)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, "Failed to verify razorpay payment", err, nil)
		return
	}

	err = c.orderUseCase.ApproveShopOrderAndClearCart(ctx, userID, shopOrderID)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to Approve order", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully razorpay payment verified", nil)
}
