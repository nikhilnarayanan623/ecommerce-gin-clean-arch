package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper"
)

type ProductUseCase interface {
	GetCategory(ctx context.Context) (helper.RespFullCategory, error)
	AddCategory(ctx context.Context, category domain.Category) (domain.Category, error)

	AddVariation(ctx context.Context, vartaion domain.Variation) (domain.Variation, error)
	AddVariationOption(ctx context.Context, variationOption domain.VariationOption) (domain.VariationOption, error)

	GetProducts(ctx context.Context) ([]helper.ResponseProduct, error)
	AddProducts(ctx context.Context, product domain.Product) (domain.Product, error)

	AddProductItem(ctx context.Context, productItem helper.ReqProductItem) (domain.ProductItem, error)
	GetProductItems(ctx context.Context, product domain.Product) ([]helper.RespProductItems, error)
}
