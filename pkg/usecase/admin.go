package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jinzhu/copier"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	service "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
	"golang.org/x/crypto/bcrypt"
)

type adminUseCase struct {
	adminRepo interfaces.AdminRepository
}

func NewAdminUseCase(repo interfaces.AdminRepository) service.AdminUseCase {

	return &adminUseCase{adminRepo: repo}
}

func (c *adminUseCase) SignUp(ctx context.Context, admin domain.Admin) error {

	if admin, err := c.adminRepo.FindAdmin(ctx, admin); err != nil {
		return err
	} else if admin.ID != 0 {
		return errors.New("can't save admin already exist with this details")
	}

	// generate a hashed password for admin
	hashPass, err := bcrypt.GenerateFromPassword([]byte(admin.Password), 10)

	if err != nil {
		return errors.New("faild to generate hashed password for admin")
	}
	// set the hashed password on the admin
	admin.Password = string(hashPass)

	return c.adminRepo.SaveAdmin(ctx, admin)
}

func (c *adminUseCase) Login(ctx context.Context, admin domain.Admin) (domain.Admin, error) {

	// get the admin from database
	dbAdmin, err := c.adminRepo.FindAdmin(ctx, admin)
	if err != nil {
		return admin, err
	} else if dbAdmin.ID == 0 {
		return admin, errors.New("admin not exist with given details")
	}

	// check db password with given password
	if bcrypt.CompareHashAndPassword([]byte(dbAdmin.Password), []byte(admin.Password)) != nil {
		return admin, errors.New("wrong password")
	}

	return dbAdmin, nil
}

func (c *adminUseCase) FindAllUser(ctx context.Context, pagination req.ReqPagination) (users []res.UserRespStrcut, err error) {

	users, err = c.adminRepo.FindAllUser(ctx, pagination)

	if err != nil {
		return nil, err
	}

	// if no error then copy users details to an array responce struct
	var responce []res.UserRespStrcut
	copier.Copy(&responce, &users)

	return responce, nil
}

// to block or unblock a user
func (c *adminUseCase) BlockUser(ctx context.Context, userID uint) error {

	return c.adminRepo.BlockUser(ctx, userID)
}

func (c *adminUseCase) GetFullSalesReport(ctx context.Context, requestData req.ReqSalesReport) (salesReport []res.SalesReport, err error) {
	salesReport, err = c.adminRepo.CreateFullSalesReport(ctx, requestData)

	if err != nil {
		return salesReport, err
	}

	log.Printf("successfully got sales report from %v to %v of limit %v",
		requestData.StartDate, requestData.EndDate, requestData.Pagination.Count)

	return salesReport, nil
}

func (c *adminUseCase) GetAllStockDetails(ctx context.Context, pagination req.ReqPagination) (stocks []res.RespStock, err error) {
	stocks, err = c.adminRepo.FindAllStockDetails(ctx, pagination)

	if err != nil {
		return stocks, err
	}

	log.Printf("successfully got stock details")
	return stocks, nil
}

func (c *adminUseCase) UpdateStock(ctx context.Context, valuesToUpdate req.ReqUpdateStock) error {

	// validate the sku
	stock, err := c.adminRepo.FindStockBySKU(ctx, valuesToUpdate.SKU)
	if err != nil {
		return err
	} else if stock.ProductName == "" {
		return fmt.Errorf("invalid sku %v", valuesToUpdate.SKU)
	}

	// update the stock detils
	err = c.adminRepo.UpdateStock(ctx, valuesToUpdate)

	if err != nil {
		return err
	}

	log.Printf("successfully updated of stock details of stock with sku %v", valuesToUpdate.SKU)
	return nil
}
