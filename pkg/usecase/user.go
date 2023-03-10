package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/jinzhu/copier"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	service "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
	"golang.org/x/crypto/bcrypt"
)

type userUserCase struct {
	userRepo interfaces.UserRepository
}

func NewUserUseCase(repo interfaces.UserRepository) service.UserUseCase {
	return &userUserCase{userRepo: repo}
}

func (c *userUserCase) Login(ctx context.Context, user domain.Users) (domain.Users, error) {

	dbUser, dberr := c.userRepo.FindUser(ctx, user)

	// check user found or not
	if dberr != nil {
		return user, dberr
	}

	// check user block_status user is blocked or not
	if user.BlockStatus {
		return user, errors.New("user blocked by admin")
	}

	//check the user password with dbPassword
	if bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)) != nil {
		return user, errors.New("entered password is wrong")
	}

	return dbUser, nil
}

func (c *userUserCase) LoginOtp(ctx context.Context, user domain.Users) (domain.Users, error) {

	user, err := c.userRepo.FindUser(ctx, user)

	if err != nil {
		return user, errors.New("can't find the user")
	}

	// check user block_status user is blocked or not
	if user.BlockStatus {
		return user, errors.New("user blocked by admin")
	}

	return user, nil
}

func (c *userUserCase) Signup(ctx context.Context, user domain.Users) (domain.Users, any) {
	//hash the password
	hashPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	if err != nil {
		return user, map[string]string{"error": "error to hash the password"}
	}
	user.Password = string(hashPass)

	return c.userRepo.SaveUser(ctx, user)
}

func (c *userUserCase) ShowAllProducts(ctx context.Context) ([]helper.ResponseProduct, any) {

	products, err := c.userRepo.GetAllProducts(ctx)

	if err != nil {
		return nil, map[string]string{"Error": "Can't get the products"}
	}

	return products, err
}
func (c *userUserCase) GetProductItems(ctx context.Context, product domain.Product) ([]domain.ProductItem, any) {

	productsItem, err := c.userRepo.GetProductItems(ctx, product)

	if err != nil {
		return nil, map[string]string{"Error": "To get products item"}
	}

	return productsItem, nil
}

func (c *userUserCase) GetCartItems(ctx context.Context, userId uint) (helper.ResCart, any) {

	return c.userRepo.GetCartItems(ctx, userId)
}

func (c *userUserCase) Home(ctx context.Context, userId uint) (helper.UserRespStrcut, any) {

	var user = domain.Users{ID: userId}
	user, err := c.userRepo.FindUser(ctx, user)

	if err != nil {
		return helper.UserRespStrcut{}, err
	}

	// everything ok then responce 200 with user details
	var response helper.UserRespStrcut
	fmt.Println(user, "user")
	copier.Copy(&response, &user) // copy required data only

	// everything is ok then return dbUser
	return response, nil

}
