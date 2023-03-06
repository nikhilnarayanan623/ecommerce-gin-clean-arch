package usecase

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/repository/interfaces"
	service "github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/usecase/interfaces"
)

type userUserCase struct {
	userRepo interfaces.UserRepository
}

func NewUserUseCase(repo interfaces.UserRepository) service.UserUseCase {
	return &userUserCase{userRepo: repo}
}

func (c *userUserCase) Login(ctx context.Context, user domain.Users) (domain.Users, any) {

	user, err := c.userRepo.FindUser(ctx, user)

	// check user found or not
	if err != nil {
		return user, map[string]string{"Error": "Can't find the user"}
	}

	// check user block_status user is blocked or not
	if user.BlockStatus {
		return user, map[string]string{"Error": "User Blocked By Admin"}
	}

	// everything is ok then return dbUser
	return user, nil
}

func (c *userUserCase) SaveUser(ctx context.Context, user domain.Users) (domain.Users, any) {

	// validate user values
	if err := validator.New().Struct(user); err != nil {

		errorMap := map[string]string{}
		for _, er := range err.(validator.ValidationErrors) {
			errorMap[er.Field()] = "Enter This field Properly"
		}

		return user, errorMap
	}

	user, err := c.userRepo.SaveUser(ctx, user)

	return user, err
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
