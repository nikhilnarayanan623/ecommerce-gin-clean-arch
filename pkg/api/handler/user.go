package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/jinzhu/copier"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/config"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper"
	service "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
)

type UserHandler struct {
	userUseCase service.UserUseCase
}

func (u *UserHandler) Login(ctx *gin.Context) {

	var user domain.Users

	if ctx.ShouldBindJSON(&user) != nil {

		ctx.JSON(404, gin.H{
			"StatusCode": 400,
			"msg":        "Enter values Properly",
			"error":      "Cant't bind the json",
		})
		return
	}

	user, err := u.userUseCase.Login(ctx, user)

	if err != nil {

		ctx.JSON(400, gin.H{
			"StatusCode": 400,
			"error":      err,
		})
		return
	}
	//create a new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
	})

	//sign the token
	signedString, err := token.SignedString([]byte(config.GetJWTCofig()))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"StatusCode": 500,
			"msg":        "Error to Create JWT",
		})
	}

	ctx.SetCookie("jwt-auth", signedString, 10*60, "", "", false, true)

	// everything ok then responce 200 with user details
	var response helper.UserRespStrcut
	copier.Copy(&response, &user) // copy required data only

	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"Status":     "Successfully Loged In",
		"user":       user,
	})
}

func (u *UserHandler) SignUp(ctx *gin.Context) {
	var user domain.Users

	if ctx.BindJSON(&user) != nil {

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
			"StatusCode": 400,
			"msg":        "Invalid Inputs",
			"error":      err,
		})
		return
	}

	ctx.JSON(200, gin.H{
		"StatusCode": 200,
		"msg":        "Successfully Account Created",
		"user":       user,
	})
}

func (u *UserHandler) Home(ctx *gin.Context) {

	products, err := u.userUseCase.ShowAllProducts(ctx)

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
		"Products":   products,
	})
}

func (u *UserHandler) Logout(ctx gin.Context) {

}
