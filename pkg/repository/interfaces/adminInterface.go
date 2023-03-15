package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
)

type AdminRepository interface {
	FindAdmin(ctx context.Context, admin domain.Admin) (domain.Admin, error)
	SaveAdmin(ctx context.Context, admin domain.Admin) (domain.Admin, error)

	FindAllUser(ctx context.Context) ([]domain.Users, error)
	BlockUser(ctx context.Context, user domain.Users) (domain.Users, error)
}
