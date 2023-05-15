package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	handlerInterface "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	service "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
)

type UserHandler struct {
	userUseCase service.UserUseCase
}

func NewUserHandler(userUsecase interfaces.UserUseCase) handlerInterface.UserHandler {
	return &UserHandler{userUseCase: userUsecase}
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

	response := res.SuccessResponse(200, "welcome to home page", nil)
	ctx.JSON(http.StatusOK, response)
}

// Logout godoc
// @summary api for user to lgout
// @description user can logout
// @security ApiKeyAuth
// @id UserLogout
// @tags User Logout
// @Router /logout [post]
// @Success 200 "successfully logged out"
func (u *UserHandler) UserLogout(ctx *gin.Context) {
	ctx.SetCookie("user-auth", "", -1, "", "", false, true)
	response := res.SuccessResponse(200, "successfully logged out", nil)
	ctx.JSON(http.StatusOK, response)
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
// @Failure 500 {object} res.Response{} "faild to get checkout items"
func (c *UserHandler) CheckOutCart(ctx *gin.Context) {

	// userId := utils.GetUserIdFromContext(ctx)

	// resCheckOut, err := c.userUseCase.CheckOutCart(ctx, userId)

	// if err != nil {
	// 	response := res.ErrorResponse(500, "faild to get checkout items", err.Error(), nil)
	// 	ctx.AbortWithStatusJSON(http.StatusInternalServerError, response)
	// 	return
	// }

	// if resCheckOut.ProductItems == nil {
	// 	response := res.ErrorResponse(401, "cart is empty can't checkout cart", "", nil)
	// 	ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
	// 	return
	// }

	// responser := res.SuccessResponse(200, "successfully got checkout data", resCheckOut)
	// ctx.JSON(http.StatusOK, responser)
}

// ! ***** for user account ***** //
// Account godoc
// @summary api for see use details
// @security ApiKeyAuth
// @id Account
// @tags User Account
// @Router /account [get]
// @Success 200 "Successfully user account details found"
// @Failure 500 {object} res.Response{} "faild to show user details"
func (u *UserHandler) Account(ctx *gin.Context) {

	userID := utils.GetUserIdFromContext(ctx)

	user, err := u.userUseCase.Account(ctx, userID)
	if err != nil {
		response := res.ErrorResponse(500, "faild to show user details", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	var data res.User
	copier.Copy(&data, &user)

	response := res.SuccessResponse(200, "Successfully user account details found", data)
	ctx.JSON(http.StatusOK, response)
}

// UpateAccount godoc
// @summary api for edit user details
// @description user can edit user details
// @security ApiKeyAuth
// @id UpateAccount
// @tags User Account
// @Param input body req.EditUser{} true "input field"
// @Router /account [put]
// @Success 200 {object} res.Response{} "successfully updated user details"
// @Failure 400 {object} res.Response{} "invalid input"
func (u *UserHandler) UpateAccount(ctx *gin.Context) {
	userID := utils.GetUserIdFromContext(ctx)

	var body req.EditUser

	if err := ctx.ShouldBindJSON(&body); err != nil { // showing epty struct which is user for know what are the fields need enter
		response := res.ErrorResponse(400, "invalid input", err.Error(), body)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	var user domain.User

	copier.Copy(&user, &body)

	user.ID = userID
	// edit the user details
	if err := u.userUseCase.EditAccount(ctx, user); err != nil {
		response := res.ErrorResponse(400, "invalid input", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := res.SuccessResponse(200, "successfully updated user details", nil)
	ctx.JSON(http.StatusOK, response)
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
// @Failure 400 {object} res.Response{} "inavlid input"
func (u *UserHandler) AddAddress(ctx *gin.Context) {

	var body req.Address
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response := res.ErrorResponse(400, "inavlid input", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
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
		response := res.ErrorResponse(400, "inavlid input", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := res.SuccessResponse(200, "successfully saved user address", body)
	ctx.JSON(http.StatusOK, response)
}

// GetAddreses godoc
// @summary api for get all address of user
// @description user can show all adderss
// @security ApiKeyAuth
// @id GetAddresses
// @tags User Address
// @Router /account/address [get]
// @Success 200 {object} res.Response{} "successfully got user addresses"
// @Failure 500 {object} res.Response{} "faild to show user addresses"
func (u *UserHandler) GetAddresses(ctx *gin.Context) {

	userID := utils.GetUserIdFromContext(ctx)

	addresses, err := u.userUseCase.GetAddresses(ctx, userID)

	if err != nil {
		response := res.ErrorResponse(500, "faild to show user addresses", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if addresses == nil {
		response := res.SuccessResponse(200, "there is no addresses to show")
		ctx.JSON(http.StatusOK, response)
		return
	}

	response := res.SuccessResponse(200, "successfully got user addresses", addresses)
	ctx.JSON(http.StatusOK, response)
}

// EditAddress godoc
// @summary api for edit user address
// @description user can change existing address
// @security ApiKeyAuth
// @id EditAddress
// @tags User Address
// @Param input body req.EditAddress{} true "Input Field"
// @Router /account/address [put]
// @Success 200 {object} res.Response{} "successfully addresses updated"
// @Failure 400 {object} res.Response{} "can't update the address"
func (u *UserHandler) EditAddress(ctx *gin.Context) {

	userID := utils.GetUserIdFromContext(ctx)
	var body req.EditAddress

	if err := ctx.ShouldBindJSON(&body); err != nil {
		respone := res.ErrorResponse(400, "invalid input", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, respone)
		return
	}

	// address is_default reference pointer need to change in future
	if body.IsDefault == nil {
		body.IsDefault = new(bool)
	}

	err := u.userUseCase.EditAddress(ctx, body, userID)
	if err != nil {
		response := res.ErrorResponse(400, "faild to update user address", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	reponse := res.SuccessResponse(200, "successfully addresses updated", body)
	ctx.JSON(http.StatusOK, reponse)

}

func (u *UserHandler) DeleteAddress(ctx *gin.Context) {

}

//todo ** wishList **

// AddToWishList godoc
// @summary api to add a productItem to wish list
// @descritpion user can add productItem to wish list
// @security ApiKeyAuth
// @id AddToWishList
// @tags Wishlist
// @Param product_id body int true "product_id"
// @Router /wishlist [post]
// @Success 200 {object} res.Response{} "successfully added product item to wishlist"
// @Failure 400 {object} res.Response{} "invalid input"
func (u *UserHandler) AddToWishList(ctx *gin.Context) {
	// get productItemID using parmas
	productItemID, err := utils.StringToUint(ctx.Param("id"))

	if err != nil {
		reponse := res.ErrorResponse(400, "invalid input", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, reponse)
		return
	}

	userID := utils.GetUserIdFromContext(ctx)

	var wishList = domain.WishList{
		ProductItemID: productItemID,
		UserID:        userID,
	}
	fmt.Println(wishList.UserID, wishList.ProductItemID)

	// add to wishlist
	if err := u.userUseCase.AddToWishList(ctx, wishList); err != nil {
		response := res.ErrorResponse(400, "fail to add product on wishlist", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	response := res.SuccessResponse(200, "successfully added product item to wishlist", nil)
	ctx.JSON(http.StatusOK, response)
}

// RemoveFromWishList godoc
// @summary api to remove a productItem from wish list
// @descritpion user can remove a productItem from wish list
// @security ApiKeyAuth
// @id RemoveFromWishList
// @tags Wishlist
// @Params product_item_id path int true "product_item_id"
// @Router /wishlist [delete]
// @Success 200 {object} res.Response{} "successfully removed product item from wishlist"
// @Failure 400 {object} res.Response{} "invalid input"
func (u *UserHandler) RemoveFromWishList(ctx *gin.Context) {

	// get productItemID using parmas
	productItemID, err := utils.StringToUint(ctx.Param("id"))

	if err != nil {
		response := res.ErrorResponse(400, "invalid input", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	userID := utils.GetUserIdFromContext(ctx)

	var wishList = domain.WishList{
		ProductItemID: productItemID,
		UserID:        userID,
	}

	// remove form wishlist
	if err := u.userUseCase.RemoveFromWishList(ctx, wishList); err != nil {
		response := res.ErrorResponse(400, "faild to remove product item from wishlist", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := res.SuccessResponse(200, "successfully removed product item from wishlist", nil)
	ctx.JSON(http.StatusOK, response)
}

// GetWishListI godoc
// @summary api get all wish list items of user
// @descritpion user get all wish list items
// @security ApiKeyAuth
// @id GetWishListI
// @tags Wishlist
// @Router /wishlist [get]
// @Success 200 "Successfully wish list items got"
// @Success 200 "Wish list is empty"
// @Failure 400  "faild to get user wish list items"
func (u *UserHandler) GetWishListI(ctx *gin.Context) {

	userID := utils.GetUserIdFromContext(ctx)
	data, err := u.userUseCase.GetWishListItems(ctx, userID)

	if err != nil {
		response := res.ErrorResponse(500, "faild to get wish list item", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	if data == nil {
		response := res.SuccessResponse(200, "wish list is empty", nil)
		ctx.JSON(http.StatusOK, response)
		return
	}

	response := res.SuccessResponse(200, "successfully got wish list item", data)
	ctx.JSON(http.StatusOK, response)
}
