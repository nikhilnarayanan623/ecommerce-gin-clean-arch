package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/middleware"
)

func AdminRoutes(api *gin.RouterGroup, admin *handler.AdminHandler, product *handler.ProductHandler) {

	api.GET("/login", admin.LoginGet)
	api.POST("/login", admin.LoginPost)

	api.GET("/signup", admin.SignUPGet)
	api.POST("/signup", admin.SignUpPost)

	api.Use(middleware.AuthenticateAdmin)
	{
		api.GET("/users", admin.Allusers)
		api.POST("/block-user", admin.BlockUser)
		api.GET("/category", product.AllCategories)
		api.POST("/category", product.AddCategory)
		// to add a new variation for category
		api.POST("/variation", product.VariationPost)

		// to add a new variation value for variation
		api.POST("/variation-option", product.VariationOptionPost)

		api.GET("/product", product.ListProducts)
		api.POST("/product", product.AddProducts)

		api.GET("/product-item", product.GetProductItems)
		api.POST("/product-item", product.AddProductItem)

	}

}
