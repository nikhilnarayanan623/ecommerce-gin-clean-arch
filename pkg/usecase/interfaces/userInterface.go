package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/res"
)

type UserUseCase interface {
	Signup(ctx context.Context, user domain.User) error
	Login(ctx context.Context, user domain.User) (domain.User, error)
	LoginOtp(ctx context.Context, user domain.User) (domain.User, error)

	Account(ctx context.Context, userId uint) (domain.User, error)
	EditAccount(ctx context.Context, user domain.User) error

	//cart side
	SaveToCart(ctx context.Context, body req.ReqCart) error             // save product_item to cart
	RemoveCartItem(ctx context.Context, body req.ReqCart) error         // remove product_item from cart
	UpdateCartItem(ctx context.Context, body req.ReqCartCount) error    // edit cartItems( quantity change )
	GetCartItems(ctx context.Context, userId uint) (res.ResCart, error) // show all cart_items

	CheckOutCart(ctx context.Context, userID uint) (res.ResCheckOut, error)
	// profile side

	//address side
	SaveAddress(ctx context.Context, address domain.Address, userID uint, isDefault bool) (domain.Address, error) // save address
	EditAddress(ctx context.Context, addressBody req.ReqEditAddress, userID uint) error
	GetAddresses(ctx context.Context, userID uint) ([]res.ResAddress, error) // to get all address of a user

	// wishlist
	AddToWishList(ctx context.Context, wishList domain.WishList) error
	RemoveFromWishList(ctx context.Context, wishList domain.WishList) error
	GetWishListItems(ctx context.Context, userID uint) ([]res.ResWishList, error)
}
