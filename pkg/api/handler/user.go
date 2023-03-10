package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/auth"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper"
	service "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/varify"
)

type UserHandler struct {
	userUseCase service.UserUseCase
}

func (u *UserHandler) LoginGet(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"msg":        "detail to enter",
		"user":       helper.LoginStruct{},
	})
}

func (u *UserHandler) LoginPost(ctx *gin.Context) {

	var body helper.LoginStruct
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 404,
			"msg":        "Cant't bind the json",
			"error":      err,
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
	var user domain.Users
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

	var user domain.Users
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
		"msg":        "OTP Successfully sented to your number",
		"userId":     user.ID,
	})

}
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

	var user domain.Users
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

func (u *UserHandler) SignUpGet(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"msg":        "enter detail for signup",
		"user":       domain.Users{},
	})
}
func (u *UserHandler) SignUpPost(ctx *gin.Context) {
	var user domain.Users
	if ctx.ShouldBindJSON(&user) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "Cant't Bind The Values",
			"user":       user,
		})

		return
	}

	user, err := u.userUseCase.Signup(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 500,
			"msg":        "Invalid Inputs",
			"error":      err,
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

	products, err := u.userUseCase.ShowAllProducts(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"StatusCode": 500,
			"error":      err,
		})
		return
	}

	//find user
	response, err := u.userUseCase.Home(ctx, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"StatusCode": 500,
			"error":      err,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"msg":        "Welcome Home",
		"user":       response,
		"Products":   products,
	})
}

func (u *UserHandler) UserCart(ctx *gin.Context) {

	userId := helper.GetUserIdFromContext(ctx)

	resCart, err := u.userUseCase.GetCartItems(ctx, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"StatusCode": 500,
			"msg":        "Faild to get user cart",
			"error":      err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"msg":        "User Cart",
		"cart":       resCart,
	})
}

func (u *UserHandler) Logout(ctx *gin.Context) {
	ctx.SetCookie("user-auth", "", -1, "", "", false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"Status":     "Successfully Loged Out",
	})
}
