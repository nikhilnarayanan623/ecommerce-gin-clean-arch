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
		signup := auth.Group("/sign-up")
		{
			signup.POST("/", authHandler.UserSignUp)
			signup.POST("/verify", authHandler.UserSignUpVerify)
		}

		login := auth.Group("/sign-in")
		{
			login.POST("/", authHandler.UserLogin)
			login.POST("/otp/send", authHandler.UserLoginOtpSend)
			login.POST("/otp/verify", authHandler.UserLoginOtpVerify)
		}

		goath := auth.Group("/google-auth")
		{
			goath.GET("/", authHandler.UserGoogleAuthLoginPage)
			goath.GET("/initialize", authHandler.UserGoogleAuthInitialize)
			goath.GET("/callback", authHandler.UserGoogleAuthCallBack)
		}

		auth.POST("/renew-access-token", authHandler.UserRenewAccessToken())

		// api.POST("/logout")

	}

	api.Use(middleware.AuthenticateUser())
	{

		// api.POST("/logout", userHandler.UserLogout)

		product := api.Group("/products")
		{
			product.GET("/", productHandler.GetAllProductsUser())

			productItem := product.Group("/:product_id/items")
			{
				productItem.GET("/", productHandler.GetAllProductItemsUser())
			}
		}

		// 	// cart
		cart := api.Group("/carts")
		{
			cart.GET("/", cartHandler.GetCart)
			cart.POST("/:product_item_id", cartHandler.AddToCart)
			cart.PUT("/", cartHandler.UpdateCart)
			cart.DELETE("/:product_item_id", cartHandler.RemoveFromCart)

			cart.PATCH("/apply-coupon", couponHandler.ApplyCouponToCart)

			cart.GET("/checkout/payment-select-page", paymentHandler.CartOrderPaymentSelectPage)
			// 		cart.GET("/payment-methods", orderHandler.GetAllPaymentMethods)
			cart.POST("/place-order", orderHandler.SaveOrder)

			// 		//cart.GET("/checkout", userHandler.CheckOutCart, orderHandler.GetAllPaymentMethods)
			cart.POST("/place-order/cod", paymentHandler.PaymentCOD)

			// razorpay payment
			cart.POST("/place-order/razorpay-checkout", paymentHandler.RazorpayCheckout)
			cart.POST("/place-order/razorpay-verify", paymentHandler.RazorpayVerify)

			// 	stripe payment
			cart.POST("/place-order/stripe-checkout", paymentHandler.StripPaymentCheckout)
			cart.POST("/place-order/stripe-verify", paymentHandler.StripePaymentVeify)
		}

		// profile
		account := api.Group("/account")
		{
			account.GET("/", userHandler.GetProfile)
			account.PUT("/", userHandler.UpdateProfile)

			account.GET("/address", userHandler.GetAllAddresses) // to show all address and // show countries
			account.POST("/address", userHandler.SaveAddress)    // to add a new address
			account.PUT("/address", userHandler.UpdateAddress)   // to edit address
			// account.DELETE("/address", userHandler.DeleteAddress)

			//wishlist
			wishList := account.Group("/wishlist")
			{
				wishList.GET("/", userHandler.GetWishList)
				wishList.POST("/:product_item_id", userHandler.SaveToWishList)
				wishList.DELETE("/:product_item_id", userHandler.RemoveFromWishList)
			}

			wallet := account.Group("/wallet")
			{
				wallet.GET("/", orderHandler.GetUserWallet)
				wallet.GET("/transactions", orderHandler.GetUserWalletTransactions)
			}

			coupons := account.Group("/coupons")
			{
				coupons.GET("/", couponHandler.GetAllCouponsForUser)
			}
		}

		paymentMethod := api.Group("/payment-methods")
		{
			paymentMethod.GET("/", paymentHandler.GetAllPaymentMethodsUser())
		}

		// 	// order
		orders := api.Group("/orders")
		{
			orders.GET("/", orderHandler.GetUserOrder)                               // get all order list for user
			orders.GET("/:shop_order_id/items", orderHandler.GetAllOrderItemsUser()) //get order items for specific order

			orders.POST("/return", orderHandler.SubmitReturnRequest)
			orders.POST("/:shop_order_id/cancel", orderHandler.CancelOrder) // cancel an order
		}

	}

}
