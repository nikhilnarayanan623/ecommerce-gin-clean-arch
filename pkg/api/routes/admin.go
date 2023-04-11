package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/middleware"
)

func AdminRoutes(api *gin.RouterGroup, adminHandler *handler.AdminHandler,
	productHandler *handler.ProductHandler, orderHandler *handler.OrderHandler,
	couponHandler *handler.CouponHandler,

) {
	// login
	login := api.Group("/login")
	{
		login.POST("/", adminHandler.AdminLogin)
	}
	// signup
	signup := api.Group("/signup")
	{
		signup.POST("/", adminHandler.AdminSignUp)
	}

	api.Use(middleware.AuthenticateAdmin)
	{
		api.GET("/", adminHandler.AdminHome)

		// sales report
		sales := api.Group("/sales")
		{
			sales.GET("/", adminHandler.FullSalesReport)
		}
		// user side
		user := api.Group("/users")
		{
			user.GET("/", adminHandler.ListUsers)
			user.PATCH("/block", adminHandler.BlockUser)
		}
		// category
		category := api.Group("/category")
		{
			category.GET("/", productHandler.GetAlllCategories)
			category.POST("/", productHandler.AddCategory)

			category.POST("/variation", productHandler.AddVariation)
			category.POST("/variation-option", productHandler.AddVariationOption)
		}
		// product
		product := api.Group("/products")
		{
			product.GET("/", productHandler.ListProducts)
			product.POST("/", productHandler.AddProducts)
			product.PUT("/", productHandler.UpdateProduct)

			product.GET("/product-item/:product_id", productHandler.GetProductItems)
			product.POST("/product-item", productHandler.AddProductItem)
		}
		// order
		order := api.Group("/orders")
		{
			order.GET("/", orderHandler.GetAllShopOrders)
			order.GET("items", orderHandler.GetOrderItemsByShopOrderItems)
			order.PUT("/", orderHandler.UdateOrderStatus)

			order.GET("/statuses", orderHandler.GetAllOrderStatuses)

			//return requests
			order.GET("/returns", orderHandler.GetAllOrderReturns)
			order.GET("/returns/pending", orderHandler.GetAllPendingReturns)
			order.PUT("/returns/pending", orderHandler.UpdateReturnRequest)
		}

		// offer
		offer := api.Group("/offers")
		{
			offer.POST("/", productHandler.AddOffer)     // add a new offer
			offer.GET("/", productHandler.ShowAllOffers) // get all offers
			offer.DELETE("/:offer_id", productHandler.RemoveOffer)

			offer.GET("/category", productHandler.GetOfferCategories) // to get all offers of categories
			offer.POST("/category", productHandler.AddOfferCategory)  // addd offer for categories
			offer.PUT("/category", productHandler.ReplaceOfferCategory)
			offer.DELETE("/category/:offer_category_id", productHandler.RemoveOfferCategory)

			offer.GET("/products", productHandler.GetOffersOfProducts) // to get all offers of products
			offer.POST("/products", productHandler.AddOfferProduct)    // add offer for products
			offer.PUT("/products", productHandler.ReplaceOfferProduct)
			offer.DELETE("/products/:offer_product_id", productHandler.RemoveOfferProduct)
		}

		// coupons
		coupons := api.Group("/coupons")
		{
			coupons.POST("/", couponHandler.AddCoupon)
			coupons.GET("/", couponHandler.GetAllCoupons)
			coupons.PUT("/", couponHandler.UpdateCoupon)
		}

		stok := api.Group("/stocks")
		{
			stok.GET("/", adminHandler.GetAllStockDetails)

			stok.PUT("/", adminHandler.UpdateStock)
		}

	}

}
