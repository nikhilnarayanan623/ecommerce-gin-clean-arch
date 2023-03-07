package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/middleware"
)

func UserRoutes(router *gin.Engine, user *handler.UserHandler) {

	router.POST("/login", user.Login)
	router.POST("/signup", user.SignUp)

	api := router.Group("/", middleware.Authentication)

	api.GET("/", user.Home)
}
