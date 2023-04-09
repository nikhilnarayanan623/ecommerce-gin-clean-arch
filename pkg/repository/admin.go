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

func (c *adminDatabase) FindAdmin(ctx context.Context, admin domain.Admin) (domain.Admin, error) {

	if c.DB.Raw("SELECT * FROM admins WHERE email=? OR user_name=?", admin.Email, admin.UserName).Scan(&admin).Error != nil {
		return admin, errors.New("faild to find admin")
	}

	return admin, nil
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

func (c *adminDatabase) BlockUser(ctx context.Context, userID uint) error {

	// first check ther user is valid or not
	var user domain.User
	c.DB.Raw("SELECT * FROM users WHERE id=?", userID).Scan(&user)
	if user.Email == "" { // here given id so check with email
		return errors.New("invalid user id user doesn't exist")
	}

	query := `UPDATE users SET block_status = $1 WHERE id = $2`
	if c.DB.Exec(query, !user.BlockStatus, userID).Error != nil {
		return fmt.Errorf("faild update user block_status to %v", !user.BlockStatus)
	}
	return nil
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
