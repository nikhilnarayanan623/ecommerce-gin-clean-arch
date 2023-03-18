package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper"
)

type UserRepository interface {
	FindUser(ctx context.Context, user domain.User) (domain.User, error)
	SaveUser(ctx context.Context, user domain.User) (domain.User, error)
	// cart
	FindProductItem(ctx context.Context, productItemID uint) (domain.ProductItem, error)
	FindCart(ctx context.Context, userId uint) (domain.Cart, error)
	UpdateCartPrice(ctx context.Context, cart domain.Cart) (domain.Cart, error)
	//cart item
	FindCartItem(ctx context.Context, cartID, productItemID uint) (domain.CartItem, error)
	SaveCartItem(ctx context.Context, cartID, productItemID uint) (domain.CartItem, error)
	RemoveCartItem(ctx context.Context, cartItem domain.CartItem) (domain.CartItem, error)
	UpdateCartItem(ctx context.Context, cartItem domain.CartItem) (domain.CartItem, error)
	GetCartItems(ctx context.Context, userId uint) (helper.ResponseCart, error)
	//address
	FindCountryByID(ctx context.Context, countryID uint) (domain.Country, error)                          // find country by id
	FindAddressByID(ctx context.Context, addressID uint) (domain.Address, error)                          // find address by id
	FindAddressByUserID(ctx context.Context, address domain.Address, userID uint) (domain.Address, error) // find address with userID and addres values
	FindAllAddressByUserID(ctx context.Context, userID uint) ([]helper.ResAddress, error)                 // to get all address of user
	SaveAddress(ctx context.Context, address domain.Address) (domain.Address, error)                      // save a full address
	UpdateAddress(ctx context.Context, address domain.Address) error
	//join table
	SaveUserAddress(ctx context.Context, userAdress domain.UserAddress) (domain.UserAddress, error) // save address for user(join table)
	UpdateUserAddress(ctx context.Context, userAddress domain.UserAddress) error
}
