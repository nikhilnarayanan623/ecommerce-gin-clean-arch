package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/utils/res"
)

type AdminRepository interface {
	FindAdmin(ctx context.Context, admin domain.Admin) (domain.Admin, error)
	SaveAdmin(ctx context.Context, admin domain.Admin) error

	FindAllUser(ctx context.Context, pagination req.ReqPagination) (users []res.UserRespStrcut, err error)
	BlockUser(ctx context.Context, userID uint) error

	CreateFullSalesReport(ctc context.Context, reqData req.ReqSalesReport) (salesReport []res.SalesReport, err error)
}
