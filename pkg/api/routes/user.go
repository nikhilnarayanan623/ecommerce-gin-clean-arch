package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/middleware"
)

func UserRoutes(router *gin.Engine, user *handler.UserHandler) {

	router.GET("/login", user.LoginGet)
	router.POST("/login", user.LoginPost)

	router.GET("/signup", user.SignUpGet)
	router.POST("/signup", user.SignUpPost)

	api := router.Group("/", middleware.AuthenticateUser)

	api.GET("/", user.Home)
	api.GET("/cart", user.UserCart)
}
