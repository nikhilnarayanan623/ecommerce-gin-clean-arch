package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/request"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/api/handler/response"
	"github.com/nikhilnarayanan623/ecommerce-gin-clean-arch/pkg/domain"
)

type ProductUseCase interface {
	FindAllCategories(ctx context.Context, pagination request.Pagination) ([]response.Category, error)
	SaveCategory(ctx context.Context, categoryName string) error
	SaveSubCategory(ctx context.Context, subCategory request.SubCategory) error

	// variations
	SaveVariation(ctx context.Context, categoryID uint, variationNames []string) error
	SaveVariationOption(ctx context.Context, variationID uint, variationOptionValues []string) error

	FindAllVariationsAndItsValues(ctx context.Context, categoryID uint) ([]response.Variation, error)

	// products
	FindAllProducts(ctx context.Context, pagination request.Pagination) (products []response.Product, err error)
	SaveProduct(ctx context.Context, product request.Product) error
	UpdateProduct(ctx context.Context, product domain.Product) error

	SaveProductItem(ctx context.Context, productID uint, productItem request.ProductItem) error
	FindAllProductItems(ctx context.Context, productID uint) ([]response.ProductItems, error)
}
