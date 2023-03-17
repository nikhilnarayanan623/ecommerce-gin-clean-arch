package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/auth"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper"
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

	var response helper.UserRespStrcut

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

	var response helper.UserRespStrcut
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
// @Success 200 {object} helper.LoginStruct "OK"
func (u *UserHandler) LoginGet(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"msg":        "detail to enter",
		"user":       helper.LoginStruct{},
	})
}

// LoginPost godoc
// @summary api for user login
// @description Enter user_name/phone/email with password
// @security ApiKeyAuth
// @tags login
// @id LoginPost
// @produce json
// @Param        inputs   body     helper.LoginStruct{}   true  "Input Field"
// @Router /login [post]
// @Success 200 "Successfully Loged In"
// @Failure 400 "faild to login"
// @Failure 500 "faild to generat JWT"
func (u *UserHandler) LoginPost(ctx *gin.Context) {

	var body helper.LoginStruct
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
// @tags login-otp
// @produce json
// @Param inputs body helper.OTPLoginStruct true "Input Field"
// @Router /login-otp-send [post]
// @Success 200 "Successfully Otp Send to registered number"
// @Failure 400 "Enter input properly"
// @Failure 500 "Faild to send otp"
func (u *UserHandler) LoginOtpSend(ctx *gin.Context) {

	var body helper.OTPLoginStruct
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
// @tags login-otp
// @produce json
// @param inputs body helper.OTPVerifyStruct true "Input Field"
// @Router /login-otp-verify [post]
// @Success 200 "Successfully Logged In"
// @Failure 400 "Invalid Otp"
// @Failure 500 "Faild to generate JWT"
func (u *UserHandler) LoginOtpVerify(ctx *gin.Context) {

	var body helper.OTPVerifyStruct
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

func (u *UserHandler) AddToCart(ctx *gin.Context) {

	var body helper.ReqCart
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

// to remove a productItem fom car
func (u UserHandler) RemoveFromCart(ctx *gin.Context) {

	var body helper.ReqCart
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

func (u *UserHandler) UpdateCart(ctx *gin.Context) {

	var body helper.ReqCartCount

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

// to show cart
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
// @Param inputs body helper.ReqAddress{} true "Input Field"
// @Router /profile/address [post]
// @Success 200 "Successfully address added"
// @Failure 400 "can't add the user addres"
func (u *UserHandler) AddAddress(ctx *gin.Context) {
	var body helper.ReqAddress
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
	}

	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"msg":        "Successfully addresses got",
		"addresses":  address,
	})
}

func (u *UserHandler) EditAddress(ctx *gin.Context) {

}

func (u *UserHandler) DeleteAddress(ctx *gin.Context) {

}
