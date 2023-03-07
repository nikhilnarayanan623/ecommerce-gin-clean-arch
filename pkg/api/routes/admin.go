package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/middleware"
)

func AdminRoutes(router *gin.Engine, admin *handler.AdminHandler) {

	router.POST("admin/login", admin.Login)

	api := router.Group("/admin", middleware.Authentication)

	api.GET("/alluser", admin.Allusers)
	api.GET("/add-product", admin.AddCategoryGET)
	api.POST("/add-product", admin.AddCategoryPOST)
}
