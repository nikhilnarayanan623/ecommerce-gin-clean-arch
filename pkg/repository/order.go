package repository

import (
	"context"
	"errors"
	"time"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/res"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	"gorm.io/gorm"
)

type OrderDatabase struct {
	DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) interfaces.OrderRepository {
	return &OrderDatabase{DB: db}
}

func (c *OrderDatabase) GetOrdersListByUserID(ctx context.Context, userID uint) ([]res.ResOrder, error) {

	var orderList []res.ResOrder

	query := `SELECT ol.product_item_id, p.product_name,p.image, pi.price,ol.qty, (pi.price * ol.qty) AS total_price,os.status FROM order_lines ol 
	JOIN shop_orders so ON ol.shop_order_id = so.id JOIN order_statuses os ON so.order_status_id=os.id 
	JOIN product_items pi ON ol.product_item_id = pi.id JOIN products p ON pi.product_id = p.id AND so.user_id=?`

	if c.DB.Raw(query, userID).Scan(&orderList).Error != nil {
		return orderList, errors.New("faild to get users order list")
	}
	return orderList, nil
}

func (c *OrderDatabase) PlaceOrderByCart(ctx context.Context, shopOrder domain.ShopOrder) error {

	trx := c.DB.Begin()

	// find the user cart
	var cart domain.Cart
	if trx.Raw("SELECT * FROM carts WHERE user_id=?", shopOrder.UserID).Scan(&cart).Error != nil {
		trx.Rollback()
		return errors.New("faild to get user cart")
	} else if cart.ID == 0 {
		return errors.New("user have no cart")
	}

	// get all cartItems of user cart which are not out of stock
	var cartItems []domain.CartItem
	query := `SELECT * FROM cart_items ci JOIN product_items pi ON ci.product_item_id = pi.id AND pi.qty_in_stock > 0  AND ci.cart_id=?`
	if trx.Raw(query, cart.ID).Scan(&cartItems).Scan(&cartItems).Error != nil {
		trx.Rollback()
		return errors.New("there is no cartItems in the user cart")
		//return errors.New("faild to get cartItems of user")
	} //else if cartItems == nil { // there is no cartItems for user
	///	return errors.New("there is no cartItems in the user cart")
	///}

	if shopOrder.AddressID == 0 { // if address id not given then get user default address
		query := `SELECT a.id AS adress_id FROM user_addresses ua JOIN addresses a ON ua.address_id = a.id AND ua.user_id=? AND ua.is_default='t'`
		if trx.Raw(query, shopOrder.UserID).Scan(&shopOrder.AddressID).Error != nil {
			trx.Rollback()
			return errors.New("faild to get user default address")
		}
	} else { // validate addressId with user
		query := `SELECT a.id FROM user_addresses ua JOIN addresses a ON ua.address_id = a.id AND ua.user_id=? AND a.id=?`
		var AddressID uint
		if trx.Raw(query, shopOrder.UserID, shopOrder.AddressID).Scan(&AddressID).Error != nil {
			trx.Rollback()
			return errors.New("faild to get user addres id")
		} else if AddressID == 0 {
			return errors.New("invalid address id")
		}
	}

	var orderStatusID int
	if trx.Raw("SELECT id FROM order_statuses WHERE status = 'pending'").Scan(&orderStatusID).Error != nil {
		trx.Rollback()
		return errors.New("faild to find order status")
	} else if orderStatusID == 0 { // pending status id not in the order status then crete new one
		if trx.Raw("INSERT INTO order_statuses (status) VALUES ('pending') RETURNING id").Scan(&orderStatusID).Error != nil {
			trx.Rollback()
			return errors.New("peding status not in the order_status and faild to insert new one")
		}
	}

	//place an full order
	query = `INSERT INTO shop_orders (user_id,order_date,address_id,order_total_price,order_status_id,cod) VALUES 
	($1,$2,$3,$4,$5,$6) RETURNING *`
	orderDate := time.Now()

	if trx.Raw(query, shopOrder.UserID, orderDate, shopOrder.AddressID, cart.TotalPrice, orderStatusID, shopOrder.COD).Scan(&shopOrder).Error != nil {
		trx.Rollback()
		return errors.New("faild to place order")
	}

	// make order line for each cartItems
	query = `INSERT INTO order_lines (product_item_id,shop_order_id,qty,price) VALUES ($1,$2,$3,$4)`
	var (
		subTotal, price uint
		orderLine       domain.OrderLine
	)
	for _, cartItem := range cartItems {
		// get productitem price
		if trx.Raw("SELECT price FROM product_items WHERE id = ?", cartItem.ProductItemID).Scan(&price).Error != nil {
			trx.Rollback()
			return errors.New("faild to get producItem's price")
		}
		// calculate cartItem subtotal price
		subTotal = price * cartItem.Qty
		// make order line
		if trx.Raw(query, cartItem.ProductItemID, shopOrder.ID, cartItem.Qty, subTotal).Scan(&orderLine).Error != nil {
			trx.Rollback()
			return errors.New("faild to create order line")
		}
	}

	// delete cart items which are not out of stock
	query = `DELETE FROM cart_items USING product_items WHERE cart_items.product_item_id = product_items.id AND product_items.qty_in_stock > 0 AND cart_items.cart_id = ?`
	if trx.Raw(query, cart.ID).Scan(&cartItems).Error != nil {
		return errors.New("faild to clear the cart of user")
	}

	// atlast commit the transaction to complete the order
	if trx.Commit().Error != nil {
		trx.Rollback()
		return errors.New("faild to commit order")
	}
	return nil
}
