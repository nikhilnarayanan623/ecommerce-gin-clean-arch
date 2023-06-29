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

	auth := api.Group("/auth")
	{
		login := auth.Group("/login")
		{
			login.POST("/", authHandler.AdminLogin)
		}

		// signup := api.Group("/signup")
		// {
		// 	signup.POST("/", adminHandler.AdminSignUp)
		// }

		auth.POST("/renew-access-token", authHandler.AdminRenewAccessToken())
	}

	api.Use(middleware.GetAdminAuthMiddleware())
	{

		// user side
		user := api.Group("/users")
		{
			user.GET("/", adminHandler.GetAllUsers)
			user.PATCH("/block", adminHandler.BlockUser)
		}
		// category
		category := api.Group("/categories")
		{
			category.GET("/", productHandler.GetAllCategories)
			category.POST("/", productHandler.SaveCategory)
			category.POST("/sub-categories", productHandler.SaveSubCategory)

			variation := category.Group("/:category_id/variations")
			{
				variation.POST("/", productHandler.SaveVariation)
				variation.GET("/", productHandler.GetAllVariations)

				variationOption := variation.Group("/:variation_id/options")
				{
					variationOption.POST("/", productHandler.SaveVariationOption)
				}
			}

		}
		// product
		product := api.Group("/products")
		{
			product.GET("/", productHandler.GetAllProductsAdmin())
			product.POST("/", productHandler.SaveProduct)
			product.PUT("/", productHandler.UpdateProduct)

			productItem := product.Group("/:product_id/items")
			{
				productItem.GET("/", productHandler.GetAllProductItemsAdmin())
				productItem.POST("/", productHandler.SaveProductItem)
			}
		}
		// 	// order
		order := api.Group("/orders")
		{
			order.GET("/all", orderHandler.GetAllShopOrders)
			order.GET("/:shop_order_id/items", orderHandler.GetAllOrderItemsAdmin())
			order.PUT("/", orderHandler.UpdateOrderStatus)

			status := order.Group("/statuses")
			{
				status.GET("/", orderHandler.GetAllOrderStatuses)
			}

			//return requests
			order.GET("/returns", orderHandler.GetAllOrderReturns)
			order.GET("/returns/pending", orderHandler.GetAllPendingReturns)
			order.PUT("/returns/pending", orderHandler.UpdateReturnRequest)
		}

		// payment_method
		paymentMethod := api.Group("/payment-methods")
		{
			paymentMethod.GET("/", paymentHandler.GetAllPaymentMethodsAdmin())
			// paymentMethod.POST("/", paymentHandler.AddPaymentMethod)
			paymentMethod.PUT("/", paymentHandler.UpdatePaymentMethod)
		}

		// offer
		offer := api.Group("/offers")
		{
			offer.POST("/", productHandler.SaveOffer)   // add a new offer
			offer.GET("/", productHandler.GetAllOffers) // get all offers
			offer.DELETE("/:offer_id", productHandler.RemoveOffer)

			offer.GET("/category", productHandler.GetAllCategoryOffers) // to get all offers of categories
			offer.POST("/category", productHandler.SaveCategoryOffer)   // add offer for categories
			offer.PATCH("/category", productHandler.ChangeCategoryOffer)
			offer.DELETE("/category/:offer_category_id", productHandler.RemoveCategoryOffer)

			offer.GET("/products", productHandler.GetAllProductsOffers) // to get all offers of products
			offer.POST("/products", productHandler.SaveProductOffer)    // add offer for products
			offer.PATCH("/products", productHandler.ChangeProductOffer)
			offer.DELETE("/products/:offer_product_id", productHandler.RemoveProductOffer)
		}

		// coupons
		coupons := api.Group("/coupons")
		{
			coupons.POST("/", couponHandler.SaveCoupon)
			coupons.GET("/", couponHandler.GetAllCouponsAdmin)
			coupons.PUT("/", couponHandler.UpdateCoupon)
		}

		// sales report
		sales := api.Group("/sales")
		{
			sales.GET("/", adminHandler.GetFullSalesReport)
		}

		stock := api.Group("/stocks")
		{
			stock.GET("/", adminHandler.GetAllStocks)

			stock.PATCH("/", adminHandler.UpdateStock)
		}

	}

}
