package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper"
)

type UserUseCase interface {
	Signup(ctx context.Context, user domain.User) (domain.User, error)
	Login(ctx context.Context, user domain.User) (domain.User, error)
	LoginOtp(ctx context.Context, user domain.User) (domain.User, error)

	Home(ctx context.Context, userId uint) (domain.User, error)

	//cart side
	SaveToCart(ctx context.Context, body helper.ReqCart) (domain.CartItem, error)          // save product_item to cart
	RemoveCartItem(ctx context.Context, body helper.ReqCart) (domain.Cart, error)          // remove product_item from cart
	UpdateCartItem(ctx context.Context, body helper.ReqCartCount) (domain.CartItem, error) // edit cartItems( quantity change )
	GetCartItems(ctx context.Context, userId uint) (helper.ResponseCart, error)            // show all cart_items

	// profile side

	//address side
	SaveAddress(ctx context.Context, address domain.Address, userID uint, isDefault bool) (domain.Address, error) // save address
	GetAddresses(ctx context.Context, userID uint) ([]helper.ResAddress, error)                                   // to get all address of a user
}
