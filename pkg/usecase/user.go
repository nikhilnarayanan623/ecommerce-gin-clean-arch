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

type userUserCase struct {
	userRepo interfaces.UserRepository
}

func NewUserUseCase(repo interfaces.UserRepository) service.UserUseCase {
	return &userUserCase{userRepo: repo}
}

func (c *userUserCase) Login(ctx context.Context, body helper.LoginStruct) (helper.UserRespStrcut, any) {

	// first validate the struct(user)
	validate := validator.New()
	validate.RegisterValidation("login", helper.CustomLoginValidator) // custom validator for login

	if err := validate.Struct(body); err != nil {
		errMap := map[string]string{}

		for _, er := range err.(validator.ValidationErrors) {
			errMap[er.Field()] = "Enter this field properly"
		}
		helper.Reset() // for loing validatin reset count in that
		return helper.UserRespStrcut{}, errMap
	}
	helper.Reset() // for loing validatin reset count in that

	// if no error in validation then copy its field int users
	var user domain.Users
	copier.Copy(&user, &body)

	dbUser, dberr := c.userRepo.FindUser(ctx, user)

	// check user found or not
	if dberr != nil {
		return helper.UserRespStrcut{}, dberr
	}

	// check user block_status user is blocked or not
	if dbUser.BlockStatus {
		return helper.UserRespStrcut{}, map[string]string{"Error": "User Blocked By Admin"}
	}

	//check the user password with dbPassword
	if bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)) != nil {
		return helper.UserRespStrcut{}, map[string]string{"Error": "Entered Password is wrong"}
	}

	// everything ok then responce 200 with user details
	var response helper.UserRespStrcut
	copier.Copy(&response, &dbUser) // copy required data only

	// everything is ok then return dbUser
	return response, nil
}

func (c *userUserCase) Signup(ctx context.Context, user domain.Users) (domain.Users, any) {

	// validate user values
	if err := validator.New().Struct(user); err != nil {

		errorMap := map[string]string{}
		for _, er := range err.(validator.ValidationErrors) {
			errorMap[er.Field()] = "Enter This field Properly"
		}

		return user, errorMap
	}

	hashPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	if err != nil {
		return user, map[string]string{"error": "error to hash the password"}
	}
	user.Password = string(hashPass)

	return c.userRepo.SaveUser(ctx, user)
}

func (c *userUserCase) ShowAllProducts(ctx context.Context) ([]domain.Product, any) {

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
