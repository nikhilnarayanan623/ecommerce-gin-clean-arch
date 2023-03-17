package usecase

import (
	"context"
	"errors"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper"
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

func (c *userUserCase) Signup(ctx context.Context, user domain.User) (domain.User, error) {
	//hash the password
	hashPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	if err != nil {
		return user, errors.New("error to hash the password")
	}
	user.Password = string(hashPass)

	return c.userRepo.SaveUser(ctx, user)
}

func (c *userUserCase) Home(ctx context.Context, userId uint) (domain.User, error) {

	var user = domain.User{ID: userId}

	return c.userRepo.FindUser(ctx, user)

}

func (c *userUserCase) SaveToCart(ctx context.Context, body helper.ReqCart) (domain.CartItem, error) {

	// get the productitem to check product is valid
	productItem, err := c.userRepo.FindProductItem(ctx, body.ProductItemID)
	if err != nil {
		return domain.CartItem{}, err
	} else if productItem.ID == 0 {
		return domain.CartItem{}, errors.New("invalid product_item id")
	}

	// check productItem is out of stock or not
	if productItem.QtyInStock == 0 {
		return domain.CartItem{}, errors.New("product is now out of stock")
	}

	// then Find the cart of user
	cart, err := c.userRepo.FindCart(ctx, body.UserID)
	if err != nil {
		return domain.CartItem{}, err
	}

	// add productItem to cartItem
	cartItem, err := c.userRepo.SaveCartItem(ctx, cart.ID, productItem.ID)
	if err != nil {
		return domain.CartItem{}, err
	}

	//update the cartTotal price
	cart, err = c.userRepo.UpdateCartPrice(ctx, cart)

	if err != nil {
		return cartItem, err
	}

	return cartItem, nil

}

func (c *userUserCase) RemoveCartItem(ctx context.Context, body helper.ReqCart) (domain.Cart, error) {

	// validate the product
	productItem, err := c.userRepo.FindProductItem(ctx, body.ProductItemID)

	if err != nil {
		return domain.Cart{}, err
	} else if productItem.ID == 0 {
		return domain.Cart{}, errors.New("invalid product_id")
	}

	// Find cart of user
	cart, err := c.userRepo.FindCart(ctx, body.UserID)
	if err != nil {
		return domain.Cart{}, err
	} else if cart.TotalPrice == 0 {
		return domain.Cart{}, errors.New("nothing to remove form cart")
	}

	// find the cartItem of user
	cartItem, err := c.userRepo.FindCartItem(ctx, cart.ID, productItem.ID)
	if err != nil {
		return cart, err
	} else if cartItem.ID == 0 {
		return cart, errors.New("this product is not in the cart")
	}

	// delete the cart Item
	cartItem, err = c.userRepo.RemoveCartItem(ctx, cartItem)
	if err != nil {
		return cart, err
	}

	//update the total price of cart
	return c.userRepo.UpdateCartPrice(ctx, cart)
}

func (c *userUserCase) UpdateCartItem(ctx context.Context, body helper.ReqCartCount) (domain.CartItem, error) {

	//validate the product
	productItem, err := c.userRepo.FindProductItem(ctx, body.ProductItemID)
	if err != nil {
		return domain.CartItem{}, err
	} else if productItem.ID == 0 {
		return domain.CartItem{}, errors.New("invalid product_item_id")
	}

	// then get the cart of user
	cart, err := c.userRepo.FindCart(ctx, body.UserID)
	if err != nil {
		return domain.CartItem{}, err
	} else if cart.ID == 0 {
		return domain.CartItem{}, errors.New("there is no cart for the user")
	}

	// get the cartitem of user
	cartItem, err := c.userRepo.FindCartItem(ctx, cart.ID, productItem.ID)
	if err != nil {
		return domain.CartItem{}, err
	} else if cartItem.ID == 0 {
		return domain.CartItem{}, errors.New("this productItem not exist in the cart")
	}

	// update the cartItem quantity
	//check the product count need to increment or not
	if *body.Increment { // to increment
		cartItem.Qty += 1
	} else if cartItem.Qty == 1 { // decremtnet last product quantity
		return cartItem, errors.New("can't decrement last count of productItem")
	} else { // decrement quantity
		cartItem.Qty -= 1
	}

	cartItem, err = c.userRepo.UpdateCartItem(ctx, cartItem)
	if err != nil {
		return cartItem, err
	}

	// update the total price of cart
	if _, err := c.userRepo.UpdateCartPrice(ctx, cart); err != nil {
		return cartItem, err
	}

	return cartItem, nil
}

func (c *userUserCase) GetCartItems(ctx context.Context, userId uint) (helper.ResponseCart, error) {

	return c.userRepo.GetCartItems(ctx, userId)
}
