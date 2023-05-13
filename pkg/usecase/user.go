package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jinzhu/copier"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	service "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
	"golang.org/x/crypto/bcrypt"
)

type userUserCase struct {
	userRepo interfaces.UserRepository
	cartRepo interfaces.CartRepository
}

func NewUserUseCase(userRepo interfaces.UserRepository, cartRepo interfaces.CartRepository) service.UserUseCase {
	return &userUserCase{
		userRepo: userRepo,
		cartRepo: cartRepo,
	}
}

func (c *userUserCase) Account(ctx context.Context, userID uint) (domain.User, error) {

	var user = domain.User{ID: userID}
	// user, err := c.userRepo.FindUser(ctx, user)

	return user, nil

}

func (c *userUserCase) EditAccount(ctx context.Context, user domain.User) error {

	// first check any other user exist with this entered unique fields
	checkUser, err := c.userRepo.FindUserByUserNameEmailOrPhoneNotID(ctx, user)
	if err != nil {
		return err
	} else if checkUser.ID != 0 { // if there is an user exist with given details then make it as error
		err = utils.CompareUserExistingDetails(user, checkUser)
		return err
	}

	// if user password given then hash the password
	if user.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
		if err != nil {
			return fmt.Errorf("faild to generate hash password for user")
		}
		user.Password = string(hash)
	}

	err = c.userRepo.UpdateUser(ctx, user)

	if err != nil {
		return err
	}

	log.Printf("successfully user details updaed for user with user_id %v", user.ID)
	return nil
}

// adddress
func (c *userUserCase) SaveAddress(ctx context.Context, userID uint, address domain.Address, isDefault bool) error {
	//check the address is already exist for the user
	address, err := c.userRepo.FindAddressByUserID(ctx, address, userID)
	if err != nil {
		return err
	} else if address.ID != 0 { // user have already this address exist
		return errors.New("user have already this address exist with same details")
	}

	//this address not exist then create it
	country, err := c.userRepo.FindCountryByID(ctx, address.CountryID)
	if err != nil {
		return err
	} else if country.ID == 0 {
		return errors.New("invalid country id")
	}

	// save the address on database
	addressID, err := c.userRepo.SaveAddress(ctx, address)
	if err != nil {
		return err
	}

	//creating a user address with this given value
	var userAdress = domain.UserAddress{
		UserID:    userID,
		AddressID: addressID,
		IsDefault: isDefault,
	}

	// then update the address with user
	err = c.userRepo.SaveUserAddress(ctx, userAdress)

	if err != nil {
		return err
	}
	log.Printf("successfully user address stored for user with user_id %v", userID)
	return nil
}

func (c *userUserCase) EditAddress(ctx context.Context, addressBody req.ReqEditAddress, userID uint) error {

	// first validate the addessId is valid or not
	address, err := c.userRepo.FindAddressByID(ctx, addressBody.ID)
	if err != nil {
		return err
	} else if address.ID == 0 {
		return errors.New("invalid address id")
	}

	// validate the country id
	country, err := c.userRepo.FindCountryByID(ctx, addressBody.CountryID)
	if err != nil {
		return err
	} else if country.ID == 0 {
		return errors.New("invalid country id")
	}
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
func (c *userUserCase) GetAddresses(ctx context.Context, userID uint) ([]res.ResAddress, error) {

	return c.userRepo.FindAllAddressByUserID(ctx, userID)
}

// to add new productItem to wishlist
func (c *userUserCase) AddToWishList(ctx context.Context, wishList domain.WishList) error {

	// first check the producItemID is valid or not
	//productItem, err := c.userRepo.FindProductItem(ctx, wishList.ProductItemID)
	var (
		productItem domain.ProductItem
		err         error
	)
	if err != nil {
		return err
	} else if productItem.ID == 0 {
		return errors.New("invalid product_id")
	}

	// check the productItem already exist on wishlist for user
	checkWishList, err := c.userRepo.FindWishListItem(ctx, wishList.ProductItemID, wishList.UserID)
	if err != nil {
		return err
	} else if checkWishList.ID != 0 {
		return errors.New("productItem already exist on wishlist")
	}

	// save productItem wishlist
	if err := c.userRepo.SaveWishListItem(ctx, wishList); err != nil {
		return err
	}

	return nil
}

// remove from wishlist
func (c *userUserCase) RemoveFromWishList(ctx context.Context, wishList domain.WishList) error {

	// first check the producItemID is valid or not
	//productItem, err := c.userRepo.FindProductItem(ctx, wishList.ProductItemID)
	var (
		productItem domain.ProductItem
		err         error
	)
	if err != nil {
		return err
	} else if productItem.ID == 0 {
		return errors.New("invalid product_id")
	}

	// check the productItem already exist on wishlist for user
	wishList, err = c.userRepo.FindWishListItem(ctx, wishList.ProductItemID, wishList.UserID)
	if err != nil {
		return err
	} else if wishList.ID == 0 {
		return errors.New("productItem not exist exist in wishlist")
	}

	// remove the productItem from user wihsList
	return c.userRepo.RemoveWishListItem(ctx, wishList)
}

func (c *userUserCase) GetWishListItems(ctx context.Context, userID uint) ([]res.ResWishList, error) {
	return c.userRepo.FindAllWishListItemsByUserID(ctx, userID)
}
