package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	handlerInterface "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
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
// @summary api for render the html page of payment select
// @security ApiKeyAuth
// @tags User Payment
// @id CartOrderPaymentSelectPage
// @Router /carts/checkout/payment-select-page [get]
// @Success 200 {object} response.Response{} "successfully order placed"
// @Failure 500 {object} response.Response{}   "Failed to render payment page"
func (c *paymentHandler) CartOrderPaymentSelectPage(ctx *gin.Context) {

	Payments, err := c.paymentUseCase.FindAllPaymentMethods(ctx)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to render payment page", err, nil)
		return
	}

	ctx.HTML(200, "paymentForm.html", Payments)
}

// AddPaymentMethod godoc
// @summary api for admin to add a new payment method
// @security ApiKeyAuth
// @tags Admin Payment
// @id AddPaymentMethod
// @Router /admin/payment-method [post]
// @Success 200 {object} response.Response{} "successfully payment added"
// @Success 400 {object} response.Response{} "Failed to bind input"
// @Failure 500 {object} response.Response{}  "Failed to add payment method"
func (c *paymentHandler) AddPaymentMethod(ctx *gin.Context) {

	var body request.PaymentMethod

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	var paymentMethod domain.PaymentMethod
	copier.Copy(&paymentMethod, &body)

	err := c.paymentUseCase.SavePaymentMethod(ctx, paymentMethod)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to add payment_method", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully added payment method", nil)
}

// UpdatePaymentMethod godoc
// @summary api for admin to update payment details
// @security ApiKeyAuth
// @tags Admin Payment
// @id UpdatePaymentMethod
// @Router /admin/payment-method  [put]
// @Success 200 {object} response.Response{} "Successfully payment method updated"
// @Success 400 {object} response.Response{} "Failed to bind input"
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

// FindAllPaymentMethods godoc
// @summary api for get all payment methods
// @security ApiKeyAuth
// @tags User Payment
// @id FindAllPaymentMethods
// @Router /admin/payment-method [get]
// @Success 200 {object} response.Response{} "Failed to find payment methods"
// @Failure 500 {object} response.Response{}   "Successfully found all payment methods"
func (c *paymentHandler) FindAllPaymentMethods(ctx *gin.Context) {

	paymentMethods, err := c.paymentUseCase.FindAllPaymentMethods(ctx)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find payment methods", err, nil)
		return
	}

	if paymentMethods == nil {
		response.SuccessResponse(ctx, http.StatusOK, "No payment methods found")
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully found all payment methods", paymentMethods)
}
