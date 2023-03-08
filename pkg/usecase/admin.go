package usecase

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper"
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

func (c *adminUseCase) SignUp(ctx context.Context, admin domain.Admin) (domain.Admin, any) {

	//validate the struct
	if err := validator.New().Struct(admin); err != nil {
		errMap := map[string]string{}

		for _, er := range err.(validator.ValidationErrors) {
			errMap[er.Field()] = "Enter this field properly"
		}
		return admin, errMap
	}

	// then hash the password
	hashPass, err := bcrypt.GenerateFromPassword([]byte(admin.Password), 10)

	if err != nil {
		return admin, map[string]string{"error": "Faild to hash the password"}
	}
	// set the hashed password on the admin
	admin.Password = string(hashPass)

	return c.adminRepo.SaveAdmin(ctx, admin)
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
		return admin, dbErr
	}

	// check db password with given password
	if bcrypt.CompareHashAndPassword([]byte(dbAdmin.Password), []byte(admin.Password)) != nil {
		return admin, map[string]string{"msg": "Entered Passsword is incorrect"}
	}

	return dbAdmin, nil
}

func (c *adminUseCase) FindAllUser(ctx context.Context) ([]helper.UserRespStrcut, error) {

	users, err := c.adminRepo.FindAllUser(ctx)

	if err != nil {
		return nil, err
	}

	// if no error then copy users details to an array responce struct
	var responce []helper.UserRespStrcut
	copier.Copy(&responce, &users)

	return responce, nil
}

func (c *adminUseCase) BlockUser(ctx context.Context, request helper.BlockStruct) (domain.Users, any) {

	// validate the struct
	if request.ID <= 0 {
		return domain.Users{}, map[string]string{"Error": "Ivalid Id"}
	}

	//copy the id from req to user
	var user domain.Users
	copier.Copy(&user, &request)

	return c.adminRepo.BlockUser(ctx, user)
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
