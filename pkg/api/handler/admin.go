package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/auth"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper"
	service "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
)

type AdminHandler struct {
	adminUseCase service.AdminUseCase
}

func (a *AdminHandler) SignUPGet(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"msg":        "admin signup details",
		"user_name":  "string(user name of admin)",
		"email":      "string(admin email)",
		"password":   "string(enter a strong password)",
	})
}

func (a *AdminHandler) SignUpPost(ctx *gin.Context) {

	var admin domain.Admin

	if ctx.ShouldBindJSON(&admin) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 500,
			"msg":        "Can't signup admin",
			"error":      "Invalid input can't bind JSON",
		})
		return
	}

	dbAdmin, err := a.adminUseCase.SignUp(ctx, admin)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 500,
			"msg":        "Can't signup admin",
			"error":      err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"msg":        "Successfully account creatd for admin",
		"admin":      dbAdmin,
	})

	//ctx.Redirect(http.StatusSeeOther, "/admin/login")
}

func (a *AdminHandler) LoginGet(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"msg":        "admin login details",
		"email":      "string(enter email)",
		"user_name":  "string(enter user name)",
		"password":   "string(enter password)",
	})
}

func (a *AdminHandler) LoginPost(ctx *gin.Context) {

	var admin domain.Admin

	if ctx.ShouldBindJSON(&admin) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "Can't bind the values invalid inputs",
		})
		return
	}

	admin, err := a.adminUseCase.Login(ctx, admin)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "Can't login",
			"err":        err,
		})
		return
	}

	// //create a new token
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
	// 	ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
	// })

	// //sign the token
	// signedString, err := token.SignedString([]byte(config.GetJWTCofig()))

	tokenString, err := auth.GenerateJWT(admin.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"StatusCode": 500,
			"msg":        "Error to Create JWT",
		})
	}

	ctx.SetCookie("admin-auth", tokenString["accessToken"], 10*60, "", "", false, true)

	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"msg":        "Successfully loged in",
		"admin":      admin,
	})
}

func (a *AdminHandler) Home(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"message":    "Welcome to Admin Panel",
	})
}
func (a *AdminHandler) Allusers(ctx *gin.Context) {

	usersResp, err := a.adminUseCase.FindAllUser(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"StatusCode": 500,
			"Error ":     err.Error(),
		})
		return
	}

	// if no error then response stats code 200 with usres
	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"users":      usersResp,
	})

}
func (a *AdminHandler) BlockUser(ctx *gin.Context) {
	var body helper.BlockStruct

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "invalid input tag",
			"err":        err,
		})
		return
	}

	user, err := a.adminUseCase.BlockUser(ctx, body)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 500,
			"error":      err,
		})
		return
	}
	fmt.Println("last")
	// if successfully blocked or unblock user then response 200
	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"msg":        "Successfully bocked or unblocked user",
		"user":       user,
	})
}

func (a *AdminHandler) AddCategoryGET(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, gin.H{
		"StatsuCode":    200,
		"msg":           "Add Product Page",
		"category_id":   "int(if you providing a sub category)",
		"categroy_name": "string(name of the category)",
	})
}
func (a *AdminHandler) AddCategoryPOST(ctx *gin.Context) {
	fmt.Println("here")

	var productCategory domain.Category

	if ctx.ShouldBindJSON(&productCategory) != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatsuCode": 400,
			"Error":      "Error to bind the input",
		})
		return
	}

	category, err := a.adminUseCase.AddCategory(ctx, productCategory)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatsuCode": 400,
			"msg":        "category can't add",
			"err":        err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"StatsuCode": 200,
		"msg":        "category added",
		"categoty":   category,
	})
}
