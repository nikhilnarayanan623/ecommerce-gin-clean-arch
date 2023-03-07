package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	service "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
)

type AdminHandler struct {
	adminUseCase service.AdminUseCase
}

func (a *AdminHandler) Login(ctx *gin.Context) {

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
	}

	ctx.JSON(http.StatusOK, gin.H{
		"StatusCode": 200,
		"msg":        "Successfully loged in",
		"admin":      admin,
	})
}



func (a *AdminHandler) Home(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Welcome to Admin Panel",
	})
}
func (a *AdminHandler) Allusers(ctx *gin.Context) {

	users, err := a.adminUseCase.FindAllUser(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"Error ": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"users": users,
	})

}
func (a *AdminHandler) BlockUser(ctx *gin.Context) {

}

func (a *AdminHandler) AddCategoryGET(ctx *gin.Context) {

	ctx.JSON(200, gin.H{
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

		ctx.JSON(400, gin.H{
			"Error": "Error to bind the input",
		})
		return
	}

	category, err := a.adminUseCase.AddCategory(ctx, productCategory)
	if err != nil {
		ctx.JSON(500, gin.H{
			"msg": "category can't add",
			"err": err,
		})
		return
	}

	ctx.JSON(200, gin.H{
		"msg":      "category added",
		"categoty": category,
	})
}
