package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	handlerInterface "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/response"
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
// @Summary Render Payment Page (User)
// @Description API for user to render payment select page
// @Security ApiKeyAuth
// @Id CartOrderPaymentSelectPage
// @Tags User Payment
// @Router /carts/checkout/payment-select-page [get]
// @Success 200 {object} response.Response{} "Successfully rendered payment page"
// @Failure 500 {object} response.Response{}   "Failed to render payment page"
func (c *paymentHandler) CartOrderPaymentSelectPage(ctx *gin.Context) {

	Payments, err := c.paymentUseCase.FindAllPaymentMethods(ctx)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to render payment page", err, nil)
		return
	}

	ctx.HTML(200, "paymentForm.html", Payments)
}

// // AddPaymentMethod godoc
// // @summary api for admin to add a new payment method
// // @security ApiKeyAuth
// // @tags Admin Payment
// // @id AddPaymentMethod
// // @Router /admin/payment-method [post]
// // @Success 200 {object} response.Response{} "successfully payment added"
// // @Success 400 {object} response.Response{} "Failed to bind input"
// // @Failure 500 {object} response.Response{}  "Failed to add payment method"
// func (c *paymentHandler) AddPaymentMethod(ctx *gin.Context) {

// 	var body request.PaymentMethod

// 	if err := ctx.ShouldBindJSON(&body); err != nil {
// 		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
// 		return
// 	}

// 	var paymentMethod domain.PaymentMethod
// 	copier.Copy(&paymentMethod, &body)

// 	err := c.paymentUseCase.SavePaymentMethod(ctx, paymentMethod)

// 	if err != nil {
// 		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to add payment_method", err, nil)
// 		return
// 	}

// 	response.SuccessResponse(ctx, http.StatusOK, "Successfully added payment method", nil)
// }

// UpdatePaymentMethod godoc
// @Summary Update payment method (Admin)
// @Description API for admin to change maximum price or block or unblock the payment method
// @security ApiKeyAuth
// @tags Admin Payment
// @id UpdatePaymentMethod
// @Router /admin/payment-method  [put]
// @Success 200 {object} response.Response{} "Successfully payment method updated"
// @Success 400 {object} response.Response{} "Invalid inputs"
// @Failure 500 {object} response.Response{}  "Failed to update payment method"
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

// FindAllPaymentMethodsAdmin godoc
// @summary Get payment methods (Admin)
// @Description API for admin to get all payment methods
// @security ApiKeyAuth
// @tags Admin Payment
// @id FindAllPaymentMethodsAdmin
// @Router /admin/payment-methods [get]
// @Success 200 {object} response.Response{} "Failed to retrieve payment methods"
// @Failure 500 {object} response.Response{}   "Successfully retrieved all payment methods"
func (c *paymentHandler) FindAllPaymentMethodsAdmin() func(ctx *gin.Context) {
	return c.findAllPaymentMethods()
}

// FindAllPaymentMethodsUser godoc
// @summary Get payment methods (User)
// @Description API for user to get all payment methods
// @security ApiKeyAuth
// @tags User Payment
// @id FindAllPaymentMethodsUser
// @Router /payment-methods [get]
// @Success 200 {object} response.Response{} "Failed to retrieve payment methods"
// @Failure 500 {object} response.Response{}   "Successfully retrieved all payment methods"
func (c *paymentHandler) FindAllPaymentMethodsUser() func(ctx *gin.Context) {
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
