package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/cmd/api/docs"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/routes"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type ServerHTTP struct {
	engine *gin.Engine
}

func NewServerHTTP(adminHandler *handler.AdminHandler, userHandler *handler.UserHandler,
	productHandler *handler.ProductHandler, orderHandler *handler.OrderHandler,
	couponHandler *handler.CouponHandler) *ServerHTTP {

	engine := gin.New()

	engine.LoadHTMLGlob("template/*.html")

	engine.Use(gin.Logger())

	// swagger docs
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// set up routes
	routes.UserRoutes(engine.Group("/"), userHandler, productHandler, orderHandler, couponHandler)
	routes.AdminRoutes(engine.Group("/admin"), adminHandler, productHandler, orderHandler, couponHandler)

	// no handler
	engine.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"StatusCode": 404,
			"msg":        "invalid url",
		})
	})

	return &ServerHTTP{engine: engine}
}

func (s *ServerHTTP) Start() {

	s.engine.Run(":8000")
}
