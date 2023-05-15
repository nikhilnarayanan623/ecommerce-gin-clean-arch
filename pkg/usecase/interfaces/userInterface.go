package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
)

type UserUseCase interface {
	Account(ctx context.Context, userId uint) (domain.User, error)
	EditAccount(ctx context.Context, user domain.User) error

	// profile side

	//address side
	SaveAddress(ctx context.Context, userID uint, address domain.Address, isDefault bool) error // save address
	EditAddress(ctx context.Context, addressBody req.EditAddress, userID uint) error
	GetAddresses(ctx context.Context, userID uint) ([]res.Address, error) // to get all address of a user

	// wishlist
	AddToWishList(ctx context.Context, wishList domain.WishList) error
	RemoveFromWishList(ctx context.Context, wishList domain.WishList) error
	GetWishListItems(ctx context.Context, userID uint) ([]res.WishList, error)
}
