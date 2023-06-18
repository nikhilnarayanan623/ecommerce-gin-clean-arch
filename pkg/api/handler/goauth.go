package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/config"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/token"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/response"
)

func (c *AuthHandler) UserGoogleAuthLoginPage(ctx *gin.Context) {

	ctx.HTML(200, "goauth.html", nil)
}

func (c *AuthHandler) UserGoogleAuthInitialize(ctx *gin.Context) {

	// setup the google provider
	goauthClientID := config.GetConfig().GoathClientID
	gouthClientSecret := config.GetConfig().GoauthClientSecret
	callbackUrl := config.GetConfig().GoauthCallbackUrl

	// setup privier
	goth.UseProviders(
		google.New(goauthClientID, gouthClientSecret, callbackUrl, "email", "profile"),
	)

	// start the google login
	gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
}

func (c *AuthHandler) UserGoogleAuthCallBack(ctx *gin.Context) {

	googleUser, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get user details from google", err, nil)
		return
	}

	var user domain.User

	copier.Copy(&user, &googleUser)

	userID, err := c.authUseCase.GoogleLogin(ctx, user)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to login with google", err, nil)
		return
	}

	c.setupTokenAndResponse(ctx, token.User, userID)
}
