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
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
)

func (c *AuthHandler) UserGoogleAuthLoginPage(ctx *gin.Context) {

	ctx.HTML(200, "goauth.html", nil)
}

func (c *AuthHandler) IntitializeGoogleAuth(ctx *gin.Context) {

	// setup the google provider
	goauthClientID := config.GetConfig().GoathClientID
	gouthClientSecret := config.GetConfig().GoauthClientSecret
	callbackUrl := config.GetConfig().GoauthCallbackUrl

	// setup privider
	goth.UseProviders(
		google.New(goauthClientID, gouthClientSecret, callbackUrl, "email", "profile"),
	)

	// start the google login
	gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
}

func (c *AuthHandler) CallbackAuth(ctx *gin.Context) {

	googleUser, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)

	if err != nil {
		response := res.ErrorResponse(500, "faild to get user details from google", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	var user domain.User

	copier.Copy(&user, &googleUser)

	userID, err := c.authUseCase.GoogleLogin(ctx, user)
	if err != nil {
		response := res.ErrorResponse(500, "faild to login with google", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	c.setupTokenAndReponse(ctx, token.TokenForUser, userID)

}
