package routes

import (
	"github.com/gin-gonic/gin"
	handlerInterface "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/middleware"
)

func UserRoutes(api *gin.RouterGroup, authHandler handlerInterface.AuthHandler, middleware middleware.Middleware,
	userHandler handlerInterface.UserHandler, cartHandler handlerInterface.CartHandler,
	ProductHandler handlerInterface.ProductHandler,
	orderHandler handlerInterface.OrderHandler, couponHandler handlerInterface.CouponHandler,
) {

	// login
	login := api.Group("/login")
	{
		login.POST("/", authHandler.UserLogin)
		// login.POST("/otp-send", userHandler.UserLoginOtpSend)
		// login.POST("/otp-verify", userHandler.UserLoginOtpVerify)

		// // 	// login page
		// login.GET("/", userHandler.GoauthLoginPage)
		// login.GET("/auth/", userHandler.IntitializeGoogleAuth)
		// login.GET("/auth/google/callback", userHandler.CallbackAuth)
	}
	api.POST("/renew-access-token", authHandler.UserRenewRefreshToken())
	// //signup
	// signup := api.Group("/signup")
	// {
	// 	signup.POST("/", userHandler.UserSignUp)
	// }

	api.Use(middleware.GetUserMiddleware())
	{

		api.GET("/", userHandler.Home)
		// 	api.POST("/logout", userHandler.UserLogout)

		// 	// products
		// 	products := api.Group("/products")
		// 	{
		// 		products.GET("/", ProductHandler.ListProducts)                            // show products
		// 		products.GET("/product-item/:product_id", ProductHandler.GetProductItems) // show product items of a product
		// 	}

		// 	// cart
		cart := api.Group("/carts")
		{
			cart.GET("/", cartHandler.ShowCart)
			cart.POST("/", cartHandler.AddToCart)
			cart.PUT("/", cartHandler.UpdateCart)
			cart.DELETE("/", cartHandler.RemoveFromCart)

			// 		cart.PATCH("/apply-coupon", couponHandler.ApplyCouponToCart)

			// 		cart.GET("/paymet-methods", orderHandler.GetAllPaymentMethods)

			cart.POST("/place-order/cod", orderHandler.PlaceOrderCartCOD) // place an order

			// 		//cart.GET("/checkout", userHandler.CheckOutCart, orderHandler.GetAllPaymentMethods)

			// 		// page for select payment method
			// 		cart.GET("/checkout/payemt-select-page", orderHandler.CartOrderPayementSelectPage)

			// 		// make razorpay order and verify
			// 		cart.POST("/place-order/razorpay-checkout", orderHandler.RazorpayCheckout)
			// 		cart.POST("/place-order/razorpay-verify", orderHandler.RazorpayVerify)

			// 		// stripe
			// 		cart.POST("/place-order/stripe-checkout", orderHandler.StripPaymentCheckout)
			// 		cart.POST("/place-order/stripe/stripe-verify", orderHandler.StripePaymentVeify)
		}

		// 	//wishlist
		// 	wishList := api.Group("/wishlist")
		// 	{
		// 		wishList.GET("/", userHandler.GetWishListI)
		// 		wishList.POST("/:id", userHandler.AddToWishList)
		// 		wishList.DELETE("/:id", userHandler.RemoveFromWishList)
		// 	}

		// 	// profile
		// 	account := api.Group("/account")
		// 	{
		// 		account.GET("/", userHandler.Account)
		// 		account.PUT("/", userHandler.UpateAccount)

		// 		account.GET("/address", userHandler.GetAddresses) // to show all address and // show countries
		// 		account.POST("/address", userHandler.AddAddress)  // to add a new address
		// 		account.PUT("/address", userHandler.EditAddress)  // to edit address
		// 		account.DELETE("/address", userHandler.DeleteAddress)

		// 		// wallet for user
		// 		account.GET("/wallet", orderHandler.GetUserWallet)
		// 		account.GET("/wallet/transactions", orderHandler.GetUserWalletTransactions)
		// 	}

		// 	// order
		orders := api.Group("/orders")
		{
			orders.GET("/", orderHandler.GetUserOrder)                       // get all order list for user
			orders.GET("/items", orderHandler.GetOrderItemsByShopOrderItems) //get order items for specific order

			orders.POST("/return", orderHandler.SubmitReturnRequest)

			orders.POST("/cancel/:shop_order_id", orderHandler.CancellOrder) // cancell an order
		}

		// 	//coupons
		// 	coupons := api.Group("/coupons")
		// 	{
		// 		coupons.GET("/", couponHandler.GetAllCouponsForUser)
		// 	}

	}

}
