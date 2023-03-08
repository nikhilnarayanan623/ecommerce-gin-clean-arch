package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/auth"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper"
	service "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
)

type UserHandler struct {
	userUseCase service.UserUseCase
}

func (u *UserHandler) LoginGet(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"msg":        "detail to enter",
		"imp":        "no user_name ",
		"user":       helper.LoginStruct{},
	})
}

func (u *UserHandler) LoginPost(ctx *gin.Context) {

	var body helper.LoginStruct

	if ctx.ShouldBindJSON(&body) != nil {

		ctx.JSON(404, gin.H{
			"StatusCode": 400,
			"msg":        "Enter values Properly",
			"error":      "Cant't bind the json",
		})
		return
	}

	responce, err := u.userUseCase.Login(ctx, body)

	if err != nil {

		ctx.JSON(400, gin.H{
			"StatusCode": 400,
			"error":      err,
		})
		return
	}

	// //create a new token
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
	// 	ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
	// 	Id:        fmt.Sprint(responce.ID),
	// })

	// //sign the token
	// signedString, err := token.SignedString([]byte(config.GetJWTCofig()))

	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{
	// 		"StatusCode": 500,
	// 		"msg":        "Error to Create JWT",
	// 	})
	// }

	tokenString, err := auth.GenerateJWT(responce.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"StatusCode": 500,
			"msg":        "Error to Create JWT",
		})
	}

	ctx.SetCookie("user-auth", tokenString["accessToken"], 10*60, "", "", false, true)

	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"Status":     "Successfully Loged In",
		"user":       responce,
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

func (u *UserHandler) UserCart(ctx *gin.Context) {

	userIdStr := ctx.GetString("userId")
	userIdInt, _ := strconv.Atoi(userIdStr)
	userId := uint(userIdInt)

	fmt.Println(userId)

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

func (u *UserHandler) Logout(ctx gin.Context) {

}
