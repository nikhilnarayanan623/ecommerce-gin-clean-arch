package usecase

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	service "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"golang.org/x/crypto/bcrypt"
)

type adminUseCase struct {
	adminRepo interfaces.AdminRepository
}

func NewAdminUseCase(repo interfaces.AdminRepository) service.AdminUseCase {

	return &adminUseCase{adminRepo: repo}
}

func (c *adminUseCase) Login(ctx context.Context, admin domain.Admin) (domain.Admin, any) {

	//validte the  admin
	err := validator.New().Struct(admin)

	if err != nil {
		errMap := map[string]string{}
		for _, er := range err.(validator.ValidationErrors) {
			errMap[er.Field()] = "Enter this field properly"
		}
		return admin, errMap
	}
	// get the admin from database
	dbAdmin, dbErr := c.adminRepo.FindAdmin(ctx, admin)

	if dbErr != nil {
		return admin, map[string]string{"msg": "db error"}
	}

	// check db password with given password
	if bcrypt.CompareHashAndPassword([]byte(dbAdmin.Password), []byte(admin.Password)) != nil {
		return admin, map[string]string{"msg": "Entered Passsword is incorrect"}
	}

	return dbAdmin, nil
}

func (c *adminUseCase) FindAllUser(ctx context.Context) ([]domain.Users, error) {

	users, err := c.adminRepo.FindAllUser(ctx)

	return users, err
}

func (c *adminUseCase) AddCategory(ctx context.Context, category domain.Category) (domain.Category, any) {

	//validate the given category name

	err := validator.New().Struct(category)

	if err != nil {
		errMap := map[string]string{}

		for _, er := range err.(validator.ValidationErrors) {
			errMap[er.Field()] = "Enter this field properly"
		}

		return category, errMap
	}

	productCategory, dbErr := c.adminRepo.AddCategory(ctx, category)

	return productCategory, dbErr
}
