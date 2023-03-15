package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	"gorm.io/gorm"
)

type userDatabse struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) interfaces.UserRepository {
	return &userDatabse{DB: DB}
}

func (c *userDatabse) FindUser(ctx context.Context, user domain.Users) (domain.Users, error) {
	fmt.Println("user", user)

	// check id,email,phone any of then match i db
	err := c.DB.Raw("SELECT * FROM users where id=? OR email=? OR phone=?", user.ID, user.Email, user.Phone).Scan(&user).Error

	if err != nil {
		return user, errors.New("faild to get user")
	}

	return user, nil
}

func (c *userDatabse) SaveUser(ctx context.Context, user domain.Users) (domain.Users, error) {

	// check whether user is already exisist
	c.DB.Raw("SELECT * FROM users WHERE email=? OR phone=?", user.Email, user.Phone).Scan(&user)
	//if exist then return message as user exist
	if user.ID != 0 {
		return user, errors.New("user already exist with this details")
	}

	//save the user details
	err := c.DB.Save(&user).Error

	return user, err
}

func (c *userDatabse) GetCartItems(ctx context.Context, userId uint) (helper.ResponseCart, error) {

	// var (
	// 	user = domain.Users{ID: userId}
	// 	// resCart = helper.ResCart{CartItems: make([]helper.ResCartItem, 0)}
	// 	resCart helper.ResCart
	// 	cart    domain.Cart
	// 	//cartItems []domain.CartItem
	// )

	// //first find the user
	// user, err := c.FindUser(ctx, user)

	// if err != nil {
	// 	return resCart, err
	// }

	// // then get cart id of user
	// if c.DB.Raw("SELECT * FROM carts WHERE users_id=?", userId).Scan(&cart); cart.ID == 0 {
	// 	return resCart, helper.SingleRespStruct{Error: "User Have no cart"} // I think I want to delete it later
	// }

	// // add total price to response
	// resCart.TotalPrice = cart.TotalPrice

	// //then get all cart items of user
	// var resCartItem []helper.ResCartItem

	// c.DB.Raw("SELECT * FROM cart_items WHERE cart_id=?", cart.ID).Scan(&resCartItem)

	// // assign it to resCart
	// resCart.CartItems = resCartItem

	// return resCart, nil

	// find the cartId of user using userId

	var cart domain.Cart
	if c.DB.Raw("SELECT * FROM carts WHERE users_id=?", userId).Scan(&cart).Error != nil {
		return helper.ResponseCart{}, errors.New("faild to get the cart of user")
	}

	return helper.ResponseCart{}, nil
}
