package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/res"
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
