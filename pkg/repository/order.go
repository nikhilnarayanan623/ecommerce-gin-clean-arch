package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/req"
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
	query := `SELECT user_id, so.id AS shop_order_id, so.order_date, so.order_total_price,so.order_status_id,os.status AS order_status,so.address_id,so.cod FROM shop_orders so 
	JOIN order_statuses os ON so.order_status_id = os.id  WHERE user_id = ?`
	if c.DB.Raw(query, userID).Scan(&shopOrders).Error != nil {
		return shopOrders, errors.New("faild to get user shop order")
	}

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

	query := `SELECT ol.product_item_id,p.product_name,p.image,ol.price, so.order_date, os.status,ol.qty, (ol.price * ol.qty) AS sub_total FROM  order_lines ol 
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
	if trx.Raw(query, cart.ID).Scan(&cartItems).Error != nil {
		trx.Rollback()
		return errors.New("there is no cartItems in the user cart")
	} else if cartItems == nil {
		trx.Rollback()
		return errors.New("invalid api call \nthere is no product_itmes in cart")
	}

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
		price     uint
		orderLine domain.OrderLine
	)
	for _, cartItem := range cartItems {
		// get productitem price if discount price is there take discount price other wise take price
		if trx.Raw("SELECT CASE WHEN discount_price > 0 THEN discount_price ELSE price END AS price FROM product_items WHERE id = ?", cartItem.ProductItemID).Scan(&price).Error != nil {
			trx.Rollback()
			return errors.New("faild to get producItem's price")
		}
		// calculate cartItem subtotal price
		// make order line
		if trx.Raw(query, cartItem.ProductItemID, shopOrder.ID, cartItem.Qty, price).Scan(&orderLine).Error != nil {
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
	}
	return orderStatus, nil
}

func (c *OrderDatabase) FindAllOrderStauses(ctx context.Context) ([]domain.OrderStatus, error) {
	var orderStatuses []domain.OrderStatus
	if c.DB.Raw("SELECT * FROM order_statuses").Scan(&orderStatuses).Error != nil {
		return orderStatuses, errors.New("faild to get all order_statuses")
	}
	return orderStatuses, nil
}

// admin side status change
func (c *OrderDatabase) UpdateShopOrderOrderStatus(ctx context.Context, shopOrderID, changeStatusID uint) error {

	// any other change the status
	query := `UPDATE shop_orders SET order_status_id = ? WHERE id = ?`
	if c.DB.Exec(query, changeStatusID, shopOrderID).Error != nil {
		return errors.New("faild to update status of order")
	}

	return nil
}

// checkout page
func (c *OrderDatabase) CheckOutCart(ctx context.Context, userId uint) (res.ResCheckOut, error) {

	var resCheckOut res.ResCheckOut
	// get all cartItems of user which are not out of stock
	query := `SELECT ci.product_item_id, p.product_name,pi.price,pi.discount_price, pi.qty_in_stock, ci.qty, 
	CASE WHEN pi.discount_price > 0 THEN (ci.qty * pi.discount_price) ELSE (ci.qty * pi.price) END AS sub_total  
	FROM cart_items ci JOIN carts c ON ci.cart_id = c.id JOIN product_items pi ON ci.product_item_id = pi.id 
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

// order return

func (c *OrderDatabase) FindOrderReturn(ctx context.Context, orderReturn domain.OrderReturn) (domain.OrderReturn, error) {
	query := `SELECT *  FROM order_returns WHERE id = ? OR shop_order_id = ?`
	if c.DB.Raw(query, orderReturn.ID, orderReturn.ShopOrderID).Scan(&orderReturn).Error != nil {
		return orderReturn, errors.New("faild to find order return")
	}
	return orderReturn, nil
}

func (c *OrderDatabase) FindAllOrderReturns(ctx context.Context, onlyPending bool) ([]domain.OrderReturn, error) {
	var orderReturns []domain.OrderReturn

	if onlyPending { // find all request which are not returned completed
		// find order_status_id for return
		var orderStatusID uint
		if c.DB.Raw("SELECT id FROM order_statuses WHERE status = 'returned'").Scan(&orderStatusID).Error != nil {
			return orderReturns, errors.New("faild to get order_status_id for return request ")
		}
		query := `SELECT ors.id, ors.shop_order_id, ors.request_date, ors.return_reason, 
		ors.return_date, ors.approval_date, ors.refund_amount, ors.is_approved,admin_comment 
		FROM order_returns ors INNER JOIN shop_orders so ON so.id = ors.shop_order_id 
		AND so.order_status_id != ?`
		if c.DB.Raw(query, orderStatusID).Scan(&orderReturns).Error != nil {
			return orderReturns, errors.New("faild to find orders of return requested")
		}
	} else {
		query := `SELECT * FROM order_returns`
		if c.DB.Raw(query).Scan(&orderReturns).Error != nil {
			return orderReturns, errors.New("faild to get order returns")
		}
	}

	return orderReturns, nil
}

// to save a return request
func (c *OrderDatabase) SaveOrderReturn(ctx context.Context, orderReturn domain.OrderReturn) error {

	trx := c.DB.Begin()

	query := `INSERT INTO order_returns (shop_order_id,return_reason,request_date,refund_amount,is_approved) 
	VALUES($1,$2,$3,$4,$5)`
	if trx.Exec(query, orderReturn.ShopOrderID, orderReturn.ReturnReason,
		orderReturn.RequestDate, orderReturn.RefundAmount, false).Error != nil {
		trx.Rollback()
		return fmt.Errorf("faild to save return for shop_order_id %d", orderReturn.ShopOrderID)
	}

	//get the returning order status id and set to order
	var orderStatus = domain.OrderStatus{Status: "return requested"}
	orderStatus, err := c.FindOrderStatus(ctx, orderStatus)
	if err != nil {
		trx.Rollback()
		return err
	} else if orderStatus.ID == 0 {
		trx.Rollback()
		return errors.New("faild get order_status_id of returning")
	}

	//update shopOrder status
	if err := c.UpdateShopOrderOrderStatus(ctx, orderReturn.ShopOrderID, orderStatus.ID); err != nil {
		trx.Rollback()
		return err
	}

	if err := trx.Commit().Error; err != nil {
		trx.Rollback()
		return fmt.Errorf("faild to complete return request for shop_order_id %d", orderReturn.ShopOrderID)
	}
	return nil
}

// update the order return
func (c *OrderDatabase) UpdateOrderReturn(ctx context.Context, body req.ReqUpdatReturnReq) error {
	trx := c.DB.Begin()

	// find the orderStatus that admin given
	orderStatus, err := c.FindOrderStatus(ctx, domain.OrderStatus{ID: body.OrderStatusID})
	if err != nil {
		trx.Rollback()
		return err
	}

	//find the orderReturn
	orderReturn, err := c.FindOrderReturn(ctx, domain.OrderReturn{ID: body.OrderReturnID})
	if err != nil {
		trx.Rollback()
		return err
	}

	// do updattion according to the order status
	if orderStatus.Status == "return cancelled" {

		query := `UPDATE order_returns SET admin_comment = $1`
		if trx.Exec(query, body.AdminComment).Error != nil {
			return fmt.Errorf("faild to update order_returns on %s", orderStatus.Status)
		}
	} else if orderStatus.Status == "return approved" {
		approvalDate := time.Now()
		returnDate := approvalDate.Add(time.Hour * 24 * 5) // manual return date set to 5 days from approval

		query := `UPDATE order_returns SET approval_date = $1, return_date = $2, admin_comment = $3,is_approved = $4`
		if trx.Exec(query, approvalDate, returnDate, body.AdminComment, true).Error != nil {
			return fmt.Errorf("faild to update order_returns on %s", orderStatus.Status)
		}
	} else if orderStatus.Status == "returned" {
		// updae the return date as current date
		returnDate := time.Now()
		query := `UPDATE order_returns SET return_date = $1, admin_comment = $2`
		if trx.Exec(query, returnDate, body.AdminComment).Error != nil {
			return fmt.Errorf("faild to update order_returns on %s", orderStatus.Status)
		}
	}

	// update the status on shop_order
	query := `UPDATE shop_orders SET order_status_id = $1 WHERE id = $2`
	if trx.Exec(query, body.OrderStatusID, orderReturn.ShopOrderID).Error != nil {
		trx.Rollback()
		return fmt.Errorf("faild to update shop_order_status on %s", orderStatus.Status)
	}
	// complete the order_return update
	if trx.Commit().Error != nil {
		return errors.New("faild to complete updation of order_return")
	}
	return nil
}
