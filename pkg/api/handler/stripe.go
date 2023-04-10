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

func (c *OrderHandler) StripPaymentCheckout(ctx *gin.Context) {

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
		response := res.ErrorResponse(400, "faild to place order on stripe", err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}
	// check payment type is  strip or not
	if paymentMethod.PaymentType != "Stripe" {
		respones := res.ErrorResponse(400, "can't place order order", "selected payment_method_id is not for Stripe ", nil)
		ctx.AbortWithStatusJSON(400, respones)
		return
	}

	body := req.ReqPlaceOrder{
		PaymentMethodID: paymentMethodID,
		AddressID:       addressID,
	}

	// get the order
	userOrder, err := c.orderUseCase.GetOrderDetails(ctx, UserID, body)
	if err != nil {
		response := res.ErrorResponse(400, "faild to place order on stripe", err.Error(), body)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// make a stripe order order
	stripeOrder, err := c.orderUseCase.GetStripeOrder(ctx, UserID, userOrder)

	if err != nil {
		response := res.ErrorResponse(500, "faild to create stripe order ", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	// create a shop order
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
		response := res.ErrorResponse(500, "faild to save order for user on place stripe", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	// set the shopOrderId on the stripOrder
	stripeOrder.ShopOrderID = shopOrderID

	ctx.JSON(200, stripeOrder)
}

func (c *OrderHandler) Success(ctx *gin.Context) {

	fmt.Println("\n\n\nsuccess")
}
func (c *OrderHandler) Cancell(ctx *gin.Context) {
	fmt.Println("\n\n\ncancell")
}
