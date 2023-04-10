package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/auth"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	service "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/varify"
)

type UserHandler struct {
	userUseCase service.UserUseCase
}

func NewUserHandler(userUsecase interfaces.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase: userUsecase}
}

// UserSignUp godoc
// @summary api for user to signup
// @security ApiKeyAuth
// @id UserSignUp
// @tags User Signup
// @Param input body req.ReqUserDetails{} true "Input Fields"
// @Router /signup [post]
// @Success 200 "Successfully account created for user"
// @Failure 400 "invalid input"
func (u *UserHandler) UserSignUp(ctx *gin.Context) {

	var body req.ReqUserDetails

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response := res.ErrorResponse(400, "invalid input", err.Error(), body)

		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	var user domain.User
	copier.Copy(&user, body)

	if err := u.userUseCase.Signup(ctx, user); err != nil {
		response := res.ErrorResponse(400, "faild to signup", err.Error(), body)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := res.SuccessResponse(200, "Successfully Account Created", body)
	ctx.JSON(200, response)
}

// UserLogin godoc
// @summary api for user to login
// @description Enter user_name | phone | email with password
// @security ApiKeyAuth
// @tags User Login
// @id UserLogin
// @Param        inputs   body     req.LoginStruct{}   true  "Input Field"
// @Router /login [post]
// @Success 200 {object} res.Response{} "successfully logged in"
// @Failure 400 {object} res.Response{}  "invalid input"
// @Failure 500 {object} res.Response{}  "faild to generat JWT"
func (u *UserHandler) UserLogin(ctx *gin.Context) {

	var body req.LoginStruct
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response := res.ErrorResponse(400, "invalid input", err.Error(), nil)

		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	//check all input field is empty
	if body.Email == "" && body.Phone == "" && body.UserName == "" {
		err := errors.New("enter atleast user_name or email or phone")
		response := res.ErrorResponse(400, "invalid input", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	//copy the body values to user
	var user domain.User
	copier.Copy(&user, &body)
	// get user from database and check password in usecase
	user, err := u.userUseCase.Login(ctx, user)
	if err != nil {
		response := res.ErrorResponse(400, "faild to login", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	// generate token using jwt in map
	tokenString, err := auth.GenerateJWT(user.ID)
	if err != nil {
		response := res.ErrorResponse(500, "faild to login", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	ctx.SetCookie("user-auth", tokenString["accessToken"], 60*60, "", "", false, true)

	response := res.SuccessResponse(200, "successfully logged in", tokenString["accessToken"])
	ctx.JSON(http.StatusOK, response)
}

// UserLoginOtpSend godoc
// @summary api for user to login with otp
// @description user can enter email/user_name/phone will send an otp to user registered phone_number
// @security ApiKeyAuth
// @id UserLoginOtpSend
// @tags User Login
// @Param inputs body req.OTPLoginStruct true "Input Field"
// @Router /login/otp-send [post]
// @Success 200 {object} res.Response{}  "Successfully Otp Send to registered number"
// @Failure 400 {object} res.Response{}  "Enter input properly"
// @Failure 500 {object} res.Response{}  "Faild to send otp"
func (u *UserHandler) UserLoginOtpSend(ctx *gin.Context) {

	var body req.OTPLoginStruct
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response := res.ErrorResponse(400, "invalid input", err.Error(), body)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	//check all input field is empty
	if body.Email == "" && body.Phone == "" && body.UserName == "" {
		err := errors.New("enter atleast user_name or email or phone")
		response := res.ErrorResponse(400, "invalid input", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	var user domain.User
	copier.Copy(&user, body)

	user, err := u.userUseCase.LoginOtp(ctx, user)

	if err != nil {
		resopnse := res.ErrorResponse(400, "can't login", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, resopnse)
		return
	}

	// if no error then send the otp
	if _, err := varify.TwilioSendOTP("+91" + user.Phone); err != nil {
		response := res.ErrorResponse(500, "faild to send otp", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	response := res.SuccessResponse(200, "successfully otp send to registered number", user.ID)
	ctx.JSON(http.StatusOK, response)
}

// UserLoginOtpVerify godoc
// @summary api for user to varify user login_otp
// @description enter your otp that send to your registered number
// @security ApiKeyAuth
// @id UserLoginOtpVerify
// @tags User Login
// @param inputs body req.OTPVerifyStruct{} true "Input Field"
// @Router /login/otp-verify [post]
// @Success 200 "successfully logged in uing otp"
// @Failure 400 "invalid login_otp"
// @Failure 500 "Faild to generate JWT"
func (u *UserHandler) UserLoginOtpVerify(ctx *gin.Context) {

	var body req.OTPVerifyStruct
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response := res.ErrorResponse(400, "invalid login_otp", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	var user = domain.User{
		ID: body.UserID,
	}

	// get the user using loginOtp useCase
	user, err := u.userUseCase.LoginOtp(ctx, user)
	if err != nil {
		response := res.ErrorResponse(400, "faild to login", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	// then varify the otp
	err = varify.TwilioVerifyOTP("+91"+user.Phone, body.OTP)
	if err != nil {
		response := res.ErrorResponse(400, "faild to login", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	// if everyting ok then generate token
	tokenString, err := auth.GenerateJWT(user.ID)
	if err != nil {
		response := res.ErrorResponse(500, "faild to login", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	ctx.SetCookie("user-auth", tokenString["accessToken"], 50*60, "", "", false, true)
	response := res.SuccessResponse(200, "successfully logged in uing otp", tokenString["accessToken"])
	ctx.JSON(http.StatusOK, response)
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

// AddToCart godoc
// @summary api for add productItem to user cart
// @description user can add a stock in product to user cart
// @security ApiKeyAuth
// @id AddToCart
// @tags User Cart
// @Param input body req.ReqCart true "Input Field"
// @Router /carts [post]
// @Success 200 "Successfully productItem added to cart"
// @Failure 400 "can't add the product item into cart"
func (u *UserHandler) AddToCart(ctx *gin.Context) {

	var body req.ReqCart
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response := res.ErrorResponse(400, "invalid input", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	// get userId and add to body
	body.UserID = utils.GetUserIdFromContext(ctx)

	err := u.userUseCase.SaveToCart(ctx, body)

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
// @Param input body req.ReqCart{} true "Input Field"
// @Router /carts [delete]
// @Success 200 {object} res.Response{} "Successfully productItem removed from cart"
// @Failure 400 {object} res.Response{}  "invalid input"
// @Failure 500 {object} res.Response{}  "can't remove product item from cart"
func (u UserHandler) RemoveFromCart(ctx *gin.Context) {

	var body req.ReqCart
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response := res.ErrorResponse(400, "invalid input", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	body.UserID = utils.GetUserIdFromContext(ctx)

	err := u.userUseCase.RemoveCartItem(ctx, body)

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
// @Param input body req.ReqCartCount{} true "Input Field"
// @Router /carts [put]
// @Success 200 "Successfully productItem count change on cart"
// @Failure 400  "invalid input"
// @Failure 500  "can't update the count of product item on cart"
func (u *UserHandler) UpdateCart(ctx *gin.Context) {

	var body req.ReqCartCount

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response := res.ErrorResponse(400, "invalid input", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	body.UserID = utils.GetUserIdFromContext(ctx)

	err := u.userUseCase.UpdateCartItem(ctx, body)

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
func (u *UserHandler) UserCart(ctx *gin.Context) {

	userId := utils.GetUserIdFromContext(ctx)

	// first get cart of user
	cart, err := u.userUseCase.GetUserCart(ctx, userId)
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
	cartItems, err := u.userUseCase.GetUserCartItems(ctx, cart.CartID)
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

	resposeCart := res.ResCart{
		CartItems:       cartItems,
		AppliedCouponID: cart.AppliedCouponID,
		TotalPrice:      cart.TotalPrice,
		DiscountAmount:  cart.DiscountAmount,
	}

	response := res.SuccessResponse(200, "successfully got user cart items", resposeCart)
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

	var data res.UserRespStrcut
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
// @Param input body req.ReqUserDetails{} true "input field"
// @Router /account [put]
// @Success 200 {object} res.Response{} "successfully updated user details"
// @Failure 400 {object} res.Response{} "invalid input"
func (u *UserHandler) UpateAccount(ctx *gin.Context) {
	userID := utils.GetUserIdFromContext(ctx)

	var body req.ReqUserDetails

	if err := ctx.ShouldBindJSON(&body); err != nil { // showing epty struct which is user for know what are the fields need enter
		response := res.ErrorResponse(400, "invalid input", err.Error(), req.ReqUserDetails{})
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

	response := res.SuccessResponse(200, "successfully updated user details", body)
	ctx.JSON(http.StatusOK, response)
}

// AddAddress godoc
// @summary api for adding a new address for user
// @description get a new address from user to store the the database
// @security ApiKeyAuth
// @id AddAddress
// @tags User Address
// @Param inputs body req.ReqAddress{} true "Input Field"
// @Router /account/address [post]
// @Success 200 {object} res.Response{} "Successfully address added"
// @Failure 400 {object} res.Response{} "inavlid input"
func (u *UserHandler) AddAddress(ctx *gin.Context) {

	var body req.ReqAddress
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
// @Param input body req.ReqEditAddress{} true "Input Field"
// @Router /account/address [put]
// @Success 200 {object} res.Response{} "successfully addresses updated"
// @Failure 400 {object} res.Response{} "can't update the address"
func (u *UserHandler) EditAddress(ctx *gin.Context) {

	userID := utils.GetUserIdFromContext(ctx)
	var body req.ReqEditAddress

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
