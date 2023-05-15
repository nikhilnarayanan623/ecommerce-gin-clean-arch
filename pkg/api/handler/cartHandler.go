package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/interfaces"
	usecase "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
)

type cartHandler struct {
	carUseCase usecase.CartUseCase
}

func NewCartHandler(cartUseCase usecase.CartUseCase) interfaces.CartHandler {
	return &cartHandler{
		carUseCase: cartUseCase,
	}
}

// AddToCart godoc
// @summary api for add productItem to user cart
// @description user can add a stock in product to user cart
// @security ApiKeyAuth
// @id AddToCart
// @tags User Cart
// @Param input body req.Cart true "Input Field"
// @Router /carts [post]
// @Success 200 "Successfully productItem added to cart"
// @Failure 400 "can't add the product item into cart"
func (u *cartHandler) AddToCart(ctx *gin.Context) {

	var body req.Cart
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response := res.ErrorResponse(400, "invalid input", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	// get userId and add to body
	body.UserID = utils.GetUserIdFromContext(ctx)

	err := u.carUseCase.SaveToCart(ctx, body)

	if err != nil {
		response := res.ErrorResponse(400, "faild to add product into cart", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := res.SuccessResponse(200, "successfully prodduct item added to cart", body.ProductItemID)
	ctx.JSON(http.StatusOK, response)
}

// RemoveFromCart godoc
// @summary api for remove a product from cart
// @description user can remove a signle productItem full quantity from cart
// @security ApiKeyAuth
// @id RemoveFromCart
// @tags User Cart
// @Param input body req.Cart{} true "Input Field"
// @Router /carts [delete]
// @Success 200 {object} res.Response{} "Successfully productItem removed from cart"
// @Failure 400 {object} res.Response{}  "invalid input"
// @Failure 500 {object} res.Response{}  "can't remove product item from cart"
func (u cartHandler) RemoveFromCart(ctx *gin.Context) {

	var body req.Cart
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response := res.ErrorResponse(400, "invalid input", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	body.UserID = utils.GetUserIdFromContext(ctx)

	err := u.carUseCase.RemoveCartItem(ctx, body)

	if err != nil {
		response := res.ErrorResponse(500, "can't remove product item from cart", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := res.SuccessResponse(200, "successfully product item removed form cart")
	ctx.JSON(http.StatusOK, response)
}

// UpdateCart godoc
// @summary api for updte productItem count
// @description user can inrement or drement count of a productItem in cart (min=1)
// @security ApiKeyAuth
// @id UpdateCart
// @tags User Cart
// @Param input body req.UpdateCartItem{} true "Input Field"
// @Router /carts [put]
// @Success 200 "Successfully productItem count change on cart"
// @Failure 400  "invalid input"
// @Failure 500  "can't update the count of product item on cart"
func (u *cartHandler) UpdateCart(ctx *gin.Context) {

	var body req.UpdateCartItem

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response := res.ErrorResponse(400, "invalid input", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	body.UserID = utils.GetUserIdFromContext(ctx)

	err := u.carUseCase.UpdateCartItem(ctx, body)

	if err != nil {
		response := res.ErrorResponse(500, "can't update the count of product item on cart", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := res.SuccessResponse(200, "successfully updated the count of product item on cart", body)
	ctx.JSON(http.StatusOK, response)
}

// UserCart godoc
// @summary api for get all cart item of user
// @description user can see all productItem that stored in cart
// @security ApiKeyAuth
// @id User Cart
// @tags User Cart
// @Router /carts [get]
// @Success 200 {object} res.Response{} "successfully got user cart items"
// @Failure 500 {object} res.Response{} "faild to get cart items"
func (u *cartHandler) ShowCart(ctx *gin.Context) {

	userId := utils.GetUserIdFromContext(ctx)

	// first get cart of user
	cart, err := u.carUseCase.GetUserCart(ctx, userId)
	if err != nil {
		response := res.ErrorResponse(500, "faild to get user cart", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	// user have not cart created
	if cart.CartID == 0 {
		respone := res.SuccessResponse(200, "user didn't add any product in cart", nil)
		ctx.JSON(http.StatusOK, respone)
		return
	}

	// get user cart items
	cartItems, err := u.carUseCase.GetUserCartItems(ctx, cart.CartID)
	if err != nil {
		response := res.ErrorResponse(500, "faild to get cart items", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if cartItems == nil {
		response := res.SuccessResponse(200, "there is no productItems in the cart", nil)
		ctx.JSON(http.StatusOK, response)
		return
	}

	resposeCart := res.Cart{
		CartItems:       cartItems,
		AppliedCouponID: cart.AppliedCouponID,
		TotalPrice:      cart.TotalPrice,
		DiscountAmount:  cart.DiscountAmount,
	}

	response := res.SuccessResponse(200, "successfully got user cart items", resposeCart)
	ctx.JSON(http.StatusOK, response)
}
