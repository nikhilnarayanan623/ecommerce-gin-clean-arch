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
}

func NewUserUseCase(repo interfaces.UserRepository) service.UserUseCase {
	return &userUserCase{userRepo: repo}
}

// google login
func (c *userUserCase) GoogleLogin(ctx context.Context, user domain.User) (domain.User, error) {

	// // first check the user already exist or not if exist then direct login
	// checkUser, err := c.userRepo.FindUserByEmail(ctx, user.Email)
	// if err != nil {
	// 	return user, err
	// } else if checkUser.ID != 0 { // user already exist so direct login
	// 	return checkUser, nil
	// }

	// // user not exit so create a user

	// //creaet a ranodm username for user
	// user.UserName = utils.GenerateRandomUserName(user.FirstName)

	// userID, err := c.userRepo.SaveUserWithGoogleDetails(ctx, user)
	// if err != nil {
	// 	return user, err
	// }

	// user.ID = userID
	return user, nil
}

// login
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

		_, err = c.userRepo.SaveUser(ctx, user)
		if err != nil {
			return err
		}
		return nil
	}
	// if user exist then check which field is exist
	return utils.CompareUsers(user, checkUser)
}

func (c *userUserCase) Account(ctx context.Context, userID uint) (domain.User, error) {

	var user = domain.User{ID: userID}
	return c.userRepo.FindUser(ctx, user)

}

func (c *userUserCase) EditAccount(ctx context.Context, user domain.User) error {

	// first check any other user exist with this entered unique fields
	checkUser, err := c.userRepo.CheckOtherUserWithDetails(ctx, user)
	if err != nil {
		return err
	} else if checkUser.ID != 0 { // if there is an user exist with given details then make it as error
		err = utils.CompareUsers(user, checkUser)
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

// get user cart(it include total price and cartId)
func (c *userUserCase) GetUserCart(ctx context.Context, userID uint) (cart domain.Cart, err error) {

	cart, err = c.userRepo.FindCartByUserID(ctx, userID)
	if err != nil {
		return cart, err
	}

	return cart, err
}

func (c *userUserCase) SaveToCart(ctx context.Context, body req.ReqCart) (err error) {

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

	// find the cart of user
	cart, err := c.userRepo.FindCartByUserID(ctx, body.UserID)
	if err != nil {
		return err
	} else if cart.CartID == 0 { // if there is no cart is available for user then create it
		cart, err = c.userRepo.SaveCart(ctx, body.UserID)
		if err != nil {
			return err
		}
		log.Println(cart.CartID)
	}

	// check the given product item is already exit in user cart
	cartItem, err := c.userRepo.FindCartItemByCartAndProductItemID(ctx, cart.CartID, body.ProductItemID)
	if err != nil {
		return err
	} else if cartItem.CartItemID != 0 {
		return errors.New("product_item already exist on the cart can't save product to cart")
	}

	// add productItem to cartItem
	if err := c.userRepo.SaveCartItem(ctx, cart.CartID, body.ProductItemID); err != nil {
		return err
	}

	log.Printf("product_item_id %v saved on cart of  user_id %v", body.ProductItemID, body.UserID)
	return nil
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
	cart, err := c.userRepo.FindCartByUserID(ctx, body.UserID)
	if err != nil {
		return err
	} else if cart.CartID == 0 {
		return errors.New("can't remove product_item from user cart \n user cart is empty")
	}

	// check the product_item exist on user cart
	cartItem, err := c.userRepo.FindCartItemByCartAndProductItemID(ctx, cart.CartID, body.ProductItemID)
	if err != nil {
		return err
	} else if cartItem.CartItemID == 0 {
		return fmt.Errorf("prduct_item with id %v is not exist on user cart", body.ProductItemID)
	}
	// then remvoe cart_item
	err = c.userRepo.DeleteCartItem(ctx, cartItem.CartItemID)
	if err != nil {
		return err
	}

	log.Printf("product_item with id %v removed form user cart with usre id %v", body.ProductItemID, body.UserID)
	return nil
}

func (c *userUserCase) UpdateCartItem(ctx context.Context, body req.ReqCartCount) error {

	//check the given product_item_id is valid or not
	productItem, err := c.userRepo.FindProductItem(ctx, body.ProductItemID)
	if err != nil {
		return err
	} else if productItem.ID == 0 {
		return errors.New("invalid product_item_id")
	}

	if body.Count < 1 {
		return fmt.Errorf("can't change cart_item qty to %v \n minuimun qty is 1", body.Count)
	}

	// check the given qty of product_item updation is morethan prodct_item qty
	if body.Count > productItem.QtyInStock {
		return errors.New("there is not this much quantity available in product_item")
	}

	// find the cart of user
	cart, err := c.userRepo.FindCartByUserID(ctx, body.UserID)
	if err != nil {
		return err
	} else if cart.CartID == 0 {
		return errors.New("user cart is empty")
	}

	// find the cart_item with given product_id and user cart_id  and check the product_item present in cart or no
	cartItem, err := c.userRepo.FindCartItemByCartAndProductItemID(ctx, cart.CartID, body.ProductItemID)
	if err != nil {
		return err
	} else if cartItem.CartItemID == 0 {
		return fmt.Errorf("product_item not exist in the cart with given product_item_id %v", body.ProductItemID)
	}

	// update the cart_item qty
	if err := c.userRepo.UpdateCartItemQty(ctx, cartItem.CartItemID, body.Count); err != nil {
		return err
	}

	log.Printf("updated user carts product_items qty of product_item_id %v , qty %v", body.ProductItemID, body.Count)
	return nil
}

func (c *userUserCase) GetUserCartItems(ctx context.Context, cartId uint) (cartItems []res.ResCartItem, err error) {
	// get the cart_items of user
	cartItems, err = c.userRepo.FindAllCartItemsByCartID(ctx, cartId)
	if err != nil {
		return cartItems, err
	}

	log.Printf("sucessfully got all cart_items of user wtih cart_id %v", cartId)
	return cartItems, nil
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
	if err := c.userRepo.SaveWishListItem(ctx, wishList); err != nil {
		return err
	}

	return nil
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
