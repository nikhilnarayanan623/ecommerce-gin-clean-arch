package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/req"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/helper/res"
)

type ProductUseCase interface {
	GetCategory(ctx context.Context) (res.RespFullCategory, error)
	AddCategory(ctx context.Context, category domain.Category) (domain.Category, error)

	AddVariation(ctx context.Context, vartaion domain.Variation) (domain.Variation, error)
	AddVariationOption(ctx context.Context, variationOption domain.VariationOption) (domain.VariationOption, error)

	GetProducts(ctx context.Context) ([]res.ResponseProduct, error)
	AddProduct(ctx context.Context, product domain.Product) (domain.Product, error)

	AddProductItem(ctx context.Context, productItem req.ReqProductItem) (domain.ProductItem, error)
	GetProductItems(ctx context.Context, product domain.Product) ([]res.RespProductItems, error)
}
