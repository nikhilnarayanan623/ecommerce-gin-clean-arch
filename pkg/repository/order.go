package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
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
func (c *OrderDatabase) FindAllShopOrdersByUserID(ctx context.Context, userID uint, pagination req.ReqPagination) ([]res.ResShopOrder, error) {

	limit := pagination.Count
	offset := (pagination.PageNumber - 1) * limit

	var shopOrders []res.ResShopOrder
	query := `SELECT so.user_id, so.id AS shop_order_id, so.order_date, so.order_total_price,so.discount, 
	so.order_status_id, os.status AS order_status,so.address_id,so.payment_method_id, pm.payment_type  
	FROM shop_orders so JOIN order_statuses os ON so.order_status_id = os.id 
	INNER JOIN payment_methods pm ON so.payment_method_id = pm.id WHERE user_id = $1 ORDER BY order_date DESC LIMIT $2 OFFSET  $3`
	if c.DB.Raw(query, userID, limit, offset).Scan(&shopOrders).Error != nil {
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
func (c *OrderDatabase) FindAllShopOrders(ctx context.Context, pagination req.ReqPagination) (shopOrders []res.ResShopOrder, err error) {

	limit := pagination.Count
	offset := (pagination.PageNumber - 1) * limit

	query := `SELECT so.user_id, so.id AS shop_order_id, so.order_date, so.order_total_price,so.discount, 
	so.order_status_id, os.status AS order_status,so.address_id,so.payment_method_id, pm.payment_type  
	FROM shop_orders so JOIN order_statuses os ON so.order_status_id = os.id 
	INNER JOIN payment_methods pm ON so.payment_method_id = pm.id 
	ORDER BY so.order_date DESC LIMIT $1 OFFSET $2`
	if c.DB.Raw(query, limit, offset).Scan(&shopOrders).Error != nil {
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
func (c *OrderDatabase) FindAllOrdersItemsByShopOrderID(ctx context.Context, shopOrderID uint, pagination req.ReqPagination) (orderItems []res.ResOrderItem, err error) {

	limit := pagination.Count
	offset := (pagination.PageNumber - 1) * limit

	query := `SELECT ol.product_item_id,p.product_name,p.image,ol.price, so.order_date, os.status,ol.qty, (ol.price * ol.qty) AS sub_total FROM  order_lines ol 
	JOIN shop_orders so ON ol.shop_order_id = so.id JOIN product_items pi ON ol.product_item_id = pi.id
	JOIN products p ON pi.product_id = p.id JOIN order_statuses os ON so.order_status_id = os.id AND ol.shop_order_id= $1 
	ORDER BY ol.qty DESC LIMIT $2 OFFSET $3`

	if c.DB.Raw(query, shopOrderID, limit, offset).Scan(&orderItems).Error != nil {
		return orderItems, errors.New("faild to get users order list")
	}
	return orderItems, nil
}

// ! order place
func (c *OrderDatabase) CheckcartIsValidForOrder(ctx context.Context, userID uint) (cart domain.Cart, err error) {

	query := `SELECT * FROM carts WHERE user_id = ?`
	if c.DB.Raw(query, userID).Scan(&cart).Error != nil {
		return cart, errors.New("faild to get cartItem of user")
	}

	// check any of the product is out of stock in cart
	var idOfOneOutOfStockProduct uint
	query = `SELECT DISTINCT pi.id FROM product_items pi 
	INNER JOIN cart_items ci ON pi.id = ci.product_item_id 
	INNER JOIN carts c ON ci.cart_id = c.cart_id 
	 WHERE c.user_id = 2 AND pi.qty_in_stock <= 0`

	err = c.DB.Raw(query).Scan(&idOfOneOutOfStockProduct).Error
	if err != nil {
		return cart, fmt.Errorf("faild to check out of stock product in user cart with usre_id %d", userID)
	}

	if idOfOneOutOfStockProduct != 0 {
		return cart, fmt.Errorf("there is an out of stock product is in cart remove it to place order")
	}

	return cart, nil
}

func (c *OrderDatabase) SaveShopOrder(ctx context.Context, shopOrder domain.ShopOrder) (domain.ShopOrder, error) {
	trx := c.DB.Begin()

	// save the shop_order
	query := `INSERT INTO shop_orders (user_id,address_id, order_total_price, discount, 
	order_status_id,order_date, payment_method_id) 
	VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	orderDate := time.Now()
	if c.DB.Raw(query, shopOrder.UserID, shopOrder.AddressID, shopOrder.OrderTotalPrice, shopOrder.Discount,
		shopOrder.OrderStatusID, orderDate, shopOrder.PaymentMethodID).Scan(&shopOrder).Error != nil {
		trx.Rollback()
		return shopOrder, errors.New("faild to save shop_order")
	}

	if trx.Commit().Error != nil {
		trx.Rollback()
		return shopOrder, errors.New("faild to complete save shop_order")
	}

	return shopOrder, nil
}

func (c *OrderDatabase) CartItemToOrderLines(ctx context.Context, userID uint) ([]domain.OrderLine, error) {
	var orderLines []domain.OrderLine
	query := `SELECT ci.product_item_id, ci.qty, CASE WHEN pi.discount_price > 0 THEN pi.discount_price ELSE pi.price END 
	FROM cart_items ci INNER JOIN product_items pi ON ci.product_item_id = pi.id 
	INNER JOIN carts c ON ci.cart_id = c.cart_id  
	WHERE c.user_id = $1`
	if c.DB.Raw(query, userID).Scan(&orderLines).Error != nil {
		return orderLines, errors.New("faild to convert cart_item to order_lines")
	}
	return orderLines, nil
}

func (c *OrderDatabase) SaveOrderLine(ctx context.Context, orderLine domain.OrderLine) error {
	query := `INSERT INTO order_lines (product_item_id, shop_order_id, qty, price) 
	VALUES ($1, $2, $3, $4)`
	if c.DB.Exec(query, orderLine.ProductItemID, orderLine.ShopOrderID, orderLine.Qty, orderLine.Price).Error != nil {
		return errors.New("faild to save orde line")
	}

	return nil
}

func (c *OrderDatabase) DeleteOrderedCartItems(ctx context.Context, userID uint) error {

	query := ` DELETE FROM cart_items ci USING carts c WHERE c.cart_id = ci.cart_id AND c.user_id = $1`
	if c.DB.Exec(query, userID).Error != nil {
		return errors.New("faild to remove cart_items on order")
	}
	return nil
}

func (c *OrderDatabase) UpdateCouponUsedForUser(ctx context.Context, userID, couponID uint) error {

	query := `INSERT INTO coupon_uses (coupon_id, user_id, used_at) VALUES ($1, $2, $3)`
	usedAt := time.Now()
	err := c.DB.Exec(query, couponID, userID, usedAt).Error
	if err != nil {
		return errors.New("faild to updte user as couon useed")
	}
	return nil
}

func (c *OrderDatabase) ValidateAddressID(ctx context.Context, addressID uint) error {

	var id uint
	if c.DB.Raw(`SELECT id FROM addresses WHERE id = $1`, addressID).Scan(&id).Error != nil {
		return errors.New("faild to validte address_id")
	} else if id == 0 {
		return errors.New("invlaid address_id")
	}
	return nil
}

// get user email and phone
func (c *OrderDatabase) GetUserEmailAndPhone(ctx context.Context, userID uint) (emailAndPhone res.ResEmailAndPhone, err error) {

	query := `SELECT email,phone FROM users WHERE id = $1`
	err = c.DB.Raw(query, userID).Scan(&emailAndPhone).Error
	if err != nil {
		return emailAndPhone, fmt.Errorf("faild to find email and phone of user with user_id %v", userID)
	}
	return emailAndPhone, nil
}

//!end

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

// order return

func (c *OrderDatabase) FindOrderReturn(ctx context.Context, orderReturn domain.OrderReturn) (domain.OrderReturn, error) {
	query := `SELECT *  FROM order_returns WHERE id = ? OR shop_order_id = ?`
	if c.DB.Raw(query, orderReturn.ID, orderReturn.ShopOrderID).Scan(&orderReturn).Error != nil {
		return orderReturn, errors.New("faild to find order return")
	}
	return orderReturn, nil
}

func (c *OrderDatabase) FindAllOrderReturns(ctx context.Context, onlyPending bool, pagination req.ReqPagination) (orderReturns []res.ResOrderReturn, errr error) {

	limit := pagination.Count
	offset := (pagination.PageNumber - 1) * limit

	// var query string
	if onlyPending { // find all request which are not returned completed
		// find order_status_id for return requested
		var orderStatusID uint
		if c.DB.Raw("SELECT id FROM order_statuses WHERE status = 'return requested'").Scan(&orderStatusID).Error != nil {
			return orderReturns, errors.New("faild to get order_status_id for return requestes")
		}
<<<<<<< HEAD
		query := `SELECT ors.id AS order_return_id, ors.shop_order_id, ors.request_date, ors.return_reason, 
=======
		if c.DB.Raw("SELECT id FROM order_statuses WHERE status = 'return approved'").Scan(&orderStatusIDApproved).Error != nil {
			return orderReturns, errors.New("faild to get order_status_id for return requestes")
		}

		query := `SELECT ors.id AS order_return_id, ors.shop_order_id, ors.request_date, ors.return_reason, ors.admin_comment, ors.approval_date, 
>>>>>>> ef7579d (changed the user side coupon show api now user can see the coupon which are used and unused with pagination and order of unused coupons first)
		os.id AS order_status_id, os.status AS order_status,ors.refund_amount  
		FROM order_returns ors INNER JOIN shop_orders so ON ors.shop_order_id =  so.id 
		INNER JOIN order_statuses os ON so.order_status_id = os.id WHERE so.order_status_id = $1 
		ORDER BY ors.request_date DESC LIMIT $2 OFFSET $3`
		if c.DB.Raw(query, orderStatusID, limit, offset).Scan(&orderReturns).Error != nil {
			return orderReturns, errors.New("faild to find orders of return requested")
		}
	} else {
		query := `SELECT ors.id AS order_return_id, ors.shop_order_id, ors.request_date, ors.return_reason, 
		os.id AS order_status_id, os.status AS order_status,ors.refund_amount, ors.admin_comment, ors.is_approved, ors.approval_date, ors.return_date 
		FROM order_returns ors INNER JOIN shop_orders so ON ors.shop_order_id =  so.id 
		INNER JOIN order_statuses os ON so.order_status_id = os.id 
		ORDER BY ors.request_date LIMIT $1 OFFSET $2`
		if c.DB.Raw(query, limit, offset).Scan(&orderReturns).Error != nil {
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
func (c *OrderDatabase) UpdateOrderReturn(ctx context.Context, body req.ReqUpdatReturnOrder) error {
	trx := c.DB.Begin()

	// find the orderStatus that admin given
	changeOrderStatus, err := c.FindOrderStatus(ctx, domain.OrderStatus{ID: body.OrderStatusID})
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

	err = fmt.Errorf("faild to change return status to %s", changeOrderStatus.Status)
	switch changeOrderStatus.Status {
	case "return cancelled":
		query := `UPDATE order_returns SET admin_comment = $1`
		if c.DB.Exec(query, body.AdminComment).Error != nil {
			trx.Rollback()
			return err
		}
	case "return approved":
		approvalTime := time.Now()
		returnDate := approvalTime.AddDate(0, 0, 5) // now its approval plus 5 days need chage it to get return date from admin later
		query := `UPDATE order_returns SET return_date = $1, approval_date = $2, admin_comment = $3`
		if c.DB.Exec(query, returnDate, approvalTime, body.AdminComment).Error != nil {
			trx.Rollback()
			return err
		}
	case "returned":
		returnDate := time.Now()
		query := `UPDATE order_returns SET return_date = $1,admin_comment = $2`
		if c.DB.Exec(query, returnDate, body.AdminComment).Error != nil {
			trx.Rollback()
			return err
		}
	}

	// update the status on shop_order
	query := `UPDATE shop_orders SET order_status_id = $1 WHERE id = $2`
	if trx.Exec(query, body.OrderStatusID, orderReturn.ShopOrderID).Error != nil {
		trx.Rollback()
		return fmt.Errorf("faild to update shop_order_status on %s", changeOrderStatus.Status)
	}
	// complete the order_return update
	if trx.Commit().Error != nil {
		trx.Rollback()
		return errors.New("faild to complete updation of order_return")
	}
	return nil
}
