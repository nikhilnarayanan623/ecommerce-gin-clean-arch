package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/response"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase"
	usecaseInterface "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
)

type UserHandler struct {
	userUseCase usecaseInterface.UserUseCase
}

func NewUserHandler(userUsecase usecaseInterface.UserUseCase) interfaces.UserHandler {
	return &UserHandler{
		userUseCase: userUsecase,
	}
}

// // Logout godoc
// // @summary api for user to logout
// // @description user can logout
// // @security ApiKeyAuth
// // @id UserLogout
// // @tags User Logout
// // @Router /logout [post]
// // @Success 200 "successfully logged out"
// func (u *UserHandler) UserLogout(ctx *gin.Context) {

// 	ctx.SetCookie("user-auth", "", -1, "", "", false, true)

// 	response.SuccessResponse(ctx, http.StatusOK, "Successfully logged out", nil)
// }

// // CheckOutCart godoc
// // @summary api for cart checkout
// // @description user can checkout user cart items
// // @Security BearerAuth
// // @id CheckOutCart
// // @tags User Cart
// // @Router /carts/checkout [get]
// // @Success 200 {object} response.Response{} "successfully got checkout data"
// // @Failure 401 {object} res.Response{} "cart is empty so user can't call this api"
// // @Failure 500 {object} res.Response{} "failed to get checkout items"
func (c *UserHandler) CheckOutCart(ctx *gin.Context) {

	// userId := utils.GetUserIdFromContext(ctx)

	// resCheckOut, err := c.userUseCase.CheckOutCart(ctx, userId)

	// if err != nil {
	// 	 response.ErrorResponse(500, "failed to get checkout items", err.Error(), nil)
	// 	ctx.AbortWithStatusJSON(http.StatusInternalServerError, response)
	// 	return
	// }

	// if resCheckOut.ProductItems == nil {
	// 	 response.ErrorResponse(401, "cart is empty can't checkout cart", "", nil)
	// 	ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
	// 	return
	// }

	// responser := res.SuccessResponse(200, "successfully got checkout data", resCheckOut)
	// ctx.JSON(http.StatusOK, responser)
}

// GetProfile godoc
//
//	@Summary		Get User Profile (User)
//	@Security		BearerAuth
//	@Description	API for user to get all user details
//	@Id				GetProfile
//	@Tags			User Profile
//	@Router			/account [get]
//	@Success		200	"Successfully retrieved user details"
//	@Failure		500	{object}	response.Response{}	"Failed to retrieve user details"
func (u *UserHandler) GetProfile(ctx *gin.Context) {

	userID := utils.GetUserIdFromContext(ctx)

	user, err := u.userUseCase.FindProfile(ctx, userID)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve user details", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully retrieved user details", user)
}

// UpdateProfile godoc
//
//	@Summary		Edit profile (User)
//	@Security		BearerAuth
//	@Description	API for user to edit user details
//	@Id				UpdateProfile
//	@Tags			User Profile
//	@Param			input	body	request.EditUser{}	true	"User details input"
//	@Router			/account [put]
//	@Success		200	{object}	response.Response{}	"Successfully profile updated"
//	@Failure		400	{object}	response.Response{}	"Invalid inputs"
//	@Failure		500	{object}	response.Response{}	"Failed to update profile"
func (u *UserHandler) UpdateProfile(ctx *gin.Context) {

	userID := utils.GetUserIdFromContext(ctx)

	var body request.EditUser

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	var user domain.User
	copier.Copy(&user, &body)
	user.ID = userID

	err := u.userUseCase.UpdateProfile(ctx, user)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update profile", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully profile updated", nil)
}

// SaveAddress godoc
//
//	@Summary		Add a new address (User)
//	@Security		BearerAuth
//	@Description	API for user to add a new address
//	@Id				SaveAddress
//	@Tags			User Profile
//	@Param			inputs	body	request.Address{}	true	"Address input"
//	@Router			/account/address [post]
//	@Success		200	{object}	response.Response{}	"Successfully address added"
//	@Failure		400	{object}	response.Response{}	"invalid input"
//	@Failure		500	{object}	response.Response{}	"Failed to save address"
func (u *UserHandler) SaveAddress(ctx *gin.Context) {

	var body request.Address
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	userID := utils.GetUserIdFromContext(ctx)

	var address domain.Address

	copier.Copy(&address, &body)

	// check is default is null
	if body.IsDefault == nil {
		body.IsDefault = new(bool)
	}

	err := u.userUseCase.SaveAddress(ctx, userID, address, *body.IsDefault)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to save address", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusCreated, "Successfully address saved")
}

// GetAllAddresses godoc
//
//	@Summary		Get all addresses (User)
//	@Security		BearerAuth
//	@Description	API for user to get all user addresses
//	@Id				GetAllAddresses
//	@Tags			User Profile
//	@Router			/account/address [get]
//	@Success		200	{object}	response.Response{}	"successfully retrieved all user addresses"
//	@Failure		500	{object}	response.Response{}	"failed to show user addresses"
func (u *UserHandler) GetAllAddresses(ctx *gin.Context) {

	userID := utils.GetUserIdFromContext(ctx)

	addresses, err := u.userUseCase.FindAddresses(ctx, userID)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get user addresses", err, nil)
		return
	}

	if addresses == nil {
		response.SuccessResponse(ctx, http.StatusOK, "No addresses found")
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully retrieved all user addresses", addresses)
}

// UpdateAddress godoc
//
//	@Summary		Update address (User)
//	@Security		BearerAuth
//	@Description	API for user to update user address
//	@Id				UpdateAddress
//	@Tags			User Profile
//	@Param			input	body	request.EditAddress{}	true	"Address input"
//	@Router			/account/address [put]
//	@Success		200	{object}	response.Response{}	"successfully addresses updated"
//	@Failure		400	{object}	response.Response{}	"can't update the address"
func (u *UserHandler) UpdateAddress(ctx *gin.Context) {

	userID := utils.GetUserIdFromContext(ctx)
	var body request.EditAddress

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	// address is_default reference pointer need to change in future
	if body.IsDefault == nil {
		body.IsDefault = new(bool)
	}

	err := u.userUseCase.UpdateAddress(ctx, body, userID)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update user address", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "successfully addresses updated", body)

}

// SaveToWishList godoc
//
//	@Summary		Add to whish list (User)
//	@Security		BearerAuth
//	@Descriptions	API for user to add product item to wish list
//	@Id				SaveToWishList
//	@Tags			User Profile
//	@Param			product_item_id	path	int	true	"Product Item ID"
//	@Router			/account/wishlist/{product_item_id} [post]
//	@Success		200	{object}	response.Response{}	"Successfully product items added to whish list"
//	@Failure		400	{object}	response.Response{}	"invalid input"
//	@Failure		409	{object}	response.Response{}	"Product item already exist on wish list"
//	@Failure		500	{object}	response.Response{}	"Failed to add product item to wishlist"
func (u *UserHandler) SaveToWishList(ctx *gin.Context) {

	productItemID, err := request.GetParamAsUint(ctx, "product_item_id")

	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindParamFailMessage, err, nil)
		return
	}

	userID := utils.GetUserIdFromContext(ctx)

	var wishList = domain.WishList{
		ProductItemID: productItemID,
		UserID:        userID,
	}

	err = u.userUseCase.SaveToWishList(ctx, wishList)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, usecase.ErrExistWishListProductItem) {
			statusCode = http.StatusConflict
		}
		response.ErrorResponse(ctx, statusCode, "Failed to add product item to wishlist", err, nil)
		return
	}
	response.SuccessResponse(ctx, http.StatusCreated, "Successfully product items added to whish list", nil)
}

// RemoveFromWishList godoc
//
//	@Summary		Remove from whish list (User)
//	@Security		BearerAuth
//	@Descriptions	API for user to remove a product item from whish list
//	@Id				RemoveFromWishList
//	@Tags			User Profile
//	@Param			product_item_id	path	int	true	"Product Item ID"
//	@Router			/account/wishlist/{product_item_id} [delete]
//	@Success		200	{object}	response.Response{}	"successfully removed product item from wishlist"
//	@Failure		400	{object}	response.Response{}	"invalid input"
func (u *UserHandler) RemoveFromWishList(ctx *gin.Context) {

	productItemID, err := request.GetParamAsUint(ctx, "product_item_id")

	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindParamFailMessage, err, nil)
		return
	}

	userID := utils.GetUserIdFromContext(ctx)

	// remove form wishlist
	if err := u.userUseCase.RemoveFromWishList(ctx, userID, productItemID); err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to remove product item from wishlist", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully removed product item from wishlist", nil)
}

// GetWishList godoc
//
//	@Summary		Get whish list product items (User)
//	@Security		BearerAuth
//	@Descriptions	API for user to get product items in the wish list
//	@Id				GetWishList
//	@Tags			User Profile
//	@Router			/account/wishlist [get]
//	@Success		200	"Successfully retrieved all product items in th wish list"
//	@Failure		500	"Failed to retrieve product items from the wish list"
func (u *UserHandler) GetWishList(ctx *gin.Context) {

	userID := utils.GetUserIdFromContext(ctx)

	wishListItems, err := u.userUseCase.FindAllWishListItems(ctx, userID)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve product items from the wish list", err, nil)
		return
	}

	if len(wishListItems) == 0 {
		response.SuccessResponse(ctx, http.StatusOK, "No wishlist items found", nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully retrieved all product items in th wish list", wishListItems)
}
