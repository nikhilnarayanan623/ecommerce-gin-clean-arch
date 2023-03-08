package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper"
)

type UserRepository interface {
	FindUser(ctx context.Context, user domain.Users) (domain.Users, any)
	SaveUser(ctx context.Context, user domain.Users) (domain.Users, any)
	GetAllProducts(ctx context.Context) ([]domain.Product, any)
	GetProductItems(ctx context.Context, product domain.Product) ([]domain.ProductItem, any)
	GetCartItems(ctx context.Context, userId uint) (helper.ResCart, any)
}
