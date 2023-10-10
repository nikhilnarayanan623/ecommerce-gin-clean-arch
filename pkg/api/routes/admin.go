package routes

import (
	"github.com/gin-gonic/gin"
	handlerInterface "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/middleware"
)

func AdminRoutes(api *gin.RouterGroup, authHandler handlerInterface.AuthHandler, middleware middleware.Middleware,
	adminHandler handlerInterface.AdminHandler, productHandler handlerInterface.ProductHandler,
	paymentHandler handlerInterface.PaymentHandler, orderHandler handlerInterface.OrderHandler,
	couponHandler handlerInterface.CouponHandler, offerHandler handlerInterface.OfferHandler,
	stockHandler handlerInterface.StockHandler, branHandler handlerInterface.BrandHandler,

) {

	auth := api.Group("/auth")
	{
		login := auth.Group("/sign-in")
		{
			login.POST("/", authHandler.AdminLogin)
		}

		// signup := api.Group("/signup")
		// {
		// 	signup.POST("/", adminHandler.AdminSignUp)
		// }

		auth.POST("/renew-access-token", authHandler.AdminRenewAccessToken())
	}

	api.Use(middleware.AuthenticateAdmin())
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
			category.POST("/", middleware.TrimSpaces(), productHandler.SaveCategory)
			category.POST("/sub-categories", middleware.TrimSpaces(), productHandler.SaveSubCategory)

			variation := category.Group("/:category_id/variations")
			{
				variation.POST("/", middleware.TrimSpaces(), productHandler.SaveVariation)
				variation.GET("/", productHandler.GetAllVariations)

				variationOption := variation.Group("/:variation_id/options")
				{
					variationOption.POST("/", middleware.TrimSpaces(), productHandler.SaveVariationOption)
				}
			}

		}
		// brand
		brand := api.Group("/brands")
		{
			brand.POST("", branHandler.Save)
			brand.GET("", branHandler.FindAll)
			brand.GET("/:brand_id", branHandler.FindOne)
			brand.PUT("/:brand_id", branHandler.Update)
			brand.DELETE("/:brand_id", branHandler.Delete)
		}

		// product
		product := api.Group("/products")
		{
			product.GET("/", productHandler.GetAllProductsAdmin())
			product.POST("/", middleware.TrimSpaces(), productHandler.SaveProduct)
			product.PUT("/", middleware.TrimSpaces(), productHandler.UpdateProduct)

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
			offer.POST("/", middleware.TrimSpaces(), offerHandler.SaveOffer) // add a new offer
			offer.GET("/", offerHandler.GetAllOffers)                        // get all offers
			offer.DELETE("/:offer_id", offerHandler.RemoveOffer)

			offer.GET("/category", offerHandler.GetAllCategoryOffers)                        // to get all offers of categories
			offer.POST("/category", middleware.TrimSpaces(), offerHandler.SaveCategoryOffer) // add offer for categories
			offer.PATCH("/category", offerHandler.ChangeCategoryOffer)
			offer.DELETE("/category/:offer_category_id", offerHandler.RemoveCategoryOffer)

			offer.GET("/products", offerHandler.GetAllProductsOffers)                       // to get all offers of products
			offer.POST("/products", middleware.TrimSpaces(), offerHandler.SaveProductOffer) // add offer for products
			offer.PATCH("/products", offerHandler.ChangeProductOffer)
			offer.DELETE("/products/:offer_product_id", offerHandler.RemoveProductOffer)
		}

		// coupons
		coupons := api.Group("/coupons")
		{
			coupons.POST("/", middleware.TrimSpaces(), couponHandler.SaveCoupon)
			coupons.GET("/", couponHandler.GetAllCouponsAdmin)
			coupons.PUT("/", middleware.TrimSpaces(), couponHandler.UpdateCoupon)
		}

		// sales report
		sales := api.Group("/sales")
		{
			sales.GET("/", adminHandler.GetFullSalesReport)
		}

		stock := api.Group("/stocks")
		{
			stock.GET("/", stockHandler.GetAllStocks)

			stock.PATCH("/", stockHandler.UpdateStock)
		}

	}

}
