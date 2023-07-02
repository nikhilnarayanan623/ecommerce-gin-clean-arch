package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jinzhu/copier"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/response"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	service "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type userUserCase struct {
	userRepo    interfaces.UserRepository
	cartRepo    interfaces.CartRepository
	productRepo interfaces.ProductRepository
}

func NewUserUseCase(userRepo interfaces.UserRepository, cartRepo interfaces.CartRepository,
	productRepo interfaces.ProductRepository) service.UserUseCase {
	return &userUserCase{
		userRepo:    userRepo,
		cartRepo:    cartRepo,
		productRepo: productRepo,
	}
}

func (c *userUserCase) FindProfile(ctx context.Context, userID uint) (domain.User, error) {

	user, err := c.userRepo.FindUserByUserID(ctx, userID)
	if err != nil {
		return domain.User{}, utils.PrependMessageToError(err, "failed to find user details")
	}

	return user, nil
}

func (c *userUserCase) UpdateProfile(ctx context.Context, user domain.User) error {

	// first check any other user exist with this entered unique fields
	checkUser, err := c.userRepo.FindUserByUserNameEmailOrPhoneNotID(ctx, user)
	if err != nil {
		return err
	}
	if checkUser.ID != 0 { // if there is an user exist with given details then make it as error
		err = utils.CompareUserExistingDetails(user, checkUser)
		fmt.Println(user)
		return err
	}

	// if user password given then hash the password
	if user.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
		if err != nil {
			return fmt.Errorf("failed to generate hash password for user")
		}
		user.Password = string(hash)
	}

	err = c.userRepo.UpdateUser(ctx, user)

	if err != nil {
		return err
	}

	return nil
}

// adddress
func (c *userUserCase) SaveAddress(ctx context.Context, userID uint, address domain.Address, isDefault bool) error {

	exist, err := c.userRepo.IsAddressAlreadyExistForUser(ctx, address, userID)
	if err != nil {
		return fmt.Errorf("failed to check address already exist \nerror:%v", err.Error())
	}
	if exist {
		return fmt.Errorf("given address already exist for user")
	}

	// //this address not exist then create it
	// country, err := c.userRepo.FindCountryByID(ctx, address.CountryID)
	// if err != nil {
	// 	return err
	// } else if country.ID == 0 {
	// 	return errors.New("invalid country id")
	// }

	// save the address on database
	addressID, err := c.userRepo.SaveAddress(ctx, address)
	if err != nil {
		return err
	}

	//creating a user address with this given value
	var userAddress = domain.UserAddress{
		UserID:    userID,
		AddressID: addressID,
		IsDefault: isDefault,
	}

	// then update the address with user
	err = c.userRepo.SaveUserAddress(ctx, userAddress)

	if err != nil {
		return err
	}

	return nil
}

func (c *userUserCase) UpdateAddress(ctx context.Context, addressBody request.EditAddress, userID uint) error {

	if exist, err := c.userRepo.IsAddressIDExist(ctx, addressBody.ID); err != nil {
		return err
	} else if !exist {
		return errors.New("invalid address id")
	}

	var address domain.Address
	copier.Copy(&address, &addressBody)

	if err := c.userRepo.UpdateAddress(ctx, address); err != nil {
		return err
	}

	// check the user address need to set default or not if it need then set it as default
	if addressBody.IsDefault != nil && *addressBody.IsDefault {
		userAddress := domain.UserAddress{
			UserID:    userID,
			AddressID: address.ID,
			IsDefault: *addressBody.IsDefault,
		}

		err := c.userRepo.UpdateUserAddress(ctx, userAddress)
		if err != nil {
			return err
		}
	}
	log.Printf("successfully address saved for user with user_id %v", userID)
	return nil
}

// get all address
func (c *userUserCase) FindAddresses(ctx context.Context, userID uint) ([]response.Address, error) {

	return c.userRepo.FindAllAddressByUserID(ctx, userID)
}

// to add new productItem to wishlist
func (c *userUserCase) SaveToWishList(ctx context.Context, wishList domain.WishList) error {

	// check the productItem already exist on wishlist for user
	checkWishList, err := c.userRepo.FindWishListItem(ctx, wishList.ProductItemID, wishList.UserID)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to check product item already exist on wish list")
	}
	if checkWishList.ID != 0 {
		return ErrExistWishListProductItem
	}

	err = c.userRepo.SaveWishListItem(ctx, wishList)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to save product item on wish list")
	}

	return nil
}

// remove from wishlist
func (c *userUserCase) RemoveFromWishList(ctx context.Context, userID, productItemID uint) error {

	err := c.userRepo.RemoveWishListItem(ctx, userID, productItemID)
	if err != nil {
		return utils.PrependMessageToError(err, "failed to remove product item form wish list")
	}

	return nil
}

func (c *userUserCase) FindAllWishListItems(ctx context.Context, userID uint) ([]response.WishListItem, error) {

	wishListItems, err := c.userRepo.FindAllWishListItemsByUserID(ctx, userID)
	if err != nil {
		return nil, utils.PrependMessageToError(err, "failed to find wish list product items")
	}

	for i, productItem := range wishListItems {
		variationValues, err := c.productRepo.FindAllVariationValuesOfProductItem(ctx, productItem.ProductItemID)
		if err != nil {
			return nil, utils.PrependMessageToError(err, "failed to find variation values product item")
		}
		wishListItems[i].VariationValues = variationValues
	}

	return wishListItems, nil
}
