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

func (c *userDatabse) FindUser(ctx context.Context, user domain.User) (domain.User, error) {
	fmt.Println("user", user)

	// check id,email,phone any of then match i db
	err := c.DB.Raw("SELECT * FROM users where id=? OR email=? OR phone=?", user.ID, user.Email, user.Phone).Scan(&user).Error

	if err != nil {
		return user, errors.New("faild to get user")
	}

	return user, nil
}

func (c *userDatabse) SaveUser(ctx context.Context, user domain.User) (domain.User, error) {

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

// to add a productItem to cart
func (c *userDatabse) AddToCart(ctx context.Context, body helper.ReqCart) (domain.Cart, error) {

	// first check the given productItem is valid or not
	var productItem domain.ProductItem
	if c.DB.Raw("SELECT * FROM product_items WHERE id=?", body.ProductItemID).Scan(&productItem).Error != nil {
		return domain.Cart{}, errors.New("faild to get productItem from database")
	} else if productItem.ID == 0 {
		return domain.Cart{}, errors.New("invalid product_item id")
	}

	// then check user have cart already exist or not
	var cart domain.Cart
	if c.DB.Raw("SELECT * FROM carts WHERE user_id=?", body.UserID).Scan(&cart).Error != nil {
		return cart, errors.New("faild to get cart of user from database")
	}

	// if user have no cart then create a new cart for user
	if cart.ID == 0 {
		querry := `INSERT INTO carts (user_id,total_price) VALUES ($1,$2) RETURNING id,user_id,total_price`
		if c.DB.Raw(querry, body.UserID, 0).Scan(&cart).Error != nil {
			return cart, errors.New("faild to create cart for user in database")
		}
	}

	// last add productId on cartItems
	// check cartItem already exist with this cartId
	var cartItems domain.CartItem
	if c.DB.Raw("SELECT * FROM cart_items WHERE cart_id=? AND product_item_id=?", cart.ID, productItem.ID).Scan(&cartItems).Error != nil {
		return cart, errors.New("faild to get cartItem")
	}

	if cartItems.ID == 0 {
		querry := `INSERT INTO cart_items (cart_id,product_item_id,qty) VALUES ($1,$2,$3) RETURNING id,cart_id,product_item_id`
		c.DB.Raw(querry, cart.ID, productItem.ID, 1).Scan(&cartItems)
	} else {
		querry := `UPDATE cart_items SET qty=? WHERE product_item_id=?`
		c.DB.Raw(querry, cartItems.Qty+1, productItem.ID).Scan(&cartItems)
	}

	// at last update total price of userCart
	totalPrice := cart.TotalPrice + productItem.Price

	if c.DB.Raw("UPDATE carts SET total_price=? WHERE id=? RETURNING id,user_id,total_price", totalPrice, cart.ID).Scan(&cart).Error != nil {
		return cart, errors.New("faild to update cart total price")
	}
	return cart, nil
}

// get all itmes from cart
func (c *userDatabse) GetCartItems(ctx context.Context, userId uint) (helper.ResponseCart, error) {

	var (
		cart     domain.Cart
		response helper.ResponseCart
	)
	// get the cart of user
	if c.DB.Raw("SELECT * FROM carts WHERE user_id=?", userId).Scan(&cart).Error != nil {
		return response, errors.New("faild to get the cart of user")
	}

	// get the cartItem of all user with subtotal
	query := `SELECT ci.product_item_id,p.product_name, ci.qty,pi.price,pi.price * ci.qty AS sub_total, (CASE WHEN pi.qty_in_stock=0 THEN 'T' ELSE 'F' END) AS out_of_stock  
				FROM cart_items ci JOIN product_items pi ON ci.product_item_id = pi.id 
				JOIN products p ON pi.product_id=p.id AND ci.cart_id=?`

	if c.DB.Raw(query, cart.ID).Scan(&response.CartItems).Error != nil {
		return response, errors.New("faild to get cartItems from database")
	}

	// atlast add the total price into response
	response.TotalPrice = cart.TotalPrice

	return response, nil
}
