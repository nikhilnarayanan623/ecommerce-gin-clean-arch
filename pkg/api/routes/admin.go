package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/middleware"
)

func AdminRoutes(api *gin.RouterGroup, adminHandler *handler.AdminHandler,
	productHandler *handler.ProductHandler, orderHandler *handler.OrderHandler,

) {
	// login
	login := api.Group("/login")
	{
		login.GET("/", adminHandler.LoginGet)
		login.POST("/", adminHandler.LoginPost)
	}
	// signup
	signup := api.Group("/signup")
	{
		signup.GET("/", adminHandler.SignUPGet)
		signup.POST("/", adminHandler.SignUpPost)
	}

	api.Use(middleware.AuthenticateAdmin)
	{
		// user side
		user := api.Group("/user")
		{
			user.GET("/", adminHandler.Allusers)
			user.POST("/block", adminHandler.BlockUser)
		}
		// category
		category := api.Group("/category")
		{
			category.GET("/", productHandler.AllCategories)
			category.POST("/", productHandler.AddCategory)

			category.POST("/variation", productHandler.VariationPost)
			category.POST("/variation-option", productHandler.VariationOptionPost)
		}
		// product
		product := api.Group("/product")
		{
			product.GET("/product", productHandler.ListProducts)
			product.POST("/product", productHandler.AddProducts)
			product.GET("/product-item", productHandler.GetProductItems)
			product.POST("/product-item", productHandler.AddProductItem)

		}
		// order
		order := api.Group("/order")
		{
			order.GET("/orders", orderHandler.GetAllShopOrders)
			order.PUT("/orders", orderHandler.UdateOrderStatus)
		}

	}

}
