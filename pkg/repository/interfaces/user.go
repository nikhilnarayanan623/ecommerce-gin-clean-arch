package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/response"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
)

//go:generate mockgen -destination=../../mock/mockrepo/user_mock.go -package=mockrepo . UserRepository
type UserRepository interface {
	FindUserByUserID(ctx context.Context, userID uint) (user domain.User, err error)
	FindUserByEmail(ctx context.Context, email string) (user domain.User, err error)
	FindUserByUserName(ctx context.Context, userName string) (user domain.User, err error)
	FindUserByPhoneNumber(ctx context.Context, phoneNumber string) (user domain.User, err error)
	FindUserByUserNameEmailOrPhoneNotID(ctx context.Context, user domain.User) (domain.User, error)

	SaveUser(ctx context.Context, user domain.User) (userID uint, err error)
	UpdateVerified(ctx context.Context, userID uint) error
	UpdateUser(ctx context.Context, user domain.User) (err error)
	UpdateBlockStatus(ctx context.Context, userID uint, blockStatus bool) error

	//address
	FindCountryByID(ctx context.Context, countryID uint) (domain.Country, error)
	FindAddressByID(ctx context.Context, addressID uint) (response.Address, error)
	IsAddressIDExist(ctx context.Context, addressID uint) (exist bool, err error)
	IsAddressAlreadyExistForUser(ctx context.Context, address domain.Address, userID uint) (bool, error)
	FindAllAddressByUserID(ctx context.Context, userID uint) ([]response.Address, error)
	SaveAddress(ctx context.Context, address domain.Address) (addressID uint, err error)
	UpdateAddress(ctx context.Context, address domain.Address) error
	// address join table
	SaveUserAddress(ctx context.Context, userAdress domain.UserAddress) error
	UpdateUserAddress(ctx context.Context, userAddress domain.UserAddress) error

	//wishlist
	FindWishListItem(ctx context.Context, productID, userID uint) (domain.WishList, error)
	FindAllWishListItemsByUserID(ctx context.Context, userID uint) ([]response.WishListItem, error)
	SaveWishListItem(ctx context.Context, wishList domain.WishList) error
	RemoveWishListItem(ctx context.Context, userID, productItemID uint) error
}
