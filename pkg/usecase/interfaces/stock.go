package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/response"
)

type StockUseCase interface {
	GetAllStockDetails(ctx context.Context, pagination request.Pagination) (stocks []response.Stock, err error)
	UpdateStockBySKU(ctx context.Context, updateDetails request.UpdateStock) error
}
