package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase"
	usecaseInterface "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/response"
)

type cartHandler struct {
	carUseCase usecaseInterface.CartUseCase
}

func NewCartHandler(cartUseCase usecaseInterface.CartUseCase) interfaces.CartHandler {
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
// @Param product_item_id path  int true "Product Item ID"
// @Router /carts/{product_item_id} [post]
// @Success 200 {object} response.Response{} "Successfully product item added to cart"
// @Failure 400 {object} response.Response{} "Failed to add product item into cart"
func (u *cartHandler) AddToCart(ctx *gin.Context) {

	productItemID, err := request.GetParamAsUint(ctx, "product_item_id")
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindParamFailMessage, err, nil)
		return
	}

	userID := utils.GetUserIdFromContext(ctx)

	err = u.carUseCase.SaveProductItemToCart(ctx, userID, productItemID)

	if err != nil {
		var statusCode int
		switch true {
		case errors.Is(err, usecase.ErrProductItemOutOfStock):
			statusCode = http.StatusBadRequest
		case errors.Is(err, usecase.ErrCartItemAlreadyExist):
			statusCode = http.StatusConflict
		default:
			statusCode = http.StatusInternalServerError
		}
		response.ErrorResponse(ctx, statusCode, "Failed to add product item into cart", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusCreated, "Successfully product item added to cart")
}

// RemoveFromCart godoc
// @summary api for remove a product from cart
// @description user can remove a signle productItem full quantity from cart
// @security ApiKeyAuth
// @id RemoveFromCart
// @tags User Cart
// @Param product_item_id path  int true "Product Item ID"
// @Router /carts/{product_item_id} [delete]
// @Success 200 {object} response.Response{} "Successfully product item removed form cart"
// @Failure 400 {object} response.Response{}  "invalid input"
// @Failure 500 {object} response.Response{}  "Failed to remove product item from cart"
func (u cartHandler) RemoveFromCart(ctx *gin.Context) {

	productItemID, err := request.GetParamAsUint(ctx, "product_item_id")
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindParamFailMessage, err, nil)
		return
	}

	userID := utils.GetUserIdFromContext(ctx)

	err = u.carUseCase.RemoveProductItemFromCartItem(ctx, userID, productItemID)

	if err != nil {

		statusCode := http.StatusInternalServerError

		if errors.Is(err, usecase.ErrCartItemNotExit) {
			statusCode = http.StatusBadRequest
		}

		response.ErrorResponse(ctx, statusCode, "Failed to remove product item from cart", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully product item removed form cart")
}

// UpdateCart godoc
// @summary api for update productItem count
// @description user can increment or decrement count of a productItem in cart (min=1)
// @security ApiKeyAuth
// @id UpdateCart
// @tags User Cart
// @Param input body request.UpdateCartItem{} true "Input Field"
// @Router /carts [put]
// @Success 200 "Successfully productItem count change on cart"
// @Failure 400 {object} response.Response{}   "invalid input"
// @Failure 500  {object} response.Response{}  "Failed to update product item in cart"
func (u *cartHandler) UpdateCart(ctx *gin.Context) {

	var body request.UpdateCartItem

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	body.UserID = utils.GetUserIdFromContext(ctx)

	err := u.carUseCase.UpdateCartItem(ctx, body)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update product item in cart", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully to update product item in cart")
}

// FindCart godoc
// @summary api for get all cart item of user
// @description user can see all productItem that stored in cart
// @security ApiKeyAuth
// @id FindCart
// @tags User Cart
// @Router /carts [get]
// @Success 200 {object} response.Response{} "Successfully find user cart items"
// @Failure 500 {object} response.Response{} "Failed to get user cart"
func (u *cartHandler) FindCart(ctx *gin.Context) {

	userId := utils.GetUserIdFromContext(ctx)

	// first get cart of user
	cart, err := u.carUseCase.GetUserCart(ctx, userId)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get user cart", err, nil)
	}

	// user have not cart created
	if cart.ID == 0 {
		response.SuccessResponse(ctx, http.StatusOK, "User cart is empty")
		return
	}

	// get user cart items
	cartItems, err := u.carUseCase.GetUserCartItems(ctx, cart.ID)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get cart items", err, nil)
		return
	}

	if len(cartItems) == 0 {
		response.SuccessResponse(ctx, http.StatusOK, "User cart is empty")
		return
	}

	responseCart := response.Cart{
		CartItems:       cartItems,
		AppliedCouponID: cart.AppliedCouponID,
		TotalPrice:      cart.TotalPrice,
		DiscountAmount:  cart.DiscountAmount,
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully find user cart items", responseCart)
}
