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

// StripPaymentCheckout godoc
//	@Summary		Stripe checkout (User)
//	@Security		BearerAuth
//	@Description	API for user to create stripe payment
//	@Tags			User Payment
//	@Id				StripPaymentCheckout
//	@Param			shop_order_id	formData	string	true	"Shop Order ID"
//	@Router			/carts/place-order/stripe-checkout [post]
//	@Success		200	{object}	response.Response{}	"successfully stripe payment order created"
//	@Failure		500	{object}	response.Response{}	"Failed to create stripe order"
func (c *paymentHandler) StripPaymentCheckout(ctx *gin.Context) {

	shopOrderID, err := request.GetFormValuesAsUint(ctx, "shop_order_id")
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindFormValueMessage, err, nil)
		return
	}

	UserID := utils.GetUserIdFromContext(ctx)

	stripeOrder, err := c.paymentUseCase.MakeStripeOrder(ctx, UserID, shopOrderID)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to create stripe order", err, nil)
		return
	}

	stripeResponse := response.OrderPayment{
		PaymentOrder: stripeOrder,
		PaymentType:  domain.StripePayment,
	}

	ctx.JSON(http.StatusOK, stripeResponse)
}

// StripePaymentVeify godoc
//	@Summary		Stripe verify (User)
//	@Security		BearerAuth
//	@Description	API for user to callback backend after stripe payment for verification
//	@Tags			User Payment
//	@Id				StripePaymentVeify
//	@Param			stripe_payment_id	formData	string	true	"Stripe payment ID"
//	@Param			shop_order_id		formData	string	true	"Shop Order ID"
//	@Router			/carts/place-order/stripe-verify [post]
//	@Success		200	{object}	response.Response{}	"Successfully stripe payment verified"
//	@Failure		402	{object}	response.Response{}	"Payment not approved"
//	@Failure		500	{object}	response.Response{}	"Failed to Approve order"
func (c *paymentHandler) StripePaymentVeify(ctx *gin.Context) {

	shopOrderID, err1 := request.GetFormValuesAsUint(ctx, "shop_order_id")
	stripePaymentID, err2 := request.GetFormValuesAsString(ctx, "stripe_payment_id")

	err := errors.Join(err1, err2)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindFormValueMessage, err, nil)
		return
	}

	userID := utils.GetUserIdFromContext(ctx)

	err = c.paymentUseCase.VerifyStripOrder(ctx, stripePaymentID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, usecase.ErrPaymentNotApproved) {
			statusCode = http.StatusPaymentRequired
		}
		response.ErrorResponse(ctx, statusCode, "Failed to verify stripe payment", err, nil)
		return
	}

	approveReq := request.ApproveOrder{
		ShopOrderID: shopOrderID,
		PaymentType: domain.StripePayment,
	}

	err = c.paymentUseCase.ApproveShopOrderAndClearCart(ctx, userID, approveReq)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to Approve order", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully stripe payment verified", nil)
}
