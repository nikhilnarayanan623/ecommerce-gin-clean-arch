package http

import (
	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/routes"
)

type ServerHTTP struct {
	engine *gin.Engine
}

func NewServerHTTP(adminHandler *handler.AdminHandler, userHandler *handler.UserHandler, productHandler *handler.ProductHandler) *ServerHTTP {

	engine := gin.New()

	engine.Use(gin.Logger())

	// admin := engine.Group("/admin")
	// user := engine.Group("/")
	//two main routes `\` -> user ; `\admin`-> admin
	routes.UserRoutes(engine.Group("/"), userHandler, productHandler)
	routes.AdminRoutes(engine.Group("/admin"), adminHandler, productHandler)

	return &ServerHTTP{engine: engine}
}

func (s *ServerHTTP) Start() {

	s.engine.Run(":8000")
}
