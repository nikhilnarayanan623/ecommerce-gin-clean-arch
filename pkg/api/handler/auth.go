package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/token"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase"
	usecaseInterface "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/response"
)

type AuthHandler struct {
	authUseCase usecaseInterface.AuthUseCase
}

func NewAuthHandler(authUsecase usecaseInterface.AuthUseCase) interfaces.AuthHandler {
	return &AuthHandler{
		authUseCase: authUsecase,
	}
}

// UserLogin godoc
// @summary api for user to login
// @description Enter user_name | phone | email with password
// @security ApiKeyAuth
// @id UserLogin
// @tags User Authentication
// @Param        inputs   body     request.Login{}   true  "Input Fields"
// @Router /login [post]
// @Success 200 {object} response.Response{} "successfully logged in"
// @Failure 400 {object} response.Response{}  "invalid input"
// @Failure 500 {object} response.Response{}  "failed to generate JWT"
func (c *AuthHandler) UserLogin(ctx *gin.Context) {

	var body request.Login

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, body)
		return
	}

	userID, err := c.authUseCase.UserLogin(ctx, body)

	if err != nil {

		var statusCode int

		switch true {
		case errors.Is(err, usecase.ErrEmptyLoginCredentials):
			statusCode = http.StatusBadRequest
		case errors.Is(err, usecase.ErrUserNotExist):
			statusCode = http.StatusNotFound
		case errors.Is(err, usecase.ErrUserBlocked):
			statusCode = http.StatusUnauthorized
		case errors.Is(err, usecase.ErrWrongPassword):
			statusCode = http.StatusUnauthorized
		default:
			statusCode = http.StatusInternalServerError
		}

		response.ErrorResponse(ctx, statusCode, "Failed to login", err, nil)
		return
	}

	// common functionality for admin and user
	c.setupTokenAndResponse(ctx, token.User, userID)
}

// UserLoginOtpSend godoc
// @summary api for user otp login send
// @description user can enter email/user_name/phone will send an otp to user registered phone_number
// @security ApiKeyAuth
// @id UserLoginOtpSend
// @tags User Authentication
// @Param inputs body request.OTPLogin{} true "Input Field"
// @Router /login/otp-send [post]
// @Success 200 {object} response.Response{}  "Successfully otp send to user's registered number"
// @Failure 400 {object} response.Response{}  "Enter input properly"
// @Failure 500 {object} response.Response{}  "Failed to send otp"
func (u *AuthHandler) UserLoginOtpSend(ctx *gin.Context) {

	var body request.OTPLogin
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, body)
		return
	}

	//check all input field is empty
	if body.Email == "" && body.Phone == "" && body.UserName == "" {
		err := errors.New("enter at least user_name or email or phone")
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, nil)
		return
	}

	otpID, err := u.authUseCase.UserLoginOtpSend(ctx, body)

	if err != nil {
		var statusCode int

		switch true {
		case errors.Is(err, usecase.ErrEmptyLoginCredentials):
			statusCode = http.StatusBadRequest
		case errors.Is(err, usecase.ErrUserNotExist):
			statusCode = http.StatusNotFound
		case errors.Is(err, usecase.ErrUserBlocked):
			statusCode = http.StatusUnauthorized
		default:
			statusCode = http.StatusInternalServerError
		}
		response.ErrorResponse(ctx, statusCode, "Failed to send otp", err, nil)
		return
	}

	optRes := response.OTPResponse{
		OtpID: otpID,
	}
	response.SuccessResponse(ctx, http.StatusOK, "Successfully otp send to user's registered number", optRes)
}

// UserLoginOtpVerify godoc
// @summary api for user to verify user login otp
// @description enter your otp that send to your registered number
// @security ApiKeyAuth
// @id UserLoginOtpVerify
// @tags User Authentication
// @param inputs body request.OTPVerify{} true "Input Field"
// @Router /login/otp-verify [post]
// @Success 200 {object} response.Response{} "successfully logged in using otp"
// @Failure 400 {object} response.Response{} "invalid login_otp"
// @Failure 500 {object} response.Response{} "failed to generate JWT"
func (c *AuthHandler) UserLoginOtpVerify(ctx *gin.Context) {

	var body request.OTPVerify
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, body)
		return
	}

	// get the user using loginOtp useCase
	userID, err := c.authUseCase.LoginOtpVerify(ctx, body)
	if err != nil {
		var statusCode int
		switch true {
		case errors.Is(err, usecase.ErrOtpExpired):
			statusCode = http.StatusBadRequest
		case errors.Is(err, usecase.ErrInvalidOtp):
			statusCode = http.StatusUnauthorized
		default:
			statusCode = http.StatusInternalServerError
		}
		response.ErrorResponse(ctx, statusCode, "Failed to verify otp", err, nil)
		return
	}

	c.setupTokenAndResponse(ctx, token.User, userID)
}

// UserSignUp godoc
// @summary api for user to signup
// @security ApiKeyAuth
// @id UserSignUp
// @tags User Authentication
// @Param input body request.UserSignUp{} true "Input Fields"
// @Router /signup [post]
// @Success 200 {object} response.Response{} "Successfully account created for user"
// @Failure 400 {object} response.Response{} "invalid input"
func (c *AuthHandler) UserSignUp(ctx *gin.Context) {

	var body request.UserSignUp

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, body)
		return
	}

	var user domain.User
	copier.Copy(&user, body)

	err := c.authUseCase.UserSignUp(ctx, user)

	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, usecase.ErrUserAlreadyExit) {
			statusCode = http.StatusConflict
		}

		response.ErrorResponse(ctx, statusCode, "Failed to signup", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusCreated, "Successfully Account Created")
}

func (c *AuthHandler) UserRenewAccessToken() gin.HandlerFunc {
	return c.renewAccessToken(token.User)
}

// UserSignUp godoc
// @summary api for admin to login
// @security ApiKeyAuth
// @id AdminLogin
// @tags Admin Authentication
// @Param input body request.Login{} true "Input Fields"
// @Router /admin/login [post]
// @Success 200 {object} response.Response{} "Successfully account created for user"
// @Failure 400 {object} response.Response{} "invalid input"
func (c *AuthHandler) AdminLogin(ctx *gin.Context) {

	var body request.Login

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, body)
		return
	}

	adminID, err := c.authUseCase.AdminLogin(ctx, body)
	if err != nil {

		var statusCode int

		switch true {
		case errors.Is(err, usecase.ErrEmptyLoginCredentials):
			statusCode = http.StatusBadRequest
		case errors.Is(err, usecase.ErrUserNotExist):
			statusCode = http.StatusNotFound
		case errors.Is(err, usecase.ErrWrongPassword):
			statusCode = http.StatusUnauthorized
		default:
			statusCode = http.StatusInternalServerError
		}

		response.ErrorResponse(ctx, statusCode, "Failed to login", err, nil)
		return
	}

	// setup token common part
	c.setupTokenAndResponse(ctx, token.Admin, adminID)
}

func (c *AuthHandler) AdminRenewAccessToken() gin.HandlerFunc {
	return c.renewAccessToken(token.Admin)
}

func (c *AuthHandler) setupTokenAndResponse(ctx *gin.Context, tokenUser token.UserType, userID uint) {

	tokenParams := usecaseInterface.GenerateTokenParams{
		UserID:   userID,
		UserType: tokenUser,
	}

	accessToken, err := c.authUseCase.GenerateAccessToken(ctx, tokenParams)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to generate access token", err, nil)
		return
	}

	refreshToken, err := c.authUseCase.GenerateRefreshToken(ctx, usecaseInterface.GenerateTokenParams{
		UserID:   userID,
		UserType: tokenUser,
	})
	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to generate refresh token", err, nil)
		return
	}

	// authorizationValue := authorizationType + " " + accessToken
	// ctx.Header(authorizationHeaderKey, authorizationValue)

	//ctx.Header("access_token", accessToken)
	//ctx.Header("refresh_token", refreshToken)
	cookieName := "auth-" + string(tokenUser)
	ctx.SetCookie(cookieName, accessToken, 15*60, "", "", false, true)

	tokenRes := response.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	response.SuccessResponse(ctx, http.StatusOK, "Successfully logged in", tokenRes)
}

func (c *AuthHandler) renewAccessToken(tokenUser token.UserType) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var body request.RefreshToken

		if err := ctx.ShouldBindJSON(&body); err != nil {
			response.ErrorResponse(ctx, http.StatusBadRequest, BindJsonFailMessage, err, body)
			return
		}

		refreshSession, err := c.authUseCase.VerifyAndGetRefreshTokenSession(ctx, body.RefreshToken, tokenUser)

		if err != nil {
			var statusCode int

			switch true {
			case errors.Is(err, usecase.ErrInvalidRefreshToken):
				statusCode = http.StatusUnauthorized
			case errors.Is(err, usecase.ErrRefreshSessionNotExist):
				statusCode = http.StatusNotFound
			case errors.Is(err, usecase.ErrRefreshSessionExpired):
				statusCode = http.StatusUnauthorized
			case errors.Is(err, usecase.ErrRefreshSessionBlocked):
				statusCode = http.StatusForbidden
			default:
				statusCode = http.StatusInternalServerError
			}
			response.ErrorResponse(ctx, statusCode, "Failed verify refresh token", err, nil)
			return
		}

		accessTokenParams := usecaseInterface.GenerateTokenParams{
			UserID:   refreshSession.UserID,
			UserType: tokenUser,
		}

		accessToken, err := c.authUseCase.GenerateAccessToken(ctx, accessTokenParams)

		if err != nil {
			response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed generate access token", err, nil)
			return
		}
		cookieName := "auth-" + string(tokenUser)
		ctx.SetCookie(cookieName, accessToken, 15*60, "", "", false, true)

		accessTokenRes := response.TokenResponse{
			AccessToken: accessToken,
		}
		response.SuccessResponse(ctx, http.StatusOK, "Successfully access token generated using refresh token", accessTokenRes)
	}
}
