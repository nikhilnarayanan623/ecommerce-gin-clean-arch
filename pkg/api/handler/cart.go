package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/response"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase"
	usecaseInterface "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
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
//
//	@Summary		Add product item to cart (User)
//	@Description	API for user to add a product item to cart
//	@Security		BearerAuth
//	@Id				AddToCart
//	@Tags			User Cart
//	@Param			product_item_id	path	int	true	"Product Item ID"
//	@Router			/carts/{product_item_id} [post]
//	@Success		200	{object}	response.Response{}	"Successfully product item added to cart"
//	@Failure		404	{object}	response.Response{}	"Product item in out of stock"
//	@Failure		409	{object}	response.Response{}	"Product item already exist in cart"
//	@Failure		500	{object}	response.Response{}	"Failed to add product item into cart"
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
		switch {
		case errors.Is(err, usecase.ErrProductItemOutOfStock):
			statusCode = http.StatusNotFound
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
//
//	@Summary		Remove product item from cart (User)
//	@Description	API for user to remove a product item from cart
//	@Security		BearerAuth
//	@Id				RemoveFromCart
//	@Tags			User Cart
//	@Param			product_item_id	path	int	true	"Product Item ID"
//	@Router			/carts/{product_item_id} [delete]
//	@Success		200	{object}	response.Response{}	"Successfully product item removed form cart"
//	@Failure		400	{object}	response.Response{}	"invalid input"
//	@Failure		404	{object}	response.Response{}	"Product item not exist in cart"
//	@Failure		500	{object}	response.Response{}	"Failed to remove product item from cart"
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
			statusCode = http.StatusNotFound
		}

		response.ErrorResponse(ctx, statusCode, "Failed to remove product item from cart", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully product item removed form cart")
}

// UpdateCart godoc
//
//	@Summary		Change Cart Qty (User)
//	@Description	API for user to update cart item quantity (minimum qty is 1)
//	@Security		BearerAuth
//	@Id				UpdateCart
//	@Tags			User Cart
//	@Param			input	body	request.UpdateCartItem{}	true	"Input Field"
//	@Router			/carts [put]
//	@Success		200	{object}	response.Response{}	"Successfully to update cart item quantity changed in cart"
//	@Failure		400	{object}	response.Response{}	"Invalid input"
//	@Failure		500	{object}	response.Response{}	"Failed to update product item in cart"
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

	response.SuccessResponse(ctx, http.StatusOK, "Successfully to update cart item quantity changed in cart")
}

// GetCart godoc
//
//	@Summary		Get cart Items (User)
//	@Description	API for user to get all cart items
//	@Security		BearerAuth
//	@Id				GetCart
//	@Tags			User Cart
//	@Router			/carts [get]
//	@Success		200	{object}	response.Response{}	"Successfully retrieved all cart items"
//	@Success		204	{object}	response.Response{}	"Cart is empty"
//	@Failure		500	{object}	response.Response{}	"Failed to get user cart"
func (u *cartHandler) GetCart(ctx *gin.Context) {

	userId := utils.GetUserIdFromContext(ctx)

	// first get cart of user
	cart, err := u.carUseCase.GetUserCart(ctx, userId)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get user cart", err, nil)
	}

	// user have not cart created
	if cart.ID == 0 {
		response.SuccessResponse(ctx, http.StatusNoContent, "User cart is empty")
		return
	}

	// get user cart items
	cartItems, err := u.carUseCase.GetUserCartItems(ctx, cart.ID)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get cart items", err, nil)
		return
	}

	if len(cartItems) == 0 {
		response.SuccessResponse(ctx, http.StatusNoContent, "User cart is empty")
		return
	}

	responseCart := response.Cart{
		CartItems:       cartItems,
		AppliedCouponID: cart.AppliedCouponID,
		TotalPrice:      cart.TotalPrice,
		DiscountAmount:  cart.DiscountAmount,
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully retrieved all cart items", responseCart)
}
