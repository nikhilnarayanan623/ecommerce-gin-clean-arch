package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	usecaseInterface "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/response"
)

type UserHandler struct {
	userUseCase usecaseInterface.UserUseCase
}

func NewUserHandler(userUsecase usecaseInterface.UserUseCase) interfaces.UserHandler {
	return &UserHandler{
		userUseCase: userUsecase,
	}
}

// Home godoc
// @summary api for showing home page of user
// @description after user login user will seen this page with user informations
// @security ApiKeyAuth
// @id User Home
// @tags Home
// @Router / [get]
// @Success 200 "Welcome Home"
func (u *UserHandler) Home(ctx *gin.Context) {

	response.SuccessResponse(ctx, http.StatusOK, "welcome to home page", nil)
}

// Logout godoc
// @summary api for user to logout
// @description user can logout
// @security ApiKeyAuth
// @id UserLogout
// @tags User Logout
// @Router /logout [post]
// @Success 200 "successfully logged out"
func (u *UserHandler) UserLogout(ctx *gin.Context) {

	ctx.SetCookie("user-auth", "", -1, "", "", false, true)

	response.SuccessResponse(ctx, http.StatusOK, "Successfully logged out", nil)
}

// CheckOutCart godoc
// @summary api for cart checkout
// @description user can checkout user cart items
// @security ApiKeyAuth
// @id CheckOutCart
// @tags User Cart
// @Router /carts/checkout [get]
// @Success 200 {object} res.Response{} "successfully got checkout data"
// @Failure 401 {object} res.Response{} "cart is empty so user can't call this api"
// @Failure 500 {object} res.Response{} "failed to get checkout items"
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

// FindUserProfile godoc
// @summary api for see use details
// @security ApiKeyAuth
// @id FindUserProfile
// @tags User GetUserProfile
// @Router /account [get]
// @Success 200 "Successfully user account details found"
// @Failure 500 {object} res.Response{} "faild to show user details"
func (u *UserHandler) FindProfile(ctx *gin.Context) {

	userID := utils.GetUserIdFromContext(ctx)

	user, err := u.userUseCase.FindProfile(ctx, userID)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find user profile details", err, nil)
		return
	}

	var data response.User
	copier.Copy(&data, &user)

	response.SuccessResponse(ctx, http.StatusOK, "Successfully found user profile details", data)
}

// UpdateUserProfile godoc
// @summary api for edit user details
// @description user can edit user details
// @security ApiKeyAuth
// @id UpdateUserProfile
// @tags User Account
// @Param input body req.EditUser{} true "input field"
// @Router /account [put]
// @Success 200 {object} res.Response{} "successfully updated user details"
// @Failure 400 {object} res.Response{} "invalid input"
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

// AddAddress godoc
// @summary api for adding a new address for user
// @description get a new address from user to store the the database
// @security ApiKeyAuth
// @id AddAddress
// @tags User Address
// @Param inputs body req.Address{} true "Input Field"
// @Router /account/address [post]
// @Success 200 {object} res.Response{} "Successfully address added"
// @Failure 400 {object} res.Response{} "invalid input"
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

// FindAllAddresses godoc
// @summary api for get all address of user
// @description user can show all address
// @security ApiKeyAuth
// @id FindAllAddresses
// @tags User Address
// @Router /account/address [get]
// @Success 200 {object} res.Response{} "successfully got user addresses"
// @Failure 500 {object} res.Response{} "failed to show user addresses"
func (u *UserHandler) FindAllAddresses(ctx *gin.Context) {

	userID := utils.GetUserIdFromContext(ctx)

	addresses, err := u.userUseCase.FindAddresses(ctx, userID)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find user addresses", err, nil)
		return
	}

	if addresses == nil {
		response.SuccessResponse(ctx, http.StatusOK, "No addresses found")
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully found user addresses", addresses)
}

// UpdateAddress godoc
// @summary api for edit user address
// @description user can change existing address
// @security ApiKeyAuth
// @id UpdateAddress
// @tags User Address
// @Param input body req.EditAddress{} true "Input Field"
// @Router /account/address [put]
// @Success 200 {object} res.Response{} "successfully addresses updated"
// @Failure 400 {object} res.Response{} "can't update the address"
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

// func (u *UserHandler) DeleteAddress(ctx *gin.Context) {

// }

//todo ** wishList **

// AddToWishList godoc
// @summary api to add a productItem to wish list
// @descriptions user can add productItem to wish list
// @security ApiKeyAuth
// @id AddToWishList
// @tags Wishlist
// @Param product_item_id body int true "product_item_id"
// @Router /wishlist [post]
// @Success 200 {object} res.Response{} "successfully added product item to wishlist"
// @Failure 400 {object} res.Response{} "invalid input"
func (u *UserHandler) AddToWishList(ctx *gin.Context) {

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
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to add product on wishlist", err, nil)
		return
	}
	response.SuccessResponse(ctx, http.StatusCreated, "successfully added product item to wishlist", nil)
}

// RemoveFromWishList godoc
// @summary api to remove a productItem from wish list
// @descriptions user can remove a productItem from wish list
// @security ApiKeyAuth
// @id RemoveFromWishList
// @tags Wishlist
// @Params product_item_id path int true "product_item_id"
// @Router /wishlist [delete]
// @Success 200 {object} res.Response{} "successfully removed product item from wishlist"
// @Failure 400 {object} res.Response{} "invalid input"
func (u *UserHandler) RemoveFromWishList(ctx *gin.Context) {

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

	// remove form wishlist
	if err := u.userUseCase.RemoveFromWishList(ctx, wishList); err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to remove product item from wishlist", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully removed product item from wishlist", nil)
}

// FindWishList godoc
// @summary api get all wish list items of user
// @descriptions user get all wish list items
// @security ApiKeyAuth
// @id FindWishList
// @tags Wishlist
// @Router /wishlist [get]
// @Success 200 "Successfully wish list items got"
// @Success 200 "Wish list is empty"
// @Failure 400  "failed to get user wish list items"
func (u *UserHandler) FindWishList(ctx *gin.Context) {

	userID := utils.GetUserIdFromContext(ctx)

	wishListItems, err := u.userUseCase.FindAllWishListItems(ctx, userID)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find wishlist item", err, nil)
		return
	}

	if len(wishListItems) == 0 {
		response.SuccessResponse(ctx, http.StatusOK, "No wishlist items found", nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "successfully got wish list item", wishListItems)
}
