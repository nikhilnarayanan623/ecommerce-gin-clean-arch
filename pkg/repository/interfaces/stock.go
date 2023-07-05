package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/response"
)

type StockRepository interface {
	FindAll(ctx context.Context, pagination request.Pagination) (stocks []response.Stock, err error)
	Update(ctx context.Context, updateValues request.UpdateStock) error
}
