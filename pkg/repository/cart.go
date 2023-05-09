package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
	"gorm.io/gorm"
)

type cartDatabase struct {
	DB *gorm.DB
}

func NewCartRepository(db *gorm.DB) interfaces.CartRepository {
	return &cartDatabase{
		DB: db,
	}
}

// find a cartItem
func (c *cartDatabase) FindCartByUserID(ctx context.Context, userID uint) (cart domain.Cart, err error) {

	query := `SELECT * FROM carts WHERE user_id = ?`
	err = c.DB.Raw(query, userID).Scan(&cart).Error

	return
}

// save cart for user
func (c *cartDatabase) SaveCart(ctx context.Context, userID uint) (cartID uint, err error) {

	query := `INSERT INTO carts (user_id,total_price) VALUES($1, $2) RETURNING cart_id`
	err = c.DB.Raw(query, userID, 0).Scan(&cartID).Error

	return cartID, err
}

// find cart_items
func (c *cartDatabase) FindCartItemByID(ctx context.Context, cartItemID uint) (cartItem domain.CartItem, err error) {
	query := `SELECT * FROM cart_items WHERE cart_item_id = ?`
	if c.DB.Raw(query, cartItemID).Scan(&cartItem).Error != nil {
		return cartItem, errors.New("faild to find cart_item with cart_item_id")
	}

	return cartItem, nil
}

// find cart_item by cart_id and product_item id (can use for checking proudct already exit )
func (c *cartDatabase) FindCartItemByCartAndProductItemID(ctx context.Context, cartID, productItemID uint) (cartItem domain.CartItem, err error) {
	query := `SELECT * FROM cart_items WHERE cart_id = $1 AND product_item_id = $2`
	if c.DB.Raw(query, cartID, productItemID).Scan(&cartItem).Error != nil {
		return cartItem, errors.New("faild to find cart_item with given cart_id and product_item_id")
	}
	return cartItem, nil
}

// add a productItem to cartitem
func (c *cartDatabase) SaveCartItem(ctx context.Context, cartId, productItemId uint) error {

	querry := `INSERT INTO cart_items (cart_id, product_item_id, qty) VALUES ($1, $2, $3)`
	if c.DB.Exec(querry, cartId, productItemId, 1).Error != nil {
		return errors.New("faild to save cart_items")
	}

	return nil
}

func (c *cartDatabase) DeleteCartItem(ctx context.Context, cartItemID uint) error {
	// delete productItem from cart
	query := `DELETE FROM cart_items WHERE cart_item_id = $1`
	if c.DB.Exec(query, cartItemID).Error != nil {
		return errors.New("faild to remove product_items from cart")
	}

	return nil
}
func (c *cartDatabase) DeleteAllCartItemsByUserID(ctx context.Context, userID uint) error {

	query := ` DELETE FROM cart_items ci USING carts c WHERE c.cart_id = ci.cart_id AND c.user_id = $1`
	err := c.DB.Exec(query, userID).Error

	return err
}

func (c *cartDatabase) DeleteAllCartItemsByCartID(ctx context.Context, cartID uint) error {

	query := ` DELETE FROM cart_items WHERE cart_id = $1`
	err := c.DB.Exec(query, cartID).Error
	return err
}

func (c *cartDatabase) UpdateCartItemQty(ctx context.Context, cartItemId, qty uint) error {

	query := `UPDATE cart_items SET qty = $1 WHERE cart_item_id = $2`
	if c.DB.Exec(query, qty, cartItemId).Error != nil {
		return errors.New("faild to update the qty of cart_item")
	}
	return nil
}

// get all itmes from cart
func (c *cartDatabase) FindAllCartItemsByCartID(ctx context.Context, cartID uint) (cartItems []res.ResCartItem, err error) {

	// get the cartItem of all user with subtotal
	query := `SELECT ci.product_item_id, p.product_name, ci.qty,pi.price ,
	 pi.discount_price, CASE WHEN pi.discount_price > 0 THEN pi.discount_price * ci.qty ELSE pi.price * ci.qty END AS sub_total,  
	 pi.qty_in_stock 
	 FROM cart_items ci INNER JOIN product_items pi ON ci.product_item_id = pi.id 
	 INNER JOIN products p ON pi.product_id = p.id AND ci.cart_id=?`

	if c.DB.Raw(query, cartID).Scan(&cartItems).Error != nil {
		return cartItems, fmt.Errorf("faild to get cart_items from cart with cart_id %v", cartID)
	}

	return cartItems, err
}

func (c *cartDatabase) CheckcartIsValidForOrder(ctx context.Context, userID uint) (cart domain.Cart, err error) {

	cart, err = c.FindCartByUserID(ctx, userID)
	if err != nil {
		return cart, err
	}

	// check any of the product is out of stock in cart
	var idOfOneOutOfStockProduct uint
	query := `SELECT DISTINCT pi.id FROM product_items pi 
	INNER JOIN cart_items ci ON pi.id = ci.product_item_id 
	INNER JOIN carts c ON ci.cart_id = c.cart_id 
	WHERE c.user_id = 2 AND pi.qty_in_stock <= 0`

	err = c.DB.Raw(query).Scan(&idOfOneOutOfStockProduct).Error
	if err != nil {
		return cart, err
	}

	if idOfOneOutOfStockProduct != 0 {
		return cart, fmt.Errorf("cart is not valid for place order some products are out of stock")
	}

	return cart, nil
}
