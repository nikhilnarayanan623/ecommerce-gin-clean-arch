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

type adminDatabase struct {
	DB *gorm.DB
}

func NewAdminRepository(DB *gorm.DB) interfaces.AdminRepository {
	return &adminDatabase{DB: DB}
}

func (c *adminDatabase) FindAdminByEmail(ctx context.Context, email string) (domain.Admin, error) {

	var admin domain.Admin
	err := c.DB.Raw("SELECT * FROM admins WHERE email = $1", email).Scan(&admin).Error

	return admin, err
}

func (c *adminDatabase) FindAdminByUserName(ctx context.Context, userName string) (domain.Admin, error) {

	var admin domain.Admin
	err := c.DB.Raw("SELECT * FROM admins WHERE user_name = $1", userName).Scan(&admin).Error

	return admin, err
}

func (c *adminDatabase) SaveAdmin(ctx context.Context, admin domain.Admin) error {

	querry := `INSERT INTO admins (user_name,email,password,created_at) VALUES ($1, $2, $3, $4)`
	createdAt := time.Now()
	if c.DB.Exec(querry, admin.UserName, admin.Email, admin.Password, createdAt).Error != nil {
		return errors.New("faild to save admin")
	}

	return nil
}

func (c *adminDatabase) FindAllUser(ctx context.Context, pagination req.ReqPagination) (users []res.UserRespStrcut, err error) {

	limit := pagination.Count
	offset := (pagination.PageNumber - 1) * limit

	query := `SELECT * FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err = c.DB.Raw(query, limit, offset).Scan(&users).Error

	return users, err
}

// sales report from order // !add  product wise report
func (c *adminDatabase) CreateFullSalesReport(ctc context.Context, reqData req.ReqSalesReport) (salesReport []res.SalesReport, err error) {

	limit := reqData.Pagination.Count
	offset := (reqData.Pagination.PageNumber - 1) * limit

	startDate := reqData.StartDate
	endDate := reqData.EndDate

	query := `SELECT u.first_name, u.email,  so.id AS shop_order_id, so.user_id, so.order_date, 
	so.order_total_price, so.discount, os.status AS order_status, pm.payment_type FROM shop_orders so
	INNER JOIN order_statuses os ON so.order_status_id = os.id 
	INNER JOIN  payment_methods pm ON so.payment_method_id = pm.id 
	INNER JOIN users u ON so.user_id = u.id 
	WHERE order_date >= $1 AND order_date <= $2
	ORDER BY so.order_date LIMIT  $3 OFFSET $4`

	if c.DB.Raw(query, startDate, endDate, limit, offset).Scan(&salesReport).Error != nil {
		return salesReport, errors.New("faild to collect data to create sales report")
	}

	return salesReport, nil
}

// stock side
func (c *adminDatabase) FindStockBySKU(ctx context.Context, sku string) (stock res.RespStock, err error) {
	query := `SELECT pi.sku, pi.qty_in_stock, pi.price, p.product_name, vo.variation_value 
	FROM product_items pi INNER JOIN products p ON p.id = pi.product_id 
	INNER JOIN product_configurations pc ON pc.product_item_id = pi.id 
	INNER JOIN variation_options vo ON vo.id = pc.variation_option_id
	WHERE pi.sku = $1`

	err = c.DB.Raw(query, sku).Scan(&stock).Error
	if err != nil {
		return stock, fmt.Errorf("faild to find stock detils of sku %v", sku)
	}

	return stock, nil
}

func (c *adminDatabase) FindAllStockDetails(ctx context.Context, pagination req.ReqPagination) (stocks []res.RespStock, err error) {

	limit := pagination.Count
	offset := (pagination.PageNumber - 1) * limit

	query := `SELECT pi.sku, pi.qty_in_stock, pi.price, p.product_name, vo.variation_value 
	FROM product_items pi INNER JOIN products p ON p.id = pi.product_id 
	INNER JOIN product_configurations pc ON pc.product_item_id = pi.id 
	INNER JOIN variation_options vo ON vo.id = pc.variation_option_id 
	ORDER BY qty_in_stock LIMIT $1 OFFSET $2`

	err = c.DB.Raw(query, limit, offset).Scan(&stocks).Error

	if err != nil {
		return stocks, fmt.Errorf("faild to find all stocks details from database")
	}

	return stocks, nil
}

func (c *adminDatabase) UpdateStock(ctx context.Context, valuesToUpdate req.ReqUpdateStock) error {

	query := `UPDATE product_items SET qty_in_stock = qty_in_stock + $1 WHERE sku = $2`

	err := c.DB.Exec(query, valuesToUpdate.QtyToAdd, valuesToUpdate.SKU).Error

	if err != nil {
		return fmt.Errorf("faild to update qty for product_item with sku %v", valuesToUpdate.SKU)
	}
	return nil
}
