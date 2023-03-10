package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/middleware"
)

func AdminRoutes(router *gin.Engine, admin *handler.AdminHandler) {

	router.GET("/admin/login", admin.LoginGet)
	router.POST("admin/login", admin.LoginPost)

	router.GET("admin/signup", admin.SignUPGet)
	router.POST("admin/signup", admin.SignUpPost)

	api := router.Group("/admin", middleware.AuthenticateAdmin)

	api.GET("/alluser", admin.Allusers)
	api.POST("/block-user", admin.BlockUser)
	api.GET("/category", admin.CategoryGET)
	api.POST("/category", admin.CategoryPOST)
	api.POST("/product", admin.AddProducts)

}
