package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	handlerInterface "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/response"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
)

type paymentHandler struct {
	paymentUseCase interfaces.PaymentUseCase
}

func NewPaymentHandler(paymentUseCase interfaces.PaymentUseCase) handlerInterface.PaymentHandler {
	return &paymentHandler{
		paymentUseCase: paymentUseCase,
	}
}

// CartOrderPaymentSelectPage godoc
//	@Summary		Render Payment Page (User)
//	@Security		BearerAuth
//	@Description	API for user to render payment select page
//	@Id				CartOrderPaymentSelectPage
//	@Tags			User Payment
//	@Router			/carts/checkout/payment-select-page [get]
//	@Success		200	{object}	response.Response{}	"Successfully rendered payment page"
//	@Failure		500	{object}	response.Response{}	"Failed to render payment page"
func (c *paymentHandler) CartOrderPaymentSelectPage(ctx *gin.Context) {

	Payments, err := c.paymentUseCase.FindAllPaymentMethods(ctx)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to render payment page", err, nil)
		return
	}

	ctx.HTML(200, "paymentForm.html", Payments)
}

// UpdatePaymentMethod godoc
//	@Summary		Update payment method (Admin)
//	@Security		BearerAuth
//	@Description	API for admin to change maximum price or block or unblock the payment method
//	@tags			Admin Payment
//	@id				UpdatePaymentMethod
//	@Router			/admin/payment-method  [put]
//	@Success		200	{object}	response.Response{}	"Successfully payment method updated"
//	@Success		400	{object}	response.Response{}	"Invalid inputs"
//	@Failure		500	{object}	response.Response{}	"Failed to update payment method"
func (c *paymentHandler) UpdatePaymentMethod(ctx *gin.Context) {

	var body request.PaymentMethodUpdate

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	err := c.paymentUseCase.UpdatePaymentMethod(ctx, body)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update payment method", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully payment method updated", nil)
}

// GetAllPaymentMethodsAdmin godoc
//	@summary		Get payment methods (Admin)
//	@Security		BearerAuth
//	@Description	API for admin to get all payment methods
//	@tags			Admin Payment
//	@id				GetAllPaymentMethodsAdmin
//	@Router			/admin/payment-methods [get]
//	@Success		200	{object}	response.Response{}	"Failed to retrieve payment methods"
//	@Failure		500	{object}	response.Response{}	"Successfully retrieved all payment methods"
func (c *paymentHandler) GetAllPaymentMethodsAdmin() func(ctx *gin.Context) {
	return c.findAllPaymentMethods()
}

// GetAllPaymentMethodsUser godoc
//	@summary		Get payment methods (User)
//	@Security		BearerAuth
//	@Description	API for user to get all payment methods
//	@tags			User Payment
//	@id				GetAllPaymentMethodsUser
//	@Router			/payment-methods [get]
//	@Success		200	{object}	response.Response{}	"Failed to retrieve payment methods"
//	@Failure		500	{object}	response.Response{}	"Successfully retrieved all payment methods"
func (c *paymentHandler) GetAllPaymentMethodsUser() func(ctx *gin.Context) {
	return c.findAllPaymentMethods()
}

func (c *paymentHandler) findAllPaymentMethods() func(ctx *gin.Context) {

	return func(ctx *gin.Context) {

		paymentMethods, err := c.paymentUseCase.FindAllPaymentMethods(ctx)
		if err != nil {
			response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve payment methods", err, nil)
			return
		}

		if paymentMethods == nil {
			response.SuccessResponse(ctx, http.StatusOK, "No payment methods found")
			return
		}

		response.SuccessResponse(ctx, http.StatusOK, "Successfully retrieved all payment methods", paymentMethods)
	}
}

// PaymentCOD godoc
//	@summary		Place order  for COD (User)
//	@Security		BearerAuth
//	@Description	API for user to place order for cash on delivery
//	@tags			User Payment
//	@id				PaymentCOD
//	@Param			shop_order_id	formData	string	true	"Shop Order ID"
//	@Router			/carts/place-order/cod [post]
//	@Success		200	{object}	response.Response{}	"successfully order placed for COD"
//	@Failure		500	{object}	response.Response{}	"Failed place order for COD"
func (c *paymentHandler) PaymentCOD(ctx *gin.Context) {

	shopOrderID, err := request.GetFormValuesAsUint(ctx, "shop_order_id")
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindFormValueMessage, err, nil)
		return
	}

	UserID := utils.GetUserIdFromContext(ctx)

	approveReq := request.ApproveOrder{
		ShopOrderID: shopOrderID,
		PaymentType: domain.CodPayment,
	}

	// approve the order and clear the user cart
	err = c.paymentUseCase.ApproveShopOrderAndClearCart(ctx, UserID, approveReq)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to approve order and clear cart", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully order placed for cod")
}
