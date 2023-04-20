package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
)

type UserRepository interface {
	FindUser(ctx context.Context, user domain.User) (domain.User, error)
	FindUserByEmail(ctx context.Context, email string) (user domain.User, err error)
	CheckOtherUserWithDetails(ctx context.Context, user domain.User) (domain.User, error) // find user exept this id
	SaveUser(ctx context.Context, user domain.User) (userID uint, err error)
	SaveUserWithGoogleDetails(ctx context.Context, user domain.User) (userID uint, err error)
	UpdateUser(ctx context.Context, user domain.User) (err error)

	// cart
	FindProductItem(ctx context.Context, productItemID uint) (domain.ProductItem, error)

	//cart item
	FindCartByUserID(ctx context.Context, userID uint) (cart domain.Cart, err error)
	SaveCart(ctx context.Context, userID uint) (cart domain.Cart, err error)
	//FindCartTotalPrice(ctx context.Context, userID uint, includeOutOfStck bool) (uint, error)

	FindCartItemByID(ctx context.Context, cartItemID uint) (cartItem domain.CartItem, err error)
	FindCartItemByCartAndProductItemID(ctx context.Context, cartID, productItemID uint) (cartItem domain.CartItem, err error)
	FindAllCartItemsByCartID(ctx context.Context, cartID uint) (cartItems []res.ResCartItem, err error)
	SaveCartItem(ctx context.Context, cartId, productItemId uint) error
	DeleteCartItem(ctx context.Context, cartItemID uint) error
	UpdateCartItemQty(ctx context.Context, cartItemId, qty uint) error

	//address
	FindCountryByID(ctx context.Context, countryID uint) (domain.Country, error)                          // find country by id
	FindAddressByID(ctx context.Context, addressID uint) (domain.Address, error)                          // find address by id
	FindAddressByUserID(ctx context.Context, address domain.Address, userID uint) (domain.Address, error) // find address with userID and addres values
	FindAllAddressByUserID(ctx context.Context, userID uint) ([]res.ResAddress, error)                    // to get all address of user
	SaveAddress(ctx context.Context, address domain.Address) (addressID uint, err error)                  // save a full address
	UpdateAddress(ctx context.Context, address domain.Address) error
	// address join table
	SaveUserAddress(ctx context.Context, userAdress domain.UserAddress) error // save address for user(join table)
	UpdateUserAddress(ctx context.Context, userAddress domain.UserAddress) error

	//wishlist
	FindWishListItem(ctx context.Context, productID, userID uint) (domain.WishList, error)
	FindAllWishListItemsByUserID(ctx context.Context, userID uint) ([]res.ResWishList, error)
	SaveWishListItem(ctx context.Context, wishList domain.WishList) error
	RemoveWishListItem(ctx context.Context, wishList domain.WishList) error
}
