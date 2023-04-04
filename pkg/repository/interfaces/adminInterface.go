package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/res"
)

type AdminRepository interface {
	FindAdmin(ctx context.Context, admin domain.Admin) (domain.Admin, error)
	SaveAdmin(ctx context.Context, admin domain.Admin) error

	FindAllUser(ctx context.Context) ([]domain.User, error)
	BlockUser(ctx context.Context, userID uint) error

	CreateFullSalesReport(ctc context.Context) ([]res.SalesReport, error)
}
