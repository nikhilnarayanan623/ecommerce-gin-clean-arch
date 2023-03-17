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

// to get productItem id
func (c *userDatabse) FindProductItem(ctx context.Context, productItemID uint) (domain.ProductItem, error) {

	var productItem domain.ProductItem
	if c.DB.Raw("SELECT * FROM product_items WHERE id=?", productItemID).Scan(&productItem).Error != nil {
		return domain.ProductItem{}, errors.New("faild to get productItem from database")
	}
	return productItem, nil
}

// get the cart of user
func (c *userDatabse) FindCart(ctx context.Context, userID uint) (domain.Cart, error) {
	// then check user have cart already exist or not
	var cart domain.Cart
	if c.DB.Raw("SELECT * FROM carts WHERE user_id=?", userID).Scan(&cart).Error != nil {
		return cart, errors.New("faild to get cart of user from database")
	}

	// if user have no cart then create a new cart for user
	if cart.ID == 0 {
		querry := `INSERT INTO carts (user_id,total_price) VALUES ($1,$2) RETURNING id,user_id,total_price`
		if c.DB.Raw(querry, userID, 0).Scan(&cart).Error != nil {
			return cart, errors.New("faild to create cart for user in database")
		}
	}

	return cart, nil
}

func (c *userDatabse) UpdateCartPrice(ctx context.Context, cart domain.Cart) (domain.Cart, error) {

	// update cartTotal Price
	query := ` SELECT SUM(ci.qty * pi.price) FROM cart_items ci 
		JOIN carts c ON ci.cart_id=c.id JOIN product_items pi 
		ON ci.product_item_id=pi.id AND c.id=? GROUP BY cart_id;`

	var TotalPrice uint

	if c.DB.Raw(query, cart.ID).Scan(&TotalPrice).Error != nil {
		return cart, errors.New("faild to calculate total price of cartItems")
	}
	fmt.Println(TotalPrice, "total price")

	//update the total price on cart
	if c.DB.Raw("UPDATE carts SET total_price = ? WHERE id=? RETURNING total_price", TotalPrice, cart.ID).Scan(&cart).Error != nil {
		return cart, errors.New("faild to update the total price of cart")
	}

	return cart, nil
}

// find a cartItem
func (c *userDatabse) FindCartItem(ctx context.Context, cartID, productItemID uint) (domain.CartItem, error) {
	// check the cartitem exist for user
	var cartItem domain.CartItem
	if c.DB.Raw("SELECT * FROM cart_items WHERE cart_id=? AND product_item_id=?", cartID, productItemID).Scan(&cartItem).Error != nil {
		return cartItem, errors.New("faild to get cartItem of user")
	}
	return cartItem, nil
}

// add a productItem to cartitem
func (c *userDatabse) SaveCartItem(ctx context.Context, cartID, productItemID uint) (domain.CartItem, error) {
	var cartItems domain.CartItem
	if c.DB.Raw("SELECT * FROM cart_items WHERE cart_id=? AND product_item_id=?", cartID, productItemID).Scan(&cartItems).Error != nil {
		return cartItems, errors.New("faild to get cartItem")
	}

	if cartItems.ID == 0 {
		querry := `INSERT INTO cart_items (cart_id,product_item_id,qty) VALUES ($1,$2,$3) RETURNING id,cart_id,product_item_id`
		if c.DB.Raw(querry, cartID, productItemID, 1).Scan(&cartItems).Error != nil {
			return cartItems, errors.New("faild to insert cartItems to cart")
		}
	} else {
		querry := `UPDATE cart_items SET qty=? WHERE product_item_id=?`
		if c.DB.Raw(querry, cartItems.Qty+1, productItemID).Scan(&cartItems).Error != nil {
			return cartItems, errors.New("faild to update cart_item")
		}
	}

	return cartItems, nil
}

func (c *userDatabse) RemoveCartItem(ctx context.Context, cartItem domain.CartItem) (domain.CartItem, error) {

	// delete productItem from cart
	if c.DB.Raw("DELETE FROM cart_items WHERE id=?", cartItem.ID).Scan(&cartItem).Error != nil {
		return cartItem, errors.New("faild to delete cart_item from database")
	}

	return cartItem, nil
}

func (c *userDatabse) UpdateCartItem(ctx context.Context, cartItem domain.CartItem) (domain.CartItem, error) {

	if c.DB.Raw("UPDATE cart_items SET qty = ? WHERE id=? RETURNING qty", cartItem.Qty, cartItem.ID).Scan(&cartItem).Error != nil {
		return cartItem, errors.New("faild to update the qty of cartItem")
	}
	return cartItem, nil
}

// get all itmes from cart
func (c *userDatabse) GetCartItems(ctx context.Context, userId uint) (helper.ResponseCart, error) {

	var (
		response helper.ResponseCart
	)
	// get the cart of user
	cart, err := c.FindCart(ctx, userId)
	if err != nil {
		return response, err
	}

	// get the cartItem of all user with subtotal
	query := `SELECT ci.product_item_id,p.product_name, ci.qty,pi.price,pi.price * ci.qty AS sub_total, (CASE WHEN pi.qty_in_stock=0 THEN 'T' ELSE 'F' END) AS out_of_stock  
				FROM cart_items ci JOIN product_items pi ON ci.product_item_id = pi.id 
				JOIN products p ON pi.product_id=p.id AND ci.cart_id=?`

	if c.DB.Raw(query, cart.ID).Scan(&response.CartItems).Error != nil {
		return response, errors.New("faild to get cartItems from database")
	}

	response.TotalPrice = cart.TotalPrice

	return response, nil
}
