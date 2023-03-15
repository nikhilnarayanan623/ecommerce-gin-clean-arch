package usecase

import (
	"context"
	"errors"

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

func (c *adminUseCase) SignUp(ctx context.Context, admin domain.Admin) (domain.Admin, error) {

	// generate a hashed password for admin
	hashPass, err := bcrypt.GenerateFromPassword([]byte(admin.Password), 10)

	if err != nil {
		return admin, errors.New("faild to generate hashed password for admin")
	}
	// set the hashed password on the admin
	admin.Password = string(hashPass)

	return c.adminRepo.SaveAdmin(ctx, admin)
}

func (c *adminUseCase) Login(ctx context.Context, admin domain.Admin) (domain.Admin, error) {

	// get the admin from database
	dbAdmin, dbErr := c.adminRepo.FindAdmin(ctx, admin)

	if dbErr != nil {
		return admin, dbErr
	}

	// check db password with given password
	if bcrypt.CompareHashAndPassword([]byte(dbAdmin.Password), []byte(admin.Password)) != nil {
		return admin, errors.New("entered passsword is incorrect")
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

// to block or unblock a user
func (c *adminUseCase) BlockUser(ctx context.Context, user domain.Users) (domain.Users, error) {

	return c.adminRepo.BlockUser(ctx, user)
}
