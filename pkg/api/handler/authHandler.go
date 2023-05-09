package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	handlerInterface "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/token"
	usecaseInterface "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
)

type AuthHandler struct {
	authUseCase usecaseInterface.AuthUseCase
}

func NewAuthHandler(authUsecase usecaseInterface.AuthUseCase) handlerInterface.AuthHandler {
	return &AuthHandler{
		authUseCase: authUsecase,
	}
}

const (
	authorizationType      = "bearer"
	authorizationHeaderKey = "authorization"
)

// UserLogin godoc
// @summary api for user to login
// @description Enter user_name | phone | email with password
// @security ApiKeyAuth
// @tags User Login
// @id UserLogin
// @Param        inputs   body     req.Login{}   true  "Input Field"
// @Router /login [post]
// @Success 200 {object} res.Response{} "successfully logged in"
// @Failure 400 {object} res.Response{}  "invalid input"
// @Failure 500 {object} res.Response{}  "faild to generat JWT"
func (c *AuthHandler) UserLogin(ctx *gin.Context) {
	var body req.Login

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		response := res.ErrorResponse(400, "faild to bind json input", err.Error(), body)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	userID, err := c.authUseCase.UserLogin(ctx, body)
	if err != nil {
		response := res.ErrorResponse(400, "faild to login", err.Error(), nil)
		ctx.JSON(400, response)
		return
	}

	accessTokenExpireDate := time.Now().Add(time.Minute * 15)
	accessToken, err := c.authUseCase.GenerateAccessToken(ctx, usecaseInterface.GenerateTokenParams{
		UserID:     userID,
		UserType:   token.TokenForUser,
		ExpireDate: accessTokenExpireDate,
	})
	if err != nil {
		respnonse := res.ErrorResponse(500, "faild to create access token", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, respnonse)
		return
	}

	refreshTokenExpireDate := time.Now().Add(time.Hour * 24 * 7)
	refreshToken, err := c.authUseCase.GenerateRefreshToken(ctx, usecaseInterface.GenerateTokenParams{
		UserID:     userID,
		UserType:   token.TokenForUser,
		ExpireDate: refreshTokenExpireDate,
	})
	if err != nil {
		respnonse := res.ErrorResponse(500, "faild to create refresh token", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, respnonse)
		return
	}

	// authorizationValue := authorizationType + " " + accessToken
	// ctx.Header(authorizationHeaderKey, authorizationValue)

	//ctx.Header("access_token", accessToken)
	//ctx.Header("refresh_token", refreshToken)
	cookieName := "auth-" + string(token.TokenForUser)
	ctx.SetCookie(cookieName, accessToken, 15*60, "", "", false, true)

	response := res.SuccessResponse(200, "successfully logged in", res.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
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
func (u *AuthHandler) UserLoginOtpSend(ctx *gin.Context) {

	var body req.OTPLogin
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

	otpRes, err := u.authUseCase.UserLoginOtpSend(ctx, body)

	if err != nil {
		var response res.Response
		if errors.Is(err, errors.New("faild to send otp")) {
			response = res.ErrorResponse(500, "faild to send otp", err.Error(), nil)
		} else {
			response = res.ErrorResponse(400, "can't login", err.Error(), nil)
		}

		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := res.SuccessResponse(200, "successfully otp send to registered number", otpRes)
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
func (c *AuthHandler) UserLoginOtpVerify(ctx *gin.Context) {

	var body req.OTPVerify
	if err := ctx.ShouldBindJSON(&body); err != nil {
		response := res.ErrorResponse(400, "invalid input", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	// get the user using loginOtp useCase
	userID, err := c.authUseCase.LoginOtpVerify(ctx, body)
	if err != nil {
		response := res.ErrorResponse(400, "faild to login", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	accessTokenExpireDate := time.Now().Add(time.Minute * 15)
	accessToken, err := c.authUseCase.GenerateAccessToken(ctx, usecaseInterface.GenerateTokenParams{
		UserID:     userID,
		UserType:   token.TokenForUser,
		ExpireDate: accessTokenExpireDate,
	})
	if err != nil {
		respnonse := res.ErrorResponse(500, "faild to create access token", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, respnonse)
		return
	}

	refreshTokenExpireDate := time.Now().Add(time.Hour * 24 * 7)
	refreshToken, err := c.authUseCase.GenerateRefreshToken(ctx, usecaseInterface.GenerateTokenParams{
		UserID:     userID,
		UserType:   token.TokenForUser,
		ExpireDate: refreshTokenExpireDate,
	})
	if err != nil {
		respnonse := res.ErrorResponse(500, "faild to create refresh token", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, respnonse)
		return
	}

	// authorizationValue := authorizationType + " " + accessToken
	// ctx.Header(authorizationHeaderKey, authorizationValue)

	//ctx.Header("access_token", accessToken)
	//ctx.Header("refresh_token", refreshToken)
	cookieName := "auth-" + string(token.TokenForUser)
	ctx.SetCookie(cookieName, accessToken, 15*60, "", "", false, true)

	response := res.SuccessResponse(200, "successfully logged in", res.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
	ctx.JSON(http.StatusOK, response)
}

func (c *AuthHandler) UserRenewRefreshToken() gin.HandlerFunc {
	return c.renewAccessToken(token.TokenForUser)
}

func (c *AuthHandler) renewAccessToken(usedFor token.UserType) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var body req.RefreshToken

		err := ctx.ShouldBindJSON(&body)
		if err != nil {
			response := res.ErrorResponse(400, "faild to bind inputs", err.Error(), nil)
			ctx.JSON(http.StatusBadRequest, response)
			return
		}

		refreshSession, err := c.authUseCase.VerifyAndGetRefreshTokenSession(ctx, body.RefreshToken, usedFor)

		if err != nil {
			response := res.ErrorResponse(400, "faild to get refresh sessions", err.Error(), nil)
			ctx.JSON(http.StatusBadRequest, response)
			return
		}
		accessTokenExpireDate := time.Now().Add(time.Minute * 15)
		accessTokenParams := usecaseInterface.GenerateTokenParams{
			UserID:     refreshSession.UserID,
			UserType:   usedFor,
			ExpireDate: accessTokenExpireDate,
		}
		accessToken, err := c.authUseCase.GenerateAccessToken(ctx, accessTokenParams)

		if err != nil {
			response := res.ErrorResponse(500, "faild to generate access token", err.Error(), nil)
			ctx.JSON(http.StatusInternalServerError, response)
			return
		}
		cookieName := "auth-" + string(token.TokenForUser)
		ctx.SetCookie(cookieName, accessToken, 15*60, "", "", false, true)

		response := res.SuccessResponse(http.StatusOK, "successfylly access token generated using refresh token",
			res.TokenResponse{
				AccessToken: accessToken,
			})

		ctx.JSON(http.StatusOK, response)
	}
}
