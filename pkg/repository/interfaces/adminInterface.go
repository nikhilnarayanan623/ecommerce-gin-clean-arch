package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper"
)

type AdminRepository interface {
	FindAdmin(ctx context.Context, admin domain.Admin) (domain.Admin, error)
	SaveAdmin(ctx context.Context, admin domain.Admin) (domain.Admin, any)
	FindAllUser(ctx context.Context) ([]domain.Users, error)
	BlockUser(ctx context.Context, user domain.Users) (domain.Users, any)
	GetCategory(ctx context.Context) ([]helper.RespCategory, any)
	AddCategory(ctx context.Context, productCategory domain.Category) (helper.RespCategory, any)
	AddProducts(ctx context.Context, product domain.Product) (domain.Product, any)
}
