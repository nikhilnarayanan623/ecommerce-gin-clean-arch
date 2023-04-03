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
	service "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
)

type AdminHandler struct {
	adminUseCase service.AdminUseCase
}

func NewAdminHandler(adminUsecase service.AdminUseCase) *AdminHandler {
	return &AdminHandler{adminUseCase: adminUsecase}
}

// // AdminSignupGet godoc
// // @summary api for admin to login
// // @id AdminSignupGet
// // @tags Admin Login
// // @Param input body domain.Admin{} true "inputs"
// // @Router /admin/login [post]
// // @Success 200 {object} res.Response{} "successfully logged in"
// // @Failure 400 {object} res.Response{} "invalid input"
// // @Failure 500 {object} res.Response{} "faild to generate jwt token"

// func (a *AdminHandler) AdminS(ctx *gin.Context) {

// 	var admin domain.Admin

// 	if ctx.ShouldBindJSON(&admin) != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"StatusCode": 500,
// 			"msg":        "Can't signup admin",
// 			"error":      "Invalid input can't bind JSON",
// 		})
// 		return
// 	}

// 	admin, err := a.adminUseCase.SignUp(ctx, admin)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"StatusCode": 500,
// 			"msg":        "Can't signup admin",
// 			"error":      err.Error(),
// 		})
// 		return
// 	}

// 	var response res.ResAdminLogin

// 	copier.Copy(&response, &admin)

// 	ctx.JSON(http.StatusOK, gin.H{
// 		"StatusCode": 200,
// 		"msg":        "Successfully account creatd for admin",
// 		"admin":      response,
// 	})
// }

// AdminLogin godoc
// @summary api for admin to login
// @id AdminLogin
// @tags Admin Login
// @Param input body req.LoginStruct{} true "inputs"
// @Router /admin/login [post]
// @Success 200 {object} res.Response{} "successfully logged in"
// @Failure 400 {object} res.Response{} "invalid input"
// @Failure 500 {object} res.Response{} "faild to generate jwt token"
func (a *AdminHandler) AdminLogin(ctx *gin.Context) {

	var body req.LoginStruct

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

	var admin domain.Admin
	copier.Copy(&admin, &body)
	admin, err := a.adminUseCase.Login(ctx, admin)

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

	ctx.SetCookie("admin-auth", tokenString["accessToken"], 60*60, "", "", false, true)

	response := res.SuccessResponse(200, "successfully logged in", nil)
	ctx.JSON(http.StatusOK, response)
}

// AdminHome godoc
// @summary api admin home
// @id AdminHome
// @tags Admin Home
// @Router /admin [get]
// @Success 200 {object} res.Response{} "successfully logged in"
func (a *AdminHandler) AdminHome(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"message":    "Welcome to Admin Home",
	})
}

// ListUsers godoc
// @summary api for admin to list users
// @id ListUsers
// @tags Admin User
// @Router /admin/users [get]
// @Success 200 {object} res.Response{} "successfully got all users"
// @Failure 500 {object} res.Response{} "faild to get all users"
func (a *AdminHandler) ListUsers(ctx *gin.Context) {

	users, err := a.adminUseCase.FindAllUser(ctx)
	if err != nil {
		respone := res.ErrorResponse(500, "faild to get all users", err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, respone)
		return
	}

	// check there is no usee
	if users == nil {
		response := res.SuccessResponse(200, "there is no users to show", nil)
		ctx.JSON(http.StatusOK, response)
		return
	}

	response := res.SuccessResponse(200, "successfully got all users", users)
	ctx.JSON(http.StatusOK, response)

}

// BlockUser godoc
// @summary api for admin to block or unblock user
// @id BlockUser
// @tags Admin User
// @Param input body req.BlockStruct{} true "inputs"
// @Router /admin/users/block [patch]
// @Success 200 {object} res.Response{} "Successfully changed user block_status"
// @Failure 400 {object} res.Response{} "invalid input"
func (a *AdminHandler) BlockUser(ctx *gin.Context) {

	var body req.BlockStruct

	if err := ctx.ShouldBindJSON(&body); err != nil {
		response := res.ErrorResponse(400, "invalid input", err.Error(), body)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	err := a.adminUseCase.BlockUser(ctx, body.UserID)
	if err != nil {
		response := res.ErrorResponse(400, "faild to change user block_status", err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	response := res.SuccessResponse(200, "Successfully changed user block_status", body)
	// if successfully blocked or unblock user then response 200
	ctx.JSON(http.StatusOK, response)
}
