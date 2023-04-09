package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
)

type UserUseCase interface {
	Signup(ctx context.Context, user domain.User) error
	Login(ctx context.Context, user domain.User) (domain.User, error)
	LoginOtp(ctx context.Context, user domain.User) (domain.User, error)

	Account(ctx context.Context, userId uint) (domain.User, error)
	EditAccount(ctx context.Context, user domain.User) error

	//cart side
	SaveToCart(ctx context.Context, body req.ReqCart) error          // save product_item to cart
	RemoveCartItem(ctx context.Context, body req.ReqCart) error      // remove product_item from cart
	UpdateCartItem(ctx context.Context, body req.ReqCartCount) error // edit cartItems( quantity change )
	GetUserCart(ctx context.Context, userID uint) (cart domain.Cart, err error)
	GetUserCartItems(ctx context.Context, cartId uint) (cartItems []res.ResCartItem, err error) // show all cart_items

	// profile side

	//address side
	SaveAddress(ctx context.Context, userID uint, address domain.Address, isDefault bool) error // save address
	EditAddress(ctx context.Context, addressBody req.ReqEditAddress, userID uint) error
	GetAddresses(ctx context.Context, userID uint) ([]res.ResAddress, error) // to get all address of a user

	// wishlist
	AddToWishList(ctx context.Context, wishList domain.WishList) error
	RemoveFromWishList(ctx context.Context, wishList domain.WishList) error
	GetWishListItems(ctx context.Context, userID uint) ([]res.ResWishList, error)
}
