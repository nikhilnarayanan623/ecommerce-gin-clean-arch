package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
)

type AdminUseCase interface {
	SignUp(ctx context.Context, admin domain.Admin) error

	FindAllUser(ctx context.Context, pagination req.Pagination) (users []res.User, err error)
	BlockOrUblockUser(ctx context.Context, blockDetails req.BlockUser) error

	GetFullSalesReport(ctx context.Context, requestData req.SalesReport) (salesReport []res.SalesReport, err error)

	// stock side
	GetAllStockDetails(ctx context.Context, pagination req.Pagination) (stocks []res.Stock, err error)
	UpdateStock(ctx context.Context, valuesToUpdate req.UpdateStock) error
}

// GetCategory(ctx context.Context) (helper.Category, any)
// 	SetCategory(ctx context.Context, body helper.Category)
