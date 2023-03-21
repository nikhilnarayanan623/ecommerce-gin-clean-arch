package repository

import (
	"context"
	"errors"
	"fmt"
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

// find a specific shop order by shopOrderID
func (c *OrderDatabase) FindShopOrderByShopOrderID(ctx context.Context, shopOrderID uint) (domain.ShopOrder, error) {

	var shopOrder domain.ShopOrder
	if c.DB.Raw("SELECT * FROM shop_orders WHERE id = ?", shopOrderID).Scan(&shopOrder).Error != nil {
		return shopOrder, errors.New("faild to get shop order")
	}

	return shopOrder, nil
}

// get all shop order of user
func (c *OrderDatabase) FindAllShopOrdersByUserID(ctx context.Context, userID uint) ([]res.ResShopOrder, error) {

	var shopOrders []res.ResShopOrder
	fmt.Println("\n\n\n jjjjjjjjjjjjjj")
	query := `SELECT so.id AS shop_order_id, so.order_date, so.order_total_price,so.order_status_id,os.status AS order_status,so.address_id,so.cod FROM shop_orders so 
	JOIN order_statuses os ON so.order_status_id = os.id  WHERE user_id = ?`
	if c.DB.Raw(query, userID).Scan(&shopOrders).Error != nil {
		return shopOrders, errors.New("faild to get user shop order")
	}
	fmt.Println("\n\n\nstetestewsts")
	// take full address and add to it
	query = `SELECT adrs.id AS address_id, adrs.name,adrs.phone_number,adrs.house,adrs.area, adrs.land_mark,adrs.city,adrs.pincode,adrs.country_id,c.country_name 
	FROM addresses adrs JOIN countries c ON adrs.country_id = c.id  WHERE adrs.id= ?`
	var address res.ResAddress
	for i, order := range shopOrders {

		if c.DB.Raw(query, order.AddressID).Scan(&address).Error != nil {
			return shopOrders, errors.New("faild to get addresses")
		}
		fmt.Println(address, order.AddressID)
		shopOrders[i].Address = address
	}
	return shopOrders, nil
}

// find all shop orders with user
func (c *OrderDatabase) FindAllShopOrders(ctx context.Context) ([]res.ResShopOrder, error) {

	var shopOrders []res.ResShopOrder
	query := `SELECT so.user_id, so.id AS shop_order_id, so.order_date, so.order_total_price,so.order_status_id,os.status AS order_status,so.address_id, so.cod FROM shop_orders so 
	JOIN order_statuses os ON so.order_status_id = os.id `
	if c.DB.Raw(query).Scan(&shopOrders).Error != nil {
		return shopOrders, errors.New("faild to get order list")
	}

	var address res.ResAddress
	query = `SELECT adrs.id, adrs.name,adrs.phone_number,adrs.house,adrs.area, adrs.land_mark,adrs.city,adrs.pincode,adrs.country_id,c.country_name 
	FROM addresses adrs JOIN countries c ON adrs.country_id = c.id  WHERE adrs.id= ?`
	for i, order := range shopOrders {

		if c.DB.Raw(query, order.AddressID).Scan(&address).Error != nil {
			return shopOrders, errors.New("faild to get addresses")
		}
		fmt.Println(address, order.AddressID)
		shopOrders[i].Address = address
	}
	return shopOrders, nil
}

// get order items of a specific order
func (c *OrderDatabase) FindAllOrdersItemsByShopOrderID(ctx context.Context, shopOrderID uint) ([]res.ResOrder, error) {

	var orderList []res.ResOrder

	query := `SELECT ol.product_item_id,p.product_name,p.image,pi.price, so.order_date, os.status,ol.qty, (pi.price * ol.qty) AS sub_total FROM  order_lines ol 
	JOIN shop_orders so ON ol.shop_order_id = so.id JOIN product_items pi ON ol.product_item_id = pi.id
	JOIN products p ON pi.product_id = p.id JOIN order_statuses os ON so.order_status_id = os.id AND ol.shop_order_id= ?`

	if c.DB.Raw(query, shopOrderID).Scan(&orderList).Error != nil {
		return orderList, errors.New("faild to get users order list")
	}
	return orderList, nil
}

// save a new order from user cart
func (c *OrderDatabase) SaveOrderByCart(ctx context.Context, shopOrder domain.ShopOrder) error {

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
	query := `SELECT ci.cart_id,ci.product_item_id, ci.qty FROM cart_items ci JOIN product_items pi ON ci.product_item_id = pi.id AND pi.qty_in_stock > 0 AND ci.cart_id=?`
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

// find order status
func (c *OrderDatabase) FindOrderStatus(ctx context.Context, orderStatus domain.OrderStatus) (domain.OrderStatus, error) {

	if c.DB.Raw("SELECT * FROM order_statuses WHERE id = ? OR status = ?", orderStatus.ID, orderStatus.Status).Scan(&orderStatus).Error != nil {
		return orderStatus, errors.New("faild to get order status")
	} else if orderStatus.ID == 0 {
		return orderStatus, errors.New("invalid order status id")
	}
	return orderStatus, nil
}

// admin side status change
func (c *OrderDatabase) UpdateOrderStatus(ctx context.Context, shopOrder domain.ShopOrder, changeStatusID uint) error {

	// any other change the status
	query := `UPDATE shop_orders SET order_status_id = ? WHERE id = ?`
	if c.DB.Raw(query, changeStatusID, shopOrder.ID).Scan(&shopOrder).Error != nil {
		return errors.New("faild to update status of order")
	}

	return nil
}

// checkout page
func (c *OrderDatabase) CheckOutCart(ctx context.Context, userId uint) (res.ResCheckOut, error) {

	var resCheckOut res.ResCheckOut
	// get all cartItems of user which are not out of stock
	query := `SELECT ci.product_item_id, p.product_name,pi.price, pi.qty_in_stock, ci.qty, (pi.price * ci.qty) AS sub_total  FROM cart_items ci JOIN carts c ON ci.cart_id = c.id JOIN product_items pi ON ci.product_item_id = pi.id 
	JOIN products p ON pi.product_id = p.id AND c.user_id = ?`
	if c.DB.Raw(query, userId).Scan(&resCheckOut.ProductItems).Error != nil {
		return resCheckOut, errors.New("faild to get cartItems for checkout")
	}

	// get user addresses
	query = `SELECT * FROM addresses adrs JOIN user_addresses uadrs ON uadrs.address_id = adrs.id AND uadrs.user_id = ?`
	if c.DB.Raw(query, userId).Scan(&resCheckOut.Addresses).Error != nil {
		return resCheckOut, errors.New("faild to get user addrss for checkout")
	}

	// find total price
	query = `SELECT total_price FROM carts WHERE user_id = ?`
	if c.DB.Raw(query, userId).Scan(&resCheckOut.TotalPrice).Error != nil {
		return resCheckOut, errors.New("faild to ge total price for user cart")
	}
	return resCheckOut, nil
}
