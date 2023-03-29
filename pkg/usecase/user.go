package usecase

import (
	"context"
	"errors"

	"github.com/jinzhu/copier"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/res"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	service "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"golang.org/x/crypto/bcrypt"
)

type userUserCase struct {
	userRepo interfaces.UserRepository
}

func NewUserUseCase(repo interfaces.UserRepository) service.UserUseCase {
	return &userUserCase{userRepo: repo}
}

func (c *userUserCase) Login(ctx context.Context, user domain.User) (domain.User, error) {

	dbUser, dberr := c.userRepo.FindUser(ctx, user)

	// check user found or not
	if dberr != nil {
		return user, dberr
	} else if dbUser.ID == 0 {
		return user, errors.New("user not exist with this details")
	}

	// check user block_status user is blocked or not
	if dbUser.BlockStatus {
		return user, errors.New("user blocked by admin")
	}

	//check the user password with dbPassword
	if bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)) != nil {
		return user, errors.New("entered password is wrong")
	}

	return dbUser, nil
}

func (c *userUserCase) LoginOtp(ctx context.Context, user domain.User) (domain.User, error) {

	user, err := c.userRepo.FindUser(ctx, user)

	if err != nil {
		return user, errors.New("can't find the user")
	} else if user.ID == 0 {
		return user, errors.New("user not exist with this details")
	}

	// check user block_status user is blocked or not
	if user.BlockStatus {
		return user, errors.New("user blocked by admin")
	}

	return user, nil
}

func (c *userUserCase) Signup(ctx context.Context, user domain.User) error {
	// check the user already exist with this details
	checkUser, err := c.userRepo.FindUser(ctx, user)
	if err != nil {
		return err
	}
	// if user not exist then create user
	if checkUser.ID == 0 {
		//hash the password
		hashPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
		if err != nil {
			return errors.New("error to hash the password")
		}
		user.Password = string(hashPass)
		return c.userRepo.SaveUser(ctx, user)
	}
	// if user exist then check which field is exist
	return helper.CompareUsers(user, checkUser)
}

func (c *userUserCase) Account(ctx context.Context, userID uint) (domain.User, error) {

	var user = domain.User{ID: userID}
	return c.userRepo.FindUser(ctx, user)

}

func (c *userUserCase) EditAccount(ctx context.Context, user domain.User) error {

	// first check any other user exist with this entered unique fields
	checkUser, err := c.userRepo.FindUserExceptID(ctx, user)
	if err != nil {
		return err
	} else if checkUser.ID == 0 { // if there is no other user exist with this detail then update it
		return c.userRepo.EditUser(ctx, user)
	}

	// if any user exist with this field then show wich field is exis
	return helper.CompareUsers(user, checkUser)
}

func (c *userUserCase) SaveToCart(ctx context.Context, body req.ReqCart) error {

	// get the productitem to check product is valid
	productItem, err := c.userRepo.FindProductItem(ctx, body.ProductItemID)
	if err != nil {
		return err
	} else if productItem.ID == 0 {
		return errors.New("invalid product_item id")
	}

	// check productItem is out of stock or not
	if productItem.QtyInStock == 0 {
		return errors.New("product is now out of stock")
	}

	// then Find the cart of user
	cart, err := c.userRepo.FindCart(ctx, body.UserID)
	if err != nil {
		return err
	}

	cartItem := domain.CartItem{
		CartID:        cart.ID,
		ProductItemID: body.ProductItemID,
	}

	//check cart_item already exist or not
	if cartItem, err := c.userRepo.FindCartItem(ctx, cartItem); err != nil {
		return err
	} else if cartItem.ID != 0 {
		return errors.New("product_item aleady exist on cart")
	}

	// add productItem to cartItem
	if err := c.userRepo.SaveCartItem(ctx, cartItem); err != nil {
		return err
	}

	//update the cartTotal price

	return c.userRepo.UpdateCartPrice(ctx, cart)

}

func (c *userUserCase) RemoveCartItem(ctx context.Context, body req.ReqCart) error {

	// validate the product
	productItem, err := c.userRepo.FindProductItem(ctx, body.ProductItemID)

	if err != nil {
		return err
	} else if productItem.ID == 0 {
		return errors.New("invalid product_id")
	}

	// Find cart of user
	cart, err := c.userRepo.FindCart(ctx, body.UserID)
	if err != nil {
		return err
	} else if cart.TotalPrice == 0 {
		return errors.New("cart is emtpy can't remove items from cart")
	}

	cartItem := domain.CartItem{
		CartID:        cart.ID,
		ProductItemID: body.ProductItemID,
	}
	// find the cartItem of user
	cartItem, err = c.userRepo.FindCartItem(ctx, cartItem)
	if err != nil {
		return err
	} else if cartItem.ID == 0 {
		return errors.New("can't remove product_item from cart \nproduct_item not present in cart")
	}

	// then remvoe cart_item
	if err := c.userRepo.RemoveCartItem(ctx, cartItem); err != nil {
		return err
	}

	// update cart total price
	return c.userRepo.UpdateCartPrice(ctx, cart)
}

func (c *userUserCase) UpdateCartItem(ctx context.Context, body req.ReqCartCount) error {

	//validate the product
	productItem, err := c.userRepo.FindProductItem(ctx, body.ProductItemID)
	if err != nil {
		return err
	} else if productItem.ID == 0 {
		return errors.New("invalid product_item_id")
	}

	// then get the cart of user
	cart, err := c.userRepo.FindCart(ctx, body.UserID)
	if err != nil {
		return err
	} else if cart.ID == 0 {
		return errors.New("user cart is empty \ncan't updte cart_item count")
	}

	// check the given count is more than produc_item quantity
	if body.Count > productItem.QtyInStock {
		return errors.New("there is not this much quantity available in product_item")
	}

	// find the cartitem of user
	cartItem, err := c.userRepo.FindCartItem(ctx, domain.CartItem{
		CartID:        cart.ID,
		ProductItemID: body.ProductItemID,
	})

	if err != nil {
		return err
	} else if cartItem.ID == 0 {
		return errors.New("this cart_item not exist in the cart")
	}

	// check given quntity is 0
	if body.Count == 0 {
		return errors.New("can't change cart_item quntity to zero")
	}
	cartItem.Qty = body.Count
	// update cart_items
	err = c.userRepo.UpdateCartItem(ctx, cartItem)
	if err != nil {
		return err
	}

	// update the total price of cart
	return c.userRepo.UpdateCartPrice(ctx, cart)
}

func (c *userUserCase) GetCartItems(ctx context.Context, userId uint) (res.ResponseCart, error) {

	return c.userRepo.GetCartItems(ctx, userId)
}

// adddress
func (c *userUserCase) SaveAddress(ctx context.Context, address domain.Address, userID uint, isDefault bool) (domain.Address, error) {
	//check the address is already exist for the user
	address, err := c.userRepo.FindAddressByUserID(ctx, address, userID)
	if err != nil {
		return address, err
	} else if address.ID != 0 { // user have already this address exist
		return address, errors.New("user have already this address exist with same details")
	}

	//this address not exist then create it
	country, err := c.userRepo.FindCountryByID(ctx, address.CountryID)
	if err != nil {
		return address, err
	} else if country.ID == 0 {
		return address, errors.New("invalid country id")
	}

	// save the address on database
	address, err = c.userRepo.SaveAddress(ctx, address)
	if err != nil {
		return address, err
	}

	//creating a user address with this given value
	var userAdress = domain.UserAddress{
		UserID:    userID,
		AddressID: address.ID,
		IsDefault: isDefault,
	}

	// then update the address with user
	c.userRepo.SaveUserAddress(ctx, userAdress)

	return address, nil
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

	if c.userRepo.UpdateAddress(ctx, address) != nil {
		return err
	}

	//update the addres with user default or not with user
	if *addressBody.IsDefault {
		userAddress := domain.UserAddress{
			UserID:    userID,
			AddressID: address.ID,
			IsDefault: *addressBody.IsDefault,
		}
		return c.userRepo.UpdateUserAddress(ctx, userAddress)
	}

	return nil
}

// get all address
func (c *userUserCase) GetAddresses(ctx context.Context, userID uint) ([]res.ResAddress, error) {

	return c.userRepo.FindAllAddressByUserID(ctx, userID)
}

// to add new productItem to wishlist
func (c *userUserCase) AddToWishList(ctx context.Context, wishList domain.WishList) error {

	// first check the producItemID is valid or not
	productItem, err := c.userRepo.FindProductItem(ctx, wishList.ProductItemID)
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
	return c.userRepo.SaveWishListItem(ctx, wishList)
}

// remove from wishlist
func (c *userUserCase) RemoveFromWishList(ctx context.Context, wishList domain.WishList) error {

	// first check the producItemID is valid or not
	productItem, err := c.userRepo.FindProductItem(ctx, wishList.ProductItemID)
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
