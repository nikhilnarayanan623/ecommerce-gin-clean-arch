package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
)

type UserRepository interface {
	FindUser(ctx context.Context, user domain.User) (domain.User, error)
	FindUserExceptID(ctx context.Context, user domain.User) (domain.User, error) // find user exept this id
	SaveUser(ctx context.Context, user domain.User) error
	EditUser(ctx context.Context, user domain.User) error
	// cart
	FindProductItem(ctx context.Context, productItemID uint) (domain.ProductItem, error)

	//cart item
	FindCart(ctx context.Context, cart domain.Cart) (domain.Cart, error)
	FindCartTotalPrice(ctx context.Context, userID uint, includeOutOfStck bool) (uint, error)
	FindAllCartItems(ctx context.Context, userID uint) ([]res.ResCartItem, error)
	SaveCartItem(ctx context.Context, cartItems domain.Cart) error
	RemoveCartItem(ctx context.Context, cartItem domain.Cart) error
	UpdateCartItem(ctx context.Context, cartItem domain.Cart) error

	// checkout for order
	CheckOutCart(ctx context.Context, userId uint) (res.ResCheckOut, error)

	//address
	FindCountryByID(ctx context.Context, countryID uint) (domain.Country, error)                          // find country by id
	FindAddressByID(ctx context.Context, addressID uint) (domain.Address, error)                          // find address by id
	FindAddressByUserID(ctx context.Context, address domain.Address, userID uint) (domain.Address, error) // find address with userID and addres values
	FindAllAddressByUserID(ctx context.Context, userID uint) ([]res.ResAddress, error)                    // to get all address of user
	SaveAddress(ctx context.Context, address domain.Address) (domain.Address, error)                      // save a full address
	UpdateAddress(ctx context.Context, address domain.Address) error
	// address join table
	SaveUserAddress(ctx context.Context, userAdress domain.UserAddress) (domain.UserAddress, error) // save address for user(join table)
	UpdateUserAddress(ctx context.Context, userAddress domain.UserAddress) error

	//wishlist
	FindWishListItem(ctx context.Context, productID, userID uint) (domain.WishList, error)
	FindAllWishListItemsByUserID(ctx context.Context, userID uint) ([]res.ResWishList, error)
	SaveWishListItem(ctx context.Context, wishList domain.WishList) error
	RemoveWishListItem(ctx context.Context, wishList domain.WishList) error
}
