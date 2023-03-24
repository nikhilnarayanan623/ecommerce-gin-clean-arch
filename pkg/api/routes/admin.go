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
		user := api.Group("/users")
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
		product := api.Group("/products")
		{
			product.GET("/", productHandler.ListProducts)
			product.POST("/", productHandler.AddProducts)
			product.GET("/product-item", productHandler.GetProductItems)
			product.POST("/product-item", productHandler.AddProductItem)

		}
		// order
		order := api.Group("/orders")
		{
			order.GET("/", orderHandler.GetAllShopOrders)
			order.PUT("/", orderHandler.UdateOrderStatus)
		}

		// offer
		offer := api.Group("/offers")
		{
			offer.POST("/", productHandler.AddOffer)     // add a new offer
			offer.GET("/", productHandler.ShowAllOffers) // get all offers
			offer.DELETE("/:offer_id", productHandler.RemoveOffer)

			offer.GET("/category", productHandler.OfferCategoryPage) // to show all offers and all categories
			offer.POST("/category", productHandler.AddOfferCategory) // addd offer for categories

			offer.POST("/category/replace", productHandler.ReplaceOfferCategory)
			offer.POST("/product", productHandler.AddOfferProduct) // add offer for products
		}

	}

}
