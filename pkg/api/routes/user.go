package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/middleware"
)

func UserRoutes(api *gin.RouterGroup, userHandler *handler.UserHandler, ProductHandler *handler.ProductHandler,
	orderHandler *handler.OrderHandler, couponHandler *handler.CouponHandler,
) {

	// login
	login := api.Group("/login")
	{
		login.GET("/", userHandler.LoginGet)
		login.POST("/", userHandler.LoginPost)
		login.POST("/otp-send", userHandler.LoginOtpSend)
		login.POST("/verify", userHandler.LoginOtpVerify)
	}
	//signup
	signup := api.Group("/signup")
	{
		signup.GET("/", userHandler.SignUpGet)
		signup.POST("/", userHandler.SignUpPost)
	}

	api.Use(middleware.AuthenticateUser)
	{
		api.GET("/", userHandler.Home)
		api.POST("/logout", userHandler.Logout)
		// products
		products := api.Group("/products")
		{
			products.GET("/", ProductHandler.ListProducts)                            // show products
			products.GET("/product-item/:product_id", ProductHandler.GetProductItems) // show product items of a product
		}

		// cart
		cart := api.Group("/carts")
		{
			cart.GET("/", userHandler.UserCart)
			cart.POST("/", userHandler.AddToCart)
			cart.PUT("/", userHandler.UpdateCart)
			cart.DELETE("/", userHandler.RemoveFromCart)

			// place order by cart
			cart.GET("/checkout", userHandler.CheckOutCart)
			cart.POST("/place-order/:address_id", orderHandler.PlaceOrderByCart) // place an order
		}

		//wishlist
		wishList := api.Group("/wishlist")
		{
			wishList.GET("/", userHandler.GetWishListI)
			wishList.POST("/:id", userHandler.AddToWishList)
			wishList.DELETE("/:id", userHandler.RemoveFromWishList)
		}

		// profile
		profile := api.Group("/profile")
		{
			profile.GET("/", userHandler.Account)
			profile.PUT("/", userHandler.EditAccount)

			profile.GET("/address", userHandler.GetAddresses) // to show all address and // show countries
			profile.POST("/address", userHandler.AddAddress)  // to add a new address
			profile.PUT("/address", userHandler.EditAddress)  // to edit address
			profile.DELETE("/address", userHandler.DeleteAddress)
		}

		// order
		orders := api.Group("/orders")
		{
			orders.GET("/", orderHandler.GetUserOrder)                             // get all order list for user
			orders.GET("/items/:shop_order_id", orderHandler.GetOrderItemsForUser) //get order items for specific order

			orders.PUT("/return", orderHandler.SubmitReturnRequest)

			orders.PUT("/cancel/:shop_order_id", orderHandler.CancellOrder) // cancell an order
		}

		// coupons
		coupons := api.Group("/coupons")
		{
			coupons.GET("/", couponHandler.GetAllUserCoupons)
			coupons.POST("/", couponHandler.ApplyUserCoupon)
		}

	}

}
