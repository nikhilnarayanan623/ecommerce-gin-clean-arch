package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
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
			"Error":      err,
		})
		return
	}

	// if there is no error then responce it
	ctx.JSON(200, gin.H{
		"StatusCode": 200,
		"Status":     "Successfully Loged In",
		"user":       user,
	})
}

func (u *UserHandler) SignUp(ctx *gin.Context) {
	var user domain.Users

	if ctx.ShouldBindJSON(&user) != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "Cant't Bind The Values",
		})

		fmt.Println(user.FirstName, user.LastName)
		return
	}

	user, err := u.userUseCase.SaveUser(ctx, user)

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
		"Products":   products,
	})
}

func (u *UserHandler) Logout(ctx gin.Context) {

}
