package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
)

type AdminRepository interface {
	FindAdmin(ctx context.Context, admin domain.Admin) (domain.Admin, any)
	SaveAdmin(ctx context.Context, admin domain.Admin) (domain.Admin, any)
	FindAllUser(ctx context.Context) ([]domain.Users, error)
	BlockUser(ctx context.Context, user domain.Users) (domain.Users, any)
	AddCategory(ctx context.Context, productCategory domain.Category) (domain.Category, any)
}
