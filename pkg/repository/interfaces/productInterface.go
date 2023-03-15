package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper"
)

type ProductRepository interface {
	GetProducts(ctx context.Context) ([]helper.ResponseProduct, error)
	AddProducts(ctx context.Context, product domain.Product) (domain.Product, error)

	GetProductItems(ctx context.Context, product domain.Product) ([]helper.RespProductItems, error)
	AddProductItem(ctx context.Context, productItem helper.ReqProductItem) (domain.ProductItem, error)

	GetCategory(ctx context.Context) (helper.RespFullCategory, error)
	AddCategory(ctx context.Context, productCategory domain.Category) (domain.Category, error)

	AddVariation(ctx context.Context, variation domain.Variation) (domain.Variation, error)
	AddVariationOption(ctx context.Context, variationOption domain.VariationOption) (domain.VariationOption, error)
}
