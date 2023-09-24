package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/response"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/service/token"
)

// UserGoogleAuthLoginPage godoc
//
//	@Summary		To load google login page (User)
//	@Description	API for user to load google login page
//	@Id				UserGoogleAuthLoginPage
//	@Tags			User Authentication
//	@Router			/auth/google-auth [get]
//	@Success		200	{object}	response.Response{}	"Successfully google login page loaded"
func (c *AuthHandler) UserGoogleAuthLoginPage(ctx *gin.Context) {

	ctx.HTML(200, "goauth.html", nil)
}

// UserGoogleAuthInitialize godoc
//
//	@Summary		Initialize google auth (User)
//	@Description	API for user to initialize google auth
//	@Id				UserGoogleAuthInitialize
//	@Tags			User Authentication
//	@Router			/auth/google-auth/initialize [get]
func (c *AuthHandler) UserGoogleAuthInitialize(ctx *gin.Context) {

	// setup the google provider
	goauthClientID := c.config.GoathClientID
	gouthClientSecret := c.config.GoauthClientSecret
	callbackUrl := c.config.GoauthCallbackUrl

	// setup privier
	goth.UseProviders(
		google.New(goauthClientID, gouthClientSecret, callbackUrl, "email", "profile"),
	)

	// start the google login
	gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
}

// UserGoogleAuthCallBack godoc
//
//	@Summary		Google auth callback (User)
//	@Description	API for google to callback after authentication
//	@Id				UserGoogleAuthCallBack
//	@Tags			User Authentication
//	@Router			/auth/google-auth/callback [post]
//	@Success		200	{object}	response.Response{}	"Successfully logged in with google"
//	@Failure		500	{object}	response.Response{}	"Failed Login with google"
func (c *AuthHandler) UserGoogleAuthCallBack(ctx *gin.Context) {

	googleUser, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get user details from google", err, nil)
		return
	}

	user := domain.User{
		FirstName:   googleUser.FirstName,
		LastName:    googleUser.LastName,
		Email:       googleUser.Email,
		GoogleImage: googleUser.AvatarURL,
	}

	userID, err := c.authUseCase.GoogleLogin(ctx, user)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to login with google", err, nil)
		return
	}

	c.setupTokenAndResponse(ctx, token.User, userID)
}
