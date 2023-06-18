package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/response"
	"gorm.io/gorm"
)

type OrderDatabase struct {
	DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) interfaces.OrderRepository {
	return &OrderDatabase{DB: db}
}

func (c *OrderDatabase) Transaction(callBack func(trxRepo interfaces.OrderRepository) error) error {

	trx := c.DB.Begin()
	transactionRepo := NewOrderRepository(trx)

	err := callBack(transactionRepo)
	if err != nil {
		trx.Rollback()
		return fmt.Errorf("faild to complete transaction \nerror:%v", err.Error())
	}

	err = trx.Commit().Error
	return err
}

func (c *OrderDatabase) IsShpoOrderIDIsValid(ctx context.Context, shopOrderID uint) error {
	var exist bool
	query := `SELECT EXISTS(SELECT 1 FROM shop_orders WHERE id = $1) AS exist`
	err := c.DB.Raw(query, shopOrderID).Scan(&exist).Error
	if err != nil {
		return err
	} else if !exist {
		return fmt.Errorf("shop_order_id not exist")
	}

	return nil
}

// find a specific shop order by shopOrderID
func (c *OrderDatabase) FindShopOrderByShopOrderID(ctx context.Context, shopOrderID uint) (shopOrder domain.ShopOrder, err error) {

	query := `SELECT * FROM shop_orders WHERE id = $1`
	err = c.DB.Raw(query, shopOrderID).Scan(&shopOrder).Error

	return shopOrder, err
}

// get all shop order of user
func (c *OrderDatabase) FindAllShopOrdersByUserID(ctx context.Context, userID uint, pagination request.Pagination) ([]response.ShopOrder, error) {

	limit := pagination.Count
	offset := (pagination.PageNumber - 1) * limit

	var shopOrders []response.ShopOrder
	query := `SELECT so.user_id, so.id AS shop_order_id, so.order_date, so.order_total_price,so.discount, 
	so.order_status_id, os.status AS order_status,so.address_id, 
	CASE WHEN so.payment_method_id != 0 
	THEN (SELECT payment_type FROM payment_methods WHERE id = so.payment_method_id) ELSE '' END AS payment_type  
	FROM shop_orders so JOIN order_statuses os ON so.order_status_id = os.id 
	WHERE user_id = $1 
	ORDER BY order_date DESC LIMIT $2 OFFSET  $3`
	err := c.DB.Raw(query, userID, limit, offset).Scan(&shopOrders).Error
	if err != nil {
		return shopOrders, err
	}
	// take full address and add to it
	query = `SELECT adrs.id AS address_id, adrs.name,adrs.phone_number,adrs.house,adrs.area, adrs.land_mark,adrs.city,adrs.pincode,adrs.country_id,c.country_name 
	FROM addresses adrs JOIN countries c ON adrs.country_id = c.id  WHERE adrs.id= ?`
	var address response.Address
	for i, order := range shopOrders {

		if c.DB.Raw(query, order.AddressID).Scan(&address).Error != nil {
			return shopOrders, errors.New("faild to get addresses")
		}
		shopOrders[i].Address = address
	}
	return shopOrders, nil
}

// find all shop orders with user
func (c *OrderDatabase) FindAllShopOrders(ctx context.Context, pagination request.Pagination) (shopOrders []response.ShopOrder, err error) {

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

	var address response.Address
	query = `SELECT adrs.id, adrs.name,adrs.phone_number,adrs.house,adrs.area, adrs.land_mark,adrs.city,adrs.pincode,adrs.country_id,c.country_name 
	FROM addresses adrs JOIN countries c ON adrs.country_id = c.id  WHERE adrs.id= ?`
	for i, order := range shopOrders {

		if c.DB.Raw(query, order.AddressID).Scan(&address).Error != nil {
			return shopOrders, errors.New("faild to get addresses")
		}
		shopOrders[i].Address = address
	}
	return shopOrders, nil
}

// get order items of a specific order
func (c *OrderDatabase) FindAllOrdersItemsByShopOrderID(ctx context.Context, shopOrderID uint, pagination request.Pagination) (orderItems []response.OrderItem, err error) {

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

func (c *OrderDatabase) SaveShopOrder(ctx context.Context, shopOrder domain.ShopOrder) (shopOrderID uint, err error) {

	// save the shop_order
	query := `INSERT INTO shop_orders (user_id,address_id, order_total_price, discount, 
	order_status_id,order_date) 
	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	orderDate := time.Now()
	if c.DB.Raw(query, shopOrder.UserID, shopOrder.AddressID, shopOrder.OrderTotalPrice, shopOrder.Discount,
		shopOrder.OrderStatusID, orderDate).Scan(&shopOrderID).Error != nil {
		return 0, errors.New("faild to save shop_order")
	}

	return shopOrderID, nil
}

func (c *OrderDatabase) SaveOrderLine(ctx context.Context, orderLine domain.OrderLine) error {
	query := `INSERT INTO order_lines (product_item_id, shop_order_id, qty, price) 
	VALUES ($1, $2, $3, $4)`
	if c.DB.Exec(query, orderLine.ProductItemID, orderLine.ShopOrderID, orderLine.Qty, orderLine.Price).Error != nil {
		return errors.New("faild to save orde line")
	}

	return nil
}

//!end

func (c *OrderDatabase) FindOrderStatusByShopOrderID(ctx context.Context, shopOrderID uint) (orderStatus domain.OrderStatus, err error) {
	query := `SELECT * FROM order_statuses
	WHERE id = (SELECT order_status_id FROM shop_orders WHERE id = $1)`
	err = c.DB.Raw(query, shopOrderID).Scan(&orderStatus).Error

	return orderStatus, err
}

// find order status
func (c *OrderDatabase) FindOrderStatusByID(ctx context.Context, orderStatusID uint) (domain.OrderStatus, error) {

	var orderStatus domain.OrderStatus
	err := c.DB.Raw("SELECT * FROM order_statuses WHERE id = $1", orderStatusID).Scan(&orderStatus).Error

	return orderStatus, err
}

func (c *OrderDatabase) FindOrderStatusByStatus(ctx context.Context, orderSatatus string) (domain.OrderStatus, error) {

	var orderStatus domain.OrderStatus
	err := c.DB.Raw("SELECT * FROM order_statuses WHERE status = $1", orderSatatus).Scan(&orderStatus).Error

	return orderStatus, err
}

func (c *OrderDatabase) FindAllOrderStauses(ctx context.Context) ([]domain.OrderStatus, error) {
	var orderStatuses []domain.OrderStatus
	if c.DB.Raw("SELECT * FROM order_statuses").Scan(&orderStatuses).Error != nil {
		return orderStatuses, errors.New("faild to get all order_statuses")
	}
	return orderStatuses, nil
}

func (c *OrderDatabase) UpdateShopOrderOrderStatus(ctx context.Context, shopOrderID, changeStatusID uint) error {

	query := `UPDATE shop_orders SET order_status_id = $1 WHERE id = $2`
	err := c.DB.Exec(query, changeStatusID, shopOrderID).Error

	return err
}

func (c *OrderDatabase) UpdateShopOrderStatusAndPaymentID(ctx context.Context, shopOrderID, statusID, paymentID uint) error {
	query := `UPDATE shop_orders SET order_status_id = $1 , payment_method_id = $2 WHERE id = $3`
	err := c.DB.Exec(query, statusID, paymentID, shopOrderID).Error

	return err
}

func (c *OrderDatabase) FindOrderReturnByReturnID(ctx context.Context, orderReturnID uint) (orderReturn domain.OrderReturn, err error) {

	query := `SELECT *  FROM order_returns WHERE id = $1`
	err = c.DB.Raw(query, orderReturnID).Scan(&orderReturn).Error

	return orderReturn, err
}
func (c *OrderDatabase) FindOrderReturnByShopOrderID(ctx context.Context, shopOrderID uint) (orderReturn domain.OrderReturn, err error) {

	query := `SELECT *  FROM order_returns WHERE shop_order_id = $1`
	err = c.DB.Raw(query, shopOrderID).Scan(&orderReturn).Error

	return orderReturn, err
}

func (c *OrderDatabase) FindAllOrderReturns(ctx context.Context, pagination request.Pagination) ([]response.OrderReturn, error) {

	limit := pagination.Count
	offset := (pagination.PageNumber - 1) * limit
	var orderReturns []response.OrderReturn

	query := `SELECT ors.id AS order_return_id, ors.shop_order_id, ors.request_date, ors.return_reason, 
		os.id AS order_status_id, os.status AS order_status,ors.refund_amount, ors.admin_comment, ors.is_approved, ors.approval_date, ors.return_date 
		FROM order_returns ors INNER JOIN shop_orders so ON ors.shop_order_id =  so.id 
		INNER JOIN order_statuses os ON so.order_status_id = os.id 
		ORDER BY ors.request_date LIMIT $1 OFFSET $2`
	err := c.DB.Raw(query, limit, offset).Scan(&orderReturns).Error

	return orderReturns, err
}

func (c *OrderDatabase) FindAllPendingOrderReturns(ctx context.Context, pagination request.Pagination) ([]response.OrderReturn, error) {
	limit := pagination.Count
	offset := (pagination.PageNumber - 1) * limit
	var pendingorderReturns []response.OrderReturn

	returnuested, err1 := c.FindOrderStatusByStatus(ctx, "return requested")
	returnApproved, err2 := c.FindOrderStatusByStatus(ctx, "return approved")
	err := errors.Join(err1, err2)
	if err != nil {
		return pendingorderReturns, err
	}

	query := `SELECT ors.id AS order_return_id, ors.shop_order_id, ors.request_date, ors.return_reason, 
	os.id AS order_status_id, os.status AS order_status,ors.refund_amount  
	FROM order_returns ors INNER JOIN shop_orders so ON ors.shop_order_id =  so.id 
	INNER JOIN order_statuses os ON so.order_status_id = os.id 
	WHERE so.order_status_id = $1 OR so.order_status_id = $2 
	ORDER BY ors.request_date DESC LIMIT $3 OFFSET $4`
	err = c.DB.Raw(query, returnuested.ID, returnApproved.ID, limit, offset).Scan(&pendingorderReturns).Error

	return pendingorderReturns, err
}

// to save a return request
func (c *OrderDatabase) SaveOrderReturn(ctx context.Context, orderReturn domain.OrderReturn) error {

	query := `INSERT INTO order_returns (shop_order_id,return_reason,request_date,refund_amount,is_approved) 
	VALUES($1,$2,$3,$4,$5)`
	err := c.DB.Exec(query, orderReturn.ShopOrderID, orderReturn.ReturnReason,
		orderReturn.RequestDate, orderReturn.RefundAmount, false).Error

	return err
}

// update the order return
func (c *OrderDatabase) UpdateOrderReturn(ctx context.Context, orderReturn domain.OrderReturn) error {

	query := `UPDATE order_returns SET admin_comment = $1, return_date = $2, 
	approval_date = $3, is_approved = $4 WHERE id = $5`
	err := c.DB.Exec(query, orderReturn.AdminComment, orderReturn.RequestDate,
		orderReturn.ApprovalDate, orderReturn.IsApproved, orderReturn.ID).Error

	return err
}
