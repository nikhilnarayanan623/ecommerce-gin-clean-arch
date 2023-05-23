package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	handlerInterface "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
)

type paymentHandler struct {
	paymentUseCase interfaces.PaymentUseCase
}

func NewPaymentHandler(paymentUseCase interfaces.PaymentUseCase) handlerInterface.PaymentHandler {
	return &paymentHandler{
		paymentUseCase: paymentUseCase,
	}
}

// CartOrderPayementSelectPage godoc
// @summary api for render the html page of payment select
// @security ApiKeyAuth
// @tags User Payment
// @id CartOrderPayementSelectPage
// @Router /carts/checkout/payemt-select-page [get]
// @Success 200 {object} res.Response{} "successfully order placed"
// @Failure 500 {object} res.Response{}   "faild to render payment page"
func (c *paymentHandler) CartOrderPayementSelectPage(ctx *gin.Context) {

	Payments, err := c.paymentUseCase.GetAllPaymentMethods(ctx)
	if err != nil {
		response := res.ErrorResponse(500, "faild to render payment page", err.Error(), nil)
		ctx.JSON(500, response)
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
// @Success 200 {object} res.Response{} "successfully payment added"
// @Success 400 {object} res.Response{} "faild to bind input"
// @Failure 500 {object} res.Response{}   "faild to add payment method"
func (c *paymentHandler) AddPaymentMethod(ctx *gin.Context) {

	var body req.PaymentMethod

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response := res.ErrorResponse(400, "invalid inputs", err.Error(), body)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	var paymentMethod domain.PaymentMethod
	copier.Copy(&paymentMethod, &body)

	err := c.paymentUseCase.AddPaymentMethod(ctx, paymentMethod)

	if err != nil {
		response := res.ErrorResponse(400, "faild to add payment_method", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := res.SuccessResponse(200, "successfully added payment_method", nil)
	ctx.JSON(http.StatusOK, response)
}

// UpdatePaymentMethod godoc
// @summary api for admin to update payment details
// @security ApiKeyAuth
// @tags Admin Payment
// @id UpdatePaymentMethod
// @Router /admin/payment-method  [put]
// @Success 200 {object} res.Response{} "successfully updated payment method"
// @Success 400 {object} res.Response{} "faild to bind input"
// @Failure 500 {object} res.Response{}   "faild to update payment details"
func (c *paymentHandler) UpdatePaymentMethod(ctx *gin.Context) {

	var body req.PaymentMethodUpdate

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response := res.ErrorResponse(400, "faild to bind input", err.Error(), body)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	err := c.paymentUseCase.EditPaymentMethod(ctx, body)

	if err != nil {
		response := res.ErrorResponse(400, "faild to update payment_method", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := res.SuccessResponse(200, "successfully updated payment_method", nil)
	ctx.JSON(http.StatusOK, response)
}

// GetAllPaymentMethods godoc
// @summary api for get all payment methods
// @security ApiKeyAuth
// @tags User Payment
// @id GetAllPaymentMethods
// @Router /admin/payment-method [get]
// @Success 200 {object} res.Response{} "successfully get payment method"
// @Failure 500 {object} res.Response{}   "faild to get all payment methods"
func (c *paymentHandler) GetAllPaymentMethods(ctx *gin.Context) {

	paymentMethods, err := c.paymentUseCase.GetAllPaymentMethods(ctx)
	if err != nil {
		response := res.ErrorResponse(500, "faild to get all payment methods", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if paymentMethods == nil {
		response := res.SuccessResponse(200, "there is no payment_methods available to show")
		ctx.JSON(http.StatusOK, response)
		return
	}

	response := res.SuccessResponse(200, "successfully got all payment methods", paymentMethods)
	ctx.JSON(http.StatusOK, response)
}
