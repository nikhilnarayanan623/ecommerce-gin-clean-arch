package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	service "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/response"
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

var (
	ErrInvalidUserID = errors.New("invalid user id")
	ErrInvalidSKU    = errors.New("invalid sku")
)

func (c *adminUseCase) SignUp(ctx context.Context, loginDetails domain.Admin) error {

	existAdmin, err := c.adminRepo.FindAdminByEmail(ctx, loginDetails.Email)
	if err != nil {
		return err
	} else if existAdmin.ID != 0 {
		return errors.New("can't save admin \nan admin already exist with this email")
	}

	existAdmin, err = c.adminRepo.FindAdminByUserName(ctx, loginDetails.UserName)
	if err != nil {
		return err
	} else if existAdmin.ID != 0 {
		return errors.New("can't save admin \nan admin already exist with this user_name")
	}

	// generate a hashed password for admin
	hashPass, err := bcrypt.GenerateFromPassword([]byte(loginDetails.Password), 10)

	if err != nil {
		return errors.New("failed to generate hashed password for admin")
	}
	// set the hashed password on the admin
	loginDetails.Password = string(hashPass)

	return c.adminRepo.SaveAdmin(ctx, loginDetails)
}

func (c *adminUseCase) FindAllUser(ctx context.Context, pagination request.Pagination) (users []response.User, err error) {

	users, err = c.adminRepo.FindAllUser(ctx, pagination)

	return users, err
}

// Block User
func (c *adminUseCase) BlockOrUnBlockUser(ctx context.Context, blockDetails request.BlockUser) error {

	userToBlock, err := c.userRepo.FindUserByUserID(ctx, blockDetails.UserID)
	if err != nil {
		return fmt.Errorf("failed to find user \nerror:%w", err)
	} else if userToBlock.ID == 0 {
		return ErrInvalidUserID
	}

	if userToBlock.BlockStatus == blockDetails.Block {
		return ErrSameBlockStatus
	}

	err = c.userRepo.UpdateBlockStatus(ctx, blockDetails.UserID, blockDetails.Block)
	if err != nil {
		return fmt.Errorf("failed to update user block status \nerror:%v", err.Error())
	}
	return nil
}

func (c *adminUseCase) GetFullSalesReport(ctx context.Context, requestData request.SalesReport) (salesReport []response.SalesReport, err error) {
	salesReport, err = c.adminRepo.CreateFullSalesReport(ctx, requestData)

	if err != nil {
		return salesReport, err
	}

	log.Printf("successfully got sales report from %v to %v of limit %v",
		requestData.StartDate, requestData.EndDate, requestData.Pagination.Count)

	return salesReport, nil
}

func (c *adminUseCase) GetAllStockDetails(ctx context.Context, pagination request.Pagination) (stocks []response.Stock, err error) {
	stocks, err = c.adminRepo.FindAllStockDetails(ctx, pagination)

	if err != nil {
		return stocks, err
	}

	log.Printf("successfully got stock details")
	return stocks, nil
}

func (c *adminUseCase) UpdateStockBySKU(ctx context.Context, updateDetails request.UpdateStock) error {

	stock, err := c.adminRepo.FindStockBySKU(ctx, updateDetails.SKU)
	if err != nil {
		return err
	}
	if stock.SKU == "" {
		return ErrInvalidSKU
	}

	err = c.adminRepo.UpdateStock(ctx, updateDetails)

	if err != nil {
		return err
	}

	log.Printf("successfully updated of stock details of stock with sku %v", updateDetails.SKU)
	return nil
}
