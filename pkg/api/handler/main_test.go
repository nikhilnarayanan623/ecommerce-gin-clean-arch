package handler

import (
	"net/http"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	handlerInterface "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/routes"
)

type ServerHTTP struct {
	Engine *gin.Engine
}

func newServerHTTP(authHandler handlerInterface.AuthHandler, adminHandler handlerInterface.AdminHandler, userHandler handlerInterface.UserHandler,
	productHandler handlerInterface.ProductHandler, orderHandler handlerInterface.OrderHandler,
	couponHandler handlerInterface.CouponHandler) *ServerHTTP {

	engine := gin.New()

	// set up routes
	routes.UserRoutes(engine.Group("/"), authHandler, userHandler, productHandler, orderHandler, couponHandler)
	routes.AdminRoutes(engine.Group("/admin"), adminHandler, productHandler, orderHandler, couponHandler)

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

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	os.Exit(m.Run())
}
