package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/auth"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/res"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	service "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
)

type AdminHandler struct {
	adminUseCase service.AdminUseCase
}

func NewAdminHandler(adminUsecase interfaces.AdminUseCase) *AdminHandler {
	return &AdminHandler{adminUseCase: adminUsecase}
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

	admin, err := a.adminUseCase.SignUp(ctx, admin)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 500,
			"msg":        "Can't signup admin",
			"error":      err.Error(),
		})
		return
	}

	var response res.ResAdminLogin

	copier.Copy(&response, &admin)

	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"msg":        "Successfully account creatd for admin",
		"admin":      response,
	})
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

// LoginPost godoc
// @summary api for admin to login
// @id LoginPost
// @tags Admin Login
// @Param input body domain.Admin{} true "inputs"
// @Router /admin/login [post]
// @Success 200 {object} res.Response{} "successfully logged in"
// @Failure 400 {object} res.Response{} "invalid input"
// @Failure 500 {object} res.Response{} "faild to generate jwt token"
func (a *AdminHandler) AdminLoginPost(ctx *gin.Context) {

	var body domain.Admin

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response := res.ErrorResponse(400, "invalid input", err.Error(), body)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	// then check all field is empty
	if body.Email == "" && body.UserName == "" {
		err := errors.New("enter email or user_name atleast")
		response := res.ErrorResponse(400, "invalid input", err.Error(), body)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	admin, err := a.adminUseCase.Login(ctx, body)

	if err != nil {
		response := res.ErrorResponse(400, "faild to login", err.Error(), admin)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	tokenString, err := auth.GenerateJWT(admin.ID)
	if err != nil {
		response := res.ErrorResponse(500, "faild to generate jwt token", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	ctx.SetCookie("admin-auth", tokenString["accessToken"], 50*60, "", "", false, true)

	response := res.SuccessResponse(200, "successfully logged in", nil)
	ctx.JSON(http.StatusOK, response)
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

	// frist check there is no user or not
	if usersResp == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"StatusCode": 200,
			"msg":        "There is no user to show",
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

	var body req.BlockStruct
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "can't bind the user id",
			"err":        err.Error(),
		})
		return
	}

	// copy into user
	var user domain.User
	copier.Copy(&user, body)

	user, err := a.adminUseCase.BlockUser(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"StatusCode": 400,
			"msg":        "can't block user",
			"error":      err.Error(),
		})
		return
	}

	var response res.UserRespStrcut
	copier.Copy(&response, &user)
	// if successfully blocked or unblock user then response 200
	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"msg":        "Successfully blocked or unblocked user",
		"user":       response,
	})
}
