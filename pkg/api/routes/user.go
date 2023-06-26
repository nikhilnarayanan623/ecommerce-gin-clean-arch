package routes

import (
	"github.com/gin-gonic/gin"
	handlerInterface "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/middleware"
)

func UserRoutes(api *gin.RouterGroup, authHandler handlerInterface.AuthHandler, middleware middleware.Middleware,
	userHandler handlerInterface.UserHandler, cartHandler handlerInterface.CartHandler,
	productHandler handlerInterface.ProductHandler, paymentHandler handlerInterface.PaymentHandler,
	orderHandler handlerInterface.OrderHandler, couponHandler handlerInterface.CouponHandler,
) {

	auth := api.Group("/auth")
	{
		signup := auth.Group("/signup")
		{
			signup.POST("/", authHandler.UserSignUp)
		}

		login := auth.Group("/login")
		{
			login.POST("/", authHandler.UserLogin)
			login.POST("/otp-send", authHandler.UserLoginOtpSend)
			login.POST("/otp-verify", authHandler.UserLoginOtpVerify)
		}

		goath := auth.Group("/google-auth")
		{
			goath.GET("/", authHandler.UserGoogleAuthLoginPage)
			goath.GET("/initialize", authHandler.UserGoogleAuthInitialize)
			goath.GET("/callback", authHandler.UserGoogleAuthCallBack)
		}

		auth.POST("/renew-access-token", authHandler.UserRenewAccessToken())

	}

	api.Use(middleware.GetUserAuthMiddleware())
	{

		api.GET("/", userHandler.Home)
		// api.POST("/logout", userHandler.UserLogout)

		product := api.Group("/products")
		{
			product.GET("/", productHandler.FindAllProductsUser())

			productItem := product.Group("/:product_id/items")
			{
				productItem.GET("/", productHandler.FindAllProductItemsUser())
			}
		}

		// 	// cart
		cart := api.Group("/carts")
		{
			cart.GET("/", cartHandler.FindCart)
			cart.POST("/:product_item_id", cartHandler.AddToCart)
			cart.PUT("/", cartHandler.UpdateCart)
			cart.DELETE("/:product_item_id", cartHandler.RemoveFromCart)

			// 		cart.PATCH("/apply-coupon", couponHandler.ApplyCouponToCart)

			// 		cart.GET("/payment-methods", orderHandler.GetAllPaymentMethods)
			cart.GET("/checkout/payment-select-page", paymentHandler.CartOrderPaymentSelectPage)

			cart.POST("/place-order", orderHandler.PlaceOrderOnCOD)
			// cart.POST("/place-order/cod", orderHandler.ApproveOrderCOD)

			// 		//cart.GET("/checkout", userHandler.CheckOutCart, orderHandler.GetAllPaymentMethods)

			// 		// make razorpay order and verify
			cart.POST("/place-order/razorpay-checkout", orderHandler.RazorpayCheckout)
			cart.POST("/place-order/razorpay-verify", orderHandler.RazorpayVerify)

			// 		// stripe
			// 		cart.POST("/place-order/stripe-checkout", orderHandler.StripPaymentCheckout)
			// 		cart.POST("/place-order/stripe/stripe-verify", orderHandler.StripePaymentVeify)
		}

		//wishlist
		wishList := api.Group("/wishlist")
		{
			wishList.GET("/", userHandler.FindWishList)
			wishList.POST("/:product_item_id", userHandler.SaveToWishList)
			wishList.DELETE("/:product_item_id", userHandler.RemoveFromWishList)
		}

		// profile
		account := api.Group("/account")
		{
			account.GET("/", userHandler.FindProfile)
			account.PUT("/", userHandler.UpdateProfile)

			account.GET("/address", userHandler.FindAllAddresses) // to show all address and // show countries
			account.POST("/address", userHandler.SaveAddress)     // to add a new address
			account.PUT("/address", userHandler.UpdateAddress)    // to edit address
			// account.DELETE("/address", userHandler.DeleteAddress)

			// wallet for user
			account.GET("/wallet", orderHandler.FindUserWallet)
			account.GET("/wallet/transactions", orderHandler.FindUserWalletTransactions)
		}

		// 	// order
		orders := api.Group("/orders")
		{
			orders.GET("/", orderHandler.FindUserOrder)          // get all order list for user
			orders.GET("/items", orderHandler.FindAllOrderItems) //get order items for specific order

			orders.POST("/return", orderHandler.SubmitReturnRequest)
			orders.POST("/cancel/:shop_order_id", orderHandler.CancelOrder) // cancell an order
		}

		//coupons
		coupons := api.Group("/coupons")
		{
			coupons.GET("/", couponHandler.FindAllCouponsForUser)
		}

	}

}
