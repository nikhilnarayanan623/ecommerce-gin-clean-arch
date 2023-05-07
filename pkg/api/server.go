package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/cmd/api/docs"
	handlerInterface "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/middleware"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/routes"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type ServerHTTP struct {
	Engine *gin.Engine
}

func NewServerHTTP(authHandler handlerInterface.AuthHandler, middleware middleware.Middleware,
	adminHandler handlerInterface.AdminHandler, userHandler handlerInterface.UserHandler,
	productHandler handlerInterface.ProductHandler, orderHandler handlerInterface.OrderHandler,
	couponHandler handlerInterface.CouponHandler) *ServerHTTP {

	engine := gin.New()

	engine.LoadHTMLGlob("views/*.html")

	engine.Use(gin.Logger())

	// swagger docs
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// set up routes
	routes.UserRoutes(engine.Group("/"), authHandler, middleware, userHandler, productHandler, orderHandler, couponHandler)
	routes.AdminRoutes(engine.Group("/admin"), authHandler, middleware, adminHandler, productHandler, orderHandler, couponHandler)

	// no handler
	engine.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"StatusCode": 404,
			"msg":        "invalid url",
		})
	})

	return &ServerHTTP{Engine: engine}
}

func (s *ServerHTTP) Start() {

	s.Engine.Run(":8000")
}
