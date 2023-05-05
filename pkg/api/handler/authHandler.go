package handler

import (
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

	accessTokenTimeDuration := time.Minute * 15
	accessToken, err := c.authUseCase.GenerateAccessToken(ctx, userID, token.TokenForUser, accessTokenTimeDuration)
	if err != nil {
		respnonse := res.ErrorResponse(500, "faild to create access token", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, respnonse)
		return
	}

	refreshTokenTimeDuration := time.Hour * 24 * 7
	refreshToken, err := c.authUseCase.GenerateRefreshToken(ctx, userID, "user", refreshTokenTimeDuration)
	if err != nil {
		respnonse := res.ErrorResponse(500, "faild to create refresh token", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, respnonse)
		return
	}

	// authorizationValue := authorizationType + " " + accessToken
	// ctx.Header(authorizationHeaderKey, authorizationValue)

	//ctx.Header("access_token", accessToken)
	//ctx.Header("refresh_token", refreshToken)
	ctx.SetCookie("user-auth", accessToken, 15*60, "", "", false, true)

	response := res.SuccessResponse(200, "successfully logged in", res.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
	ctx.JSON(http.StatusOK, response)
}
