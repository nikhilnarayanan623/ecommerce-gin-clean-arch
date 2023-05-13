package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
)

type AdminRepository interface {
	FindAdminByEmail(ctx context.Context, email string) (domain.Admin, error)
	FindAdminByUserName(ctx context.Context, userName string) (domain.Admin, error)
	SaveAdmin(ctx context.Context, admin domain.Admin) error

	FindAllUser(ctx context.Context, pagination req.ReqPagination) (users []res.UserRespStrcut, err error)
	BlockUser(ctx context.Context, userID uint) error

	CreateFullSalesReport(ctc context.Context, reqData req.ReqSalesReport) (salesReport []res.SalesReport, err error)

	//stock side
	FindStockBySKU(ctx context.Context, sku string) (stock res.RespStock, err error)
	FindAllStockDetails(ctx context.Context, pagination req.ReqPagination) (stocks []res.RespStock, err error)
	UpdateStock(ctx context.Context, valuesToUpdate req.ReqUpdateStock) error
}
