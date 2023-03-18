package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/middleware"
)

func UserRoutes(api *gin.RouterGroup, user *handler.UserHandler, product *handler.ProductHandler) {

	api.GET("/login", user.LoginGet)
	api.POST("/login", user.LoginPost)
	api.POST("/login-otp-send", user.LoginOtpSend)
	api.POST("/login-otp-verify", user.LoginOtpVerify)

	api.GET("/signup", user.SignUpGet)
	api.POST("/signup", user.SignUpPost)

	api.Use(middleware.AuthenticateUser)
	{
		api.GET("/", user.Home)
		api.GET("/product", product.ListProducts)         // show products
		api.GET("/product-item", product.GetProductItems) // show product items of a product
		// cart
		api.GET("/cart", user.UserCart)
		api.POST("/cart", user.AddToCart)
		api.PUT("/cart", user.UpdateCart)
		api.DELETE("/cart", user.RemoveFromCart)

		//wishlist
		api.GET("/wishlist")

		//profile address
		api.GET("/profile/address", user.GetAddresses) // to show all address and // show countries
		api.POST("/profile/address", user.AddAddress)  // to add a new address
		api.PUT("/profile/address", user.EditAddress)  // to edit address
		api.DELETE("profile/address", user.DeleteAddress)

		api.POST("/logout", user.Logout)
	}

}
