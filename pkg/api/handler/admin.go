package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/auth"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper"
	service "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
)

type AdminHandler struct {
	adminUseCase service.AdminUseCase
}

func (a *AdminHandler) SignUPGet(ctx *gin.Context) {

	// getting the validation engine and type casting it.

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
	if err := ctx.ShouldBindJSON(&admin); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "Can't bind the values invalid inputs",
			"error":      err.Error(),
		})
		return
	}

	// then check all field is empty
	if admin.Email == "" && admin.UserName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "Enter atleast user_name or email",
		})
		return
	}

	admin, err := a.adminUseCase.Login(ctx, admin)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "Can't login",
			"err":        err.Error(),
		})
		return
	}

	tokenString, err := auth.GenerateJWT(admin.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"StatusCode": 500,
			"msg":        "Error to Create JWT",
		})
	}

	// if no error then copy the admin details to response
	var response helper.ResAdminLogin
	copier.Copy(&response, &admin)

	ctx.SetCookie("admin-auth", tokenString["accessToken"], 10*60, "", "", false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"msg":        "Successfully loged in",
		"admin":      response,
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
	// if successfully blocked or unblock user then response 200
	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"msg":        "Successfully bocked or unblocked user",
		"user":       user,
	})
}

func (a *AdminHandler) CategoryGET(ctx *gin.Context) {

	categories, err := a.adminUseCase.GetCategory(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"StatusCode": 500,
			"msg":        "Faild to get categories",
			"error":      err,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"StatsuCode": 200,
		"msg":        "Category Page",
		"categories": categories,
	})
}
func (a *AdminHandler) CategoryPOST(ctx *gin.Context) {

	var productCategory domain.Category

	if ctx.ShouldBindJSON(&productCategory) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatsuCode": 400,
			"Error":      "Error to bind the input",
		})
		return
	}

	respose, err := a.adminUseCase.AddCategory(ctx, productCategory)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatsuCode": 400,
			"msg":        "category can't be add",
			"err":        err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"StatsuCode": 200,
		"msg":        "category added",
		"categoty":   respose,
	})
}

func (a *AdminHandler) ShowAllProducts(ctx *gin.Context) {

}

func (a *AdminHandler) AddProducts(ctx *gin.Context) {

	var body helper.ProductRequest
	if ctx.ShouldBindJSON(&body) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatsuCode": 400,
			"Error":      "Error to bind the input",
		})
		return
	}

	product, err := a.adminUseCase.AddProducts(ctx, body)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatsuCode": 400,
			"msg":        "product can't be add",
			"err":        err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"StatsuCode": 200,
		"msg":        "product added",
		"product":    product,
	})
}
