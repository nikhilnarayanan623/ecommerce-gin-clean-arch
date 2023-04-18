package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/auth"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/config"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
)

func (c *UserHandler) GoauthLoginPage(ctx *gin.Context) {

	ctx.HTML(200, "goauth.html", nil)
}

func (c *UserHandler) IntitializeGoogleAuth(ctx *gin.Context) {

	// setup the google provider
	goauthClientID := config.GetCofig().GoathClientID
	gouthClientSecret := config.GetCofig().GoauthClientSecret
	var callbackUrl string
	// for using both local host and server side check path and  set different url
	if ctx.Request.Host == "localhost:8000" {
		callbackUrl = "http://localhost:8000/login/auth/google/callback" // copy from redirect url from google auth client creation

	} else { // if this call from server side then change the callback url
		callbackUrl = config.GetCofig().GoauthCallbackUrl
	}

	// setup privider
	goth.UseProviders(
		google.New(goauthClientID, gouthClientSecret, callbackUrl, "email", "profile"),
	)

	// start the google login
	gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
}

func (c *UserHandler) CallbackAuth(ctx *gin.Context) {

	googleUser, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)

	if err != nil {
		response := res.ErrorResponse(500, "faild to get user details from google", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	var user domain.User

	copier.Copy(&user, &googleUser)
	fmt.Println("copied google user", user)

	user, err = c.userUseCase.GoogleLogin(ctx, user)

	if err != nil {
		response := res.ErrorResponse(500, "faild to login with google", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
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
	response := res.SuccessResponse(200, "successfully logged using google auth", tokenString["accessToken"])
	ctx.JSON(http.StatusOK, response)
}
