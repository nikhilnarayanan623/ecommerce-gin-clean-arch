package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper"
)

type UserUseCase interface {
	Signup(ctx context.Context, user domain.Users) (domain.Users, any)
	Login(ctx context.Context, user domain.Users) (domain.Users, error)
	LoginOtp(ctx context.Context, user domain.Users) (domain.Users, error)

	Home(ctx context.Context, id uint) (helper.UserRespStrcut, any)
	ShowAllProducts(ctx context.Context) ([]helper.ResponseProduct, any)                     // show all products
	GetProductItems(ctx context.Context, product domain.Product) ([]domain.ProductItem, any) // to get all product items of a specific product
	GetCartItems(ctx context.Context, userId uint) (helper.ResCart, any)
}
