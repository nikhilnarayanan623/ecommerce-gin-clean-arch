package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
)

type AdminUseCase interface {
	Login(ctx context.Context, admin domain.Admin) (domain.Admin, error)
	SignUp(ctx context.Context, admin domain.Admin) error

	FindAllUser(ctx context.Context, pagination req.ReqPagination) (users []res.UserRespStrcut, err error)
	BlockUser(ctx context.Context, userID uint) error

	GetFullSalesReport(ctx context.Context, requestData req.ReqSalesReport) (salesReport []res.SalesReport, err error)

	// stock side
	GetAllStockDetails(ctx context.Context, pagination req.ReqPagination) (stocks []res.RespStock, err error)
	UpdateStock(ctx context.Context, valuesToUpdate req.ReqUpdateStock) error
}

// GetCategory(ctx context.Context) (helper.ReqCategory, any)
// 	SetCategory(ctx context.Context, body helper.ReqCategory)
