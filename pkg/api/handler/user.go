package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/auth"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/res"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	service "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/varify"
)

type UserHandler struct {
	userUseCase service.UserUseCase
}

func NewUserHandler(userUsecase interfaces.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase: userUsecase}
}

// SignUpGet godoc
// @summary api for user to signup page
// @description user can see what are the fields to enter to create a new account
// @security ApiKeyAuth
// @id SignUpGet
// @tags signup
// @produce json
// @Router /signup [get]
// @Success 200 {object} domain.User{} "OK"
func (u *UserHandler) SignUpGet(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"msg":        "enter detail for signup",
		"user":       domain.User{},
	})
}

// SignUpPost godoc
// @summary api for user to post the user details
// @description user can send user details and validate and create new account
// @security ApiKeyAuth
// @id SignUpPost
// @tags signup
// @produce json
// @Router /signup [post]
// @Success 200 "Successfully account created"
// @Failure 400 "Faild to create account"
func (u *UserHandler) SignUpPost(ctx *gin.Context) {
	var user domain.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "Cant't Bind The Values",
			"error":      err.Error(),
		})

		return
	}

	user, err := u.userUseCase.Signup(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "can't signup",
			"error":      err.Error(),
		})
		return
	}

	var response res.UserRespStrcut

	copier.Copy(&response, &user)

	ctx.JSON(200, gin.H{
		"StatusCode": 200,
		"msg":        "Successfully Account Created",
		"user":       response,
	})
}

func (u *UserHandler) Home(ctx *gin.Context) {

	userId := helper.GetUserIdFromContext(ctx)

	user, err := u.userUseCase.Home(ctx, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"StatusCode": 500,
			"error":      err.Error(),
		})
		return
	}

	var response res.UserRespStrcut
	copier.Copy(&response, &user)

	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"msg":        "Welcome Home",
		"user":       response,
	})
}

// LoginGet godoc
// @summary to get the json format for login
// @description Enter this fields on login page post
// @tags login
// @security ApiKeyAuth
// @id LoginGet
// @produce json
// @Router /login [get]
// @Success 200 {object} req.LoginStruct "OK"
func (u *UserHandler) LoginGet(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"msg":        "detail to enter",
		"user":       req.LoginStruct{},
	})
}

// LoginPost godoc
// @summary api for user login
// @description Enter user_name/phone/email with password
// @security ApiKeyAuth
// @tags login
// @id LoginPost
// @produce json
// @Param        inputs   body     req.LoginStruct{}   true  "Input Field"
// @Router /login [post]
// @Success 200 "Successfully Loged In"
// @Failure 400 "faild to login"
// @Failure 500 "faild to generat JWT"
func (u *UserHandler) LoginPost(ctx *gin.Context) {

	var body req.LoginStruct
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "Cant't bind the json",
			"error":      err.Error(),
		})
		return
	}

	//check all input field is empty
	if body.Email == "" && body.Phone == "" && body.UserName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "can't login",
			"error":      "Enter atleast user_name or email or phone",
		})
		return
	}
	//copy the body values to user
	var user domain.User
	copier.Copy(&user, &body)

	// get user from database and check password in usecase
	user, err := u.userUseCase.Login(ctx, user)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "can't login",
			"error":      err.Error(),
		})
		return
	}
	// generate token using jwt in map
	tokenString, err := auth.GenerateJWT(user.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"StatusCode": 500,
			"msg":        "can't login",
			"error":      "faild to generat JWT",
		})
		return
	}

	ctx.SetCookie("user-auth", tokenString["accessToken"], 10*60, "", "", false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"msg":        "Successfully Loged In",
	})
}

// LoginOtpSend godoc
// @summary api for user login with otp
// @description user can enter email/user_name/phone will send an otp to user phone
// @security ApiKeyAuth
// @id LoginOtpSend
// @tags login
// @produce json
// @Param inputs body req.OTPLoginStruct true "Input Field"
// @Router /login-otp-send [post]
// @Success 200 "Successfully Otp Send to registered number"
// @Failure 400 "Enter input properly"
// @Failure 500 "Faild to send otp"
func (u *UserHandler) LoginOtpSend(ctx *gin.Context) {

	var body req.OTPLoginStruct
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "Cant't bind the json",
			"error":      err.Error(),
		})
		return
	}

	//check all input field is empty
	if body.Email == "" && body.Phone == "" && body.UserName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "can't login",
			"error":      "Enter atleast user_name or email or phone",
		})
		return
	}

	var user domain.User
	copier.Copy(&user, body)

	user, err := u.userUseCase.LoginOtp(ctx, user)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "Can't login",
			"error":      err.Error(),
		})
		return
	}

	// if no error then send the otp
	if _, err := varify.TwilioSendOTP("+91" + user.Phone); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"StatusCode": 500,
			"msg":        "Can't login",
			"error":      "faild to sent OTP",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"msg":        "Successfully Otp Send to registered number",
		"userId":     user.ID,
	})

}

// LoginOtpVerify godoc
// @summary varify user login otp
// @description enter your otp that send to your registered number
// @security ApiKeyAuth
// @id LoginOtpVerify
// @tags login
// @produce json
// @param inputs body req.OTPVerifyStruct{} true "Input Field"
// @Router /login-otp-verify [post]
// @Success 200 "Successfully Logged In"
// @Failure 400 "Invalid Otp"
// @Failure 500 "Faild to generate JWT"
func (u *UserHandler) LoginOtpVerify(ctx *gin.Context) {

	var body req.OTPVerifyStruct
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "Enter values Properly",
			"error":      err.Error(),
		})
		return
	}

	var user domain.User
	copier.Copy(&user, &body)

	// get the user using loginOtp useCase
	user, err := u.userUseCase.LoginOtp(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "Can't login",
			"error":      err.Error(),
		})
		return
	}

	// then varify the otp
	err = varify.TwilioVerifyOTP("+91"+user.Phone, body.OTP)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "Can't login",
			"error":      "invalid OTP",
		})
		return
	}

	// if everyting ok then generate token
	tokenString, err := auth.GenerateJWT(user.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"StatusCode": 500,
			"msg":        "can't login",
			"error":      "faild to generat JWT",
		})
		return
	}

	ctx.SetCookie("user-auth", tokenString["accessToken"], 10*60, "", "", false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"Status":     "Successfully Loged In",
	})
}

// Logout godoc
// @summary api for user to lgout
// @description user can logout
// @security ApiKeyAuth
// @id Logout
// @tags logout
// @produce json
// @Router /logout [post]
// @Success 200 "Successfully logout"
func (u *UserHandler) Logout(ctx *gin.Context) {
	ctx.SetCookie("user-auth", "", -1, "", "", false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"Status":     "Successfully Loged Out",
	})
}

// AddToCart godoc
// @summary api for add productItem to user cart
// @description user can add a stock in product to user cart
// @security ApiKeyAuth
// @id AddToCart
// @tags cart
// @produce json
// @Param input body req.ReqCart true "Input Field"
// @Router /cart [post]
// @Success 200 "Successfully productItem added to cart"
// @Failure 400 "can't add the product item into cart"
func (u *UserHandler) AddToCart(ctx *gin.Context) {

	var body req.ReqCart
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "can't bind the json",
			"error":      err.Error(),
		})
		return
	}

	// get userId and add to body
	body.UserID = helper.GetUserIdFromContext(ctx)

	_, err := u.userUseCase.SaveToCart(ctx, body)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "can't add the product item into cart",
			"error":      err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"msg":        "Successfully productItem added to cart",
		"product_id": body.ProductItemID,
	})
}

// RemoveFromCart godoc
// @summary api for remove a product from cart
// @description user can remove a signle productItem full quantity from cart
// @security ApiKeyAuth
// @id RemoveFromCart
// @tags cart
// @produce json
// @Param input body req.ReqCart{} true "Input Field"
// @Router /cart [delete]
// @Success 200 "Successfully productItem removed from cart"
// @Failure 400  "can't remove product item into cart"
func (u UserHandler) RemoveFromCart(ctx *gin.Context) {

	var body req.ReqCart
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "can't add the product into cart",
			"error":      err.Error(),
		})
		return
	}

	body.UserID = helper.GetUserIdFromContext(ctx)

	_, err := u.userUseCase.RemoveCartItem(ctx, body)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "can't remove product item into cart",
			"error":      err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"msg":        "Successfully productItem removed from cart",
		"product_id": body.ProductItemID,
	})

}

// UpdateCart godoc
// @summary api for updte productItem count
// @description user can inrement or drement count of a productItem in cart (min=1)
// @security ApiKeyAuth
// @id UpdateCart
// @tags cart
// @produce json
// @Param input body req.ReqCartCount{} true "Input Field"
// @Router /cart [put]
// @Success 200 "Successfully productItem count change on cart"
// @Failure 400  "can't change count of product item on cart"
func (u *UserHandler) UpdateCart(ctx *gin.Context) {

	var body req.ReqCartCount

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "can't bind the json",
			"error":      err.Error(),
		})
		return
	}

	body.UserID = helper.GetUserIdFromContext(ctx)

	cartItem, err := u.userUseCase.UpdateCartItem(ctx, body)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "cant update the count of productItem in cart",
			"error":      err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"msg":        "Successfully ProductItem count changed in cart",
		"product_id": cartItem.ProductItemID,
	})
}

// UserCart godoc
// @summary api for get all cart item of user
// @description user can see all productItem that stored in cart
// @security ApiKeyAuth
// @id UserCart
// @tags cart
// @produce json
// @Router /cart [get]
// @Success 200 "there is no productItems in the cart"
// @Success 200 {object} res.ResponseCart{} "there is no productItems in the cart"
// @Failure 500 "Faild to get user cart"
func (u *UserHandler) UserCart(ctx *gin.Context) {

	userId := helper.GetUserIdFromContext(ctx)

	resCart, err := u.userUseCase.GetCartItems(ctx, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"StatusCode": 500,
			"msg":        "Faild to get user cart",
			"error":      err.Error(),
		})
		return
	}

	if resCart.CartItems == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"StatusCode": 200,
			"msg":        "there is no productItems in the cart",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"msg":        "User Cart",
		"cart":       resCart,
	})
}

// ***** for user profiler ***** //

// AddAddress godoc
// @summary api for adding a new address for user
// @description get a new address from user to store the the database
// @security ApiKeyAuth
// @id AddAddress
// @tags address
// @produce json
// @Param inputs body req.ReqAddress{} true "Input Field"
// @Router /profile/address [post]
// @Success 200 "Successfully address added"
// @Failure 400 "can't add the user addres"
func (u *UserHandler) AddAddress(ctx *gin.Context) {
	var body req.ReqAddress
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "can't bind the json",
			"error":      err.Error(),
		})
		return
	}
	userID := helper.GetUserIdFromContext(ctx)

	var address domain.Address

	copier.Copy(&address, &body)

	address, err := u.userUseCase.SaveAddress(ctx, address, userID, *body.IsDefault)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "can't add the user address",
			"error":      err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"msg":        "Successfully addressed added",
	})
}

// GetAddreses godoc
// @summary api for get all address of user
// @description user can show all adderss
// @security ApiKeyAuth
// @id GetAddresses
// @tags address
// @produce json
// @Router /profile/address [get]
// @Success 200 "Successfully address got"
// @Failure 500 "Faild to get address of user"
func (u *UserHandler) GetAddresses(ctx *gin.Context) {

	userID := helper.GetUserIdFromContext(ctx)

	address, err := u.userUseCase.GetAddresses(ctx, userID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"StatusCode": 500,
			"msg":        "can't show addresses of user",
			"error":      err.Error(),
		})
		return
	}

	if address == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"StatusCode": 200,
			"msg":        "There is no address available to show",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"msg":        "Successfully addresses got",
		"addresses":  address,
	})
}

// EditAddress godoc
// @summary api for edit user address
// @description user can change existing address
// @security ApiKeyAuth
// @id EditAddress
// @tags address
// @produce json
// @Param input body req.ReqEditAddress true "Input Field"
// @Router /profile/address [put]
// @Success 200 "Successfully addresses updated"
// @Failure 400 "can't update the address"
func (u *UserHandler) EditAddress(ctx *gin.Context) {

	var body req.ReqEditAddress

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "can't bind the json",
			"error":      err.Error(),
		})
		return
	}
	userID := helper.GetUserIdFromContext(ctx)

	if err := u.userUseCase.EditAddress(ctx, body, userID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "can't update the address",
			"error":      err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"msg":        "Successfully addresses updated",
	})

}

func (u *UserHandler) DeleteAddress(ctx *gin.Context) {

}

// ** wishList **

// AddToWishList godoc
// @summary api to add a productItem to wish list
// @descritpion user can add productItem to wish list
// @security ApiKeyAuth
// @id AddToWishList
// @tags wishlist
// @produce json
// @Param product_id body int true "product_id"
// @Router /wishlist [post]
// @Success 200 "Successfully product_item added to wishlist"
// @Failure 400 "Faild to add product_item to wishlist"
func (u *UserHandler) AddToWishList(ctx *gin.Context) {
	// get productItemID using parmas
	productItemID, err := helper.StringToUint(ctx.Param("id"))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "invalid input",
		})
		return
	}

	userID := helper.GetUserIdFromContext(ctx)

	var wishList = domain.WishList{
		ProductItemID: productItemID,
		UserID:        userID,
	}
	fmt.Println(wishList.UserID, wishList.ProductItemID)

	// add to wishlist
	if err := u.userUseCase.AddToWishList(ctx, wishList); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "can't add product_item to wishlist",
			"error":      err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"msg":        "Successfully product_item added to wishlist",
	})
}

// RemoveFromWishList godoc
// @summary api to remove a productItem from wish list
// @descritpion user can remove a productItem from wish list
// @security ApiKeyAuth
// @id RemoveFromWishList
// @tags wishlist
// @produce json
// @Params product_item_id path int true "product_item_id"
// @Router /wishlist [post]
// @Success 200 "Successfully product_item remvoed from wishlist"
// @Failure 400 "Faild to remove product_item from wishlist"
func (u *UserHandler) RemoveFromWishList(ctx *gin.Context) {

	// get productItemID using parmas
	productItemID, err := helper.StringToUint(ctx.Param("id"))
	fmt.Println(productItemID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "invalid input",
		})
		return
	}

	userID := helper.GetUserIdFromContext(ctx)

	var wishList = domain.WishList{
		ProductItemID: productItemID,
		UserID:        userID,
	}

	// remove form wishlist
	if err := u.userUseCase.RemoveFromWishList(ctx, wishList); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "can't remove product_item from wishlist",
			"error":      err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"msg":        "Successfully product_item renoved from wishlist",
	})
}

// GetWishListI godoc
// @summary api get all wish list items of user
// @descritpion user get all wish list items
// @security ApiKeyAuth
// @id GetWishListI
// @tags wishlist
// @produce json
// @Router /wishlist [get]
// @Success 200 "Successfully wish list items got"
// @Success 200 "Wish list is empty"
// @Failure 400  "faild to get user wish list items"
func (u *UserHandler) GetWishListI(ctx *gin.Context) {

	userID := helper.GetUserIdFromContext(ctx)
	wishlists, err := u.userUseCase.GetWishListItems(ctx, userID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"StatusCode": 500,
			"msg":        "faild to get user wish list items",
			"error":      err.Error(),
		})
		return
	}

	if wishlists == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"StatusCode": 200,
			"msg":        "Wish list is empty",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"msg":        "Successfully wish list items got",
		"wishLists":  wishlists,
	})
}
