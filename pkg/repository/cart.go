package repository

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/response"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
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

	query := `INSERT INTO carts (user_id,total_price) VALUES($1, $2) RETURNING id`
	err = c.DB.Raw(query, userID, 0).Scan(&cartID).Error

	return cartID, err
}

func (c *cartDatabase) UpdateCart(ctx context.Context, cartId, discountAmount, couponID uint) error {

	query := `UPDATE carts SET discount_amount = $1, applied_coupon_id = $2 WHERE id = $3`
	err := c.DB.Exec(query, discountAmount, couponID, cartId).Error

	return err
}

// find cart_items
func (c *cartDatabase) FindCartItemByID(ctx context.Context, cartItemID uint) (cartItem domain.CartItem, err error) {
	query := `SELECT * FROM cart_items WHERE id = ?`
	err = c.DB.Raw(query, cartItemID).Scan(&cartItem).Error

	return
}

func (c *cartDatabase) FindCartItemByCartAndProductItemID(ctx context.Context, cartID, productItemID uint) (cartItem domain.CartItem, err error) {
	query := `SELECT * FROM cart_items WHERE cart_id = $1 AND product_item_id = $2`
	err = c.DB.Raw(query, cartID, productItemID).Scan(&cartItem).Error

	return cartItem, err
}

func (c *cartDatabase) SaveCartItem(ctx context.Context, cartId, productItemId uint) error {

	query := `INSERT INTO cart_items (cart_id, product_item_id, qty) VALUES ($1, $2, $3)`
	err := c.DB.Exec(query, cartId, productItemId, 1).Error

	return err
}

func (c *cartDatabase) DeleteCartItem(ctx context.Context, cartItemID uint) error {

	query := `DELETE FROM cart_items WHERE id = $1`
	err := c.DB.Exec(query, cartItemID).Error

	return err
}

func (c *cartDatabase) DeleteAllCartItemsByCartID(ctx context.Context, cartID uint) error {

	query := ` DELETE FROM cart_items WHERE cart_id = $1`
	err := c.DB.Exec(query, cartID).Error
	return err
}

func (c *cartDatabase) UpdateCartItemQty(ctx context.Context, cartItemId, qty uint) error {

	query := `UPDATE cart_items SET qty = $1 WHERE id = $2`
	err := c.DB.Exec(query, qty, cartItemId).Error

	return err
}

func (c *cartDatabase) FindAllCartItemsByCartID(ctx context.Context, cartID uint) (cartItems []response.CartItem, err error) {

	// get the cartItem of all user with subtotal
	query := `SELECT ci.product_item_id, p.name AS product_name, ci.qty,pi.price ,
	 pi.discount_price, pi.qty_in_stock,
	 CASE WHEN pi.discount_price > 0 THEN pi.discount_price * ci.qty ELSE pi.price * ci.qty END AS sub_total   
	 FROM cart_items ci INNER JOIN product_items pi ON ci.product_item_id = pi.id 
	 INNER JOIN products p ON pi.product_id = p.id AND ci.cart_id=?`

	err = c.DB.Raw(query, cartID).Scan(&cartItems).Error

	return
}

func (c *cartDatabase) IsCartValidForOrder(ctx context.Context, userID uint) (valid bool, err error) {

	var outOfStockExist bool
	query := `SELECT 
		EXISTS( SELECT DISTINCT pi.id FROM product_items pi 
		INNER JOIN cart_items ci ON pi.id = ci.product_item_id 
		INNER JOIN carts c ON ci.cart_id = c.id 
		WHERE c.user_id = $1 AND pi.qty_in_stock <= 0) AS valid FROM carts`

	err = c.DB.Raw(query, userID).Scan(&outOfStockExist).Error

	// if error or a product is found a product is out
	if err != nil || outOfStockExist {
		return false, err
	}

	return true, nil
}
