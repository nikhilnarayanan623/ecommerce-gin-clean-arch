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
	userRepo  interfaces.UserRepository
}

func NewAdminUseCase(repo interfaces.AdminRepository, userRepo interfaces.UserRepository) service.AdminUseCase {

	return &adminUseCase{
		adminRepo: repo,
		userRepo:  userRepo,
	}
}

func (c *adminUseCase) SignUp(ctx context.Context, loginDetails domain.Admin) error {

	admin, err := c.adminRepo.FindAdminByEmail(ctx, loginDetails.Email)
	if err != nil {
		return err
	} else if admin.ID != 0 {
		return errors.New("can't save admin \nan admin already exist with this email")
	}

	admin, err = c.adminRepo.FindAdminByUserName(ctx, loginDetails.UserName)
	if err != nil {
		return err
	} else if admin.ID != 0 {
		return errors.New("can't save admin \nan admin already exist with this user_name")
	}

	// generate a hashed password for admin
	hashPass, err := bcrypt.GenerateFromPassword([]byte(loginDetails.Password), 10)

	if err != nil {
		return errors.New("faild to generate hashed password for admin")
	}
	// set the hashed password on the admin
	admin.Password = string(hashPass)

	return c.adminRepo.SaveAdmin(ctx, admin)
}

func (c *adminUseCase) FindAllUser(ctx context.Context, pagination req.Pagination) (users []res.User, err error) {

	users, err = c.adminRepo.FindAllUser(ctx, pagination)

	if err != nil {
		return nil, err
	}

	// if no error then copy users details to an array responce struct
	var responce []res.User
	copier.Copy(&responce, &users)

	return responce, nil
}

func (c *adminUseCase) BlockOrUblockUser(ctx context.Context, blockDetails req.BlockUser) error {

	userToBlock, err := c.userRepo.FindUserByUserID(ctx, blockDetails.UserID)
	if err != nil {
		return fmt.Errorf("faild to find user \nerror:%v", err.Error())
	} else if userToBlock.ID == 0 {
		return fmt.Errorf("invalid user_id")
	}

	if userToBlock.BlockStatus == blockDetails.Block {
		return fmt.Errorf("user block status already in given status")
	}

	err = c.userRepo.UpdateBlockStatus(ctx, blockDetails.UserID, blockDetails.Block)
	if err != nil {
		return fmt.Errorf("faild to update user block status \nerror:%v", err.Error())
	}
	return nil
}

func (c *adminUseCase) GetFullSalesReport(ctx context.Context, requestData req.SalesReport) (salesReport []res.SalesReport, err error) {
	salesReport, err = c.adminRepo.CreateFullSalesReport(ctx, requestData)

	if err != nil {
		return salesReport, err
	}

	log.Printf("successfully got sales report from %v to %v of limit %v",
		requestData.StartDate, requestData.EndDate, requestData.Pagination.Count)

	return salesReport, nil
}

func (c *adminUseCase) GetAllStockDetails(ctx context.Context, pagination req.Pagination) (stocks []res.Stock, err error) {
	stocks, err = c.adminRepo.FindAllStockDetails(ctx, pagination)

	if err != nil {
		return stocks, err
	}

	log.Printf("successfully got stock details")
	return stocks, nil
}

func (c *adminUseCase) UpdateStock(ctx context.Context, valuesToUpdate req.UpdateStock) error {

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
