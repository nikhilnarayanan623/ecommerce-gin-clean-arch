package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
)

//go:generate mockgen -destination=../../mock/mockRepository/userRepoMock.go -package=mockRepository . UserRepository
type UserRepository interface {
	FindUserByUserID(ctx context.Context, userID uint) (user domain.User, err error)
	FindUserByEmail(ctx context.Context, email string) (user domain.User, err error)
	FindUserByUserName(ctx context.Context, userName string) (user domain.User, err error)
	FindUserByPhoneNumber(ctx context.Context, phoneNumber string) (user domain.User, err error)

	CheckOtherUserWithDetails(ctx context.Context, user domain.User) (domain.User, error) // find user exept this id
	SaveUser(ctx context.Context, user domain.User) (userID uint, err error)
	UpdateUser(ctx context.Context, user domain.User) (err error)

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
