package routes

import (
	"github.com/gin-gonic/gin"
	handlerInterface "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/middleware"
)

func AdminRoutes(api *gin.RouterGroup, authHandler handlerInterface.AuthHandler, middleware middleware.Middleware,
	adminHandler handlerInterface.AdminHandler, productHandler handlerInterface.ProductHandler,
	paymentHandler handlerInterface.PaymentHandler,
	orderHandler handlerInterface.OrderHandler, couponHandler handlerInterface.CouponHandler,

) {
	// login
	login := api.Group("/login")
	{
		login.POST("/", authHandler.AdminLogin)
	}
	// signup
	signup := api.Group("/signup")
	{
		signup.POST("/", adminHandler.AdminSignUp)
	}
	api.POST("/renew-access-token", authHandler.AdminRenewAccessToken())

	api.Use(middleware.GetAdminAuthMiddleware())
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
			user.GET("/", adminHandler.FindAllUsers)
			user.PATCH("/block", adminHandler.BlockUser)
		}
		// category
		category := api.Group("/category")
		{
			category.GET("/", productHandler.FindAllCategories)
			category.POST("/", productHandler.SaveCategory)
			category.POST("/sub-category", productHandler.SaveSubCategory)

			category.GET("/variations/:category_id", productHandler.FindAllVariations)

			category.POST("/variation", productHandler.SaveVariation)
			category.POST("/variation-option", productHandler.SaveVariationOption)
		}
		// product
		product := api.Group("/products")
		{
			product.GET("/", productHandler.FindAllProducts)
			product.POST("/", productHandler.SaveProduct)
			product.PUT("/", productHandler.UpdateProduct)

			product.GET("/product-item/:product_id", productHandler.FindAllProductItems)
			product.POST("/product-item", productHandler.SaveProductItem)
		}
		// 	// order
		order := api.Group("/orders")
		{
			order.GET("/", orderHandler.FindAllShopOrders)
			order.GET("items", orderHandler.FindAllOrderItems)
			order.PUT("/", orderHandler.UpdateOrderStatus)

			order.GET("/statuses", orderHandler.FindAllOrderStatuses)

			//return requests
			order.GET("/returns", orderHandler.FindAllOrderReturns)
			order.GET("/returns/pending", orderHandler.FindAllPendingReturns)
			order.PUT("/returns/pending", orderHandler.UpdateReturnRequest)
		}

		// payment_method
		paymentMethod := api.Group("/payment-method")
		{
			paymentMethod.GET("/", paymentHandler.FindAllPaymentMethods)
			paymentMethod.POST("/", paymentHandler.AddPaymentMethod)
			paymentMethod.PUT("/", paymentHandler.UpdatePaymentMethod)
		}

		// 	// offer
		// 	offer := api.Group("/offers")
		// 	{
		// 		offer.POST("/", productHandler.AddOffer)     // add a new offer
		// 		offer.GET("/", productHandler.ShowAllOffers) // get all offers
		// 		offer.DELETE("/:offer_id", productHandler.RemoveOffer)

		// 		offer.GET("/category", productHandler.GetOfferCategories) // to get all offers of categories
		// 		offer.POST("/category", productHandler.AddOfferCategory)  // add offer for categories
		// 		offer.PUT("/category", productHandler.ReplaceOfferCategory)
		// 		offer.DELETE("/category/:offer_category_id", productHandler.RemoveOfferCategory)

		// 		offer.GET("/products", productHandler.GetOffersOfProducts) // to get all offers of products
		// 		offer.POST("/products", productHandler.AddOfferProduct)    // add offer for products
		// 		offer.PUT("/products", productHandler.ReplaceOfferProduct)
		// 		offer.DELETE("/products/:offer_product_id", productHandler.RemoveOfferProduct)
		// 	}

		// 	// coupons
		// 	coupons := api.Group("/coupons")
		// 	{
		// 		coupons.POST("/", couponHandler.AddCoupon)
		// 		coupons.GET("/", couponHandler.GetAllCoupons)
		// 		coupons.PUT("/", couponHandler.UpdateCoupon)
		// 	}

		// 	stock := api.Group("/stocks")
		// 	{
		// 		stock.GET("/", adminHandler.FindAllStockDetails)

		// 		stock.PATCH("/", adminHandler.UpdateStock)
		// 	}

	}

}
